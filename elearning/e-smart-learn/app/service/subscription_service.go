package service

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
	billingportalsession "github.com/stripe/stripe-go/v83/billingportal/session"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/paymentintent"
	"github.com/stripe/stripe-go/v83/price"
	"github.com/stripe/stripe-go/v83/product"
	"github.com/stripe/stripe-go/v83/promotioncode"
	stripesubscription "github.com/stripe/stripe-go/v83/subscription"
	"github.com/stripe/stripe-go/v83/webhook"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	GetTransactionsHistory(ctx context.Context, userID uuid.UUID) (*dto.ListTransactionsHistoryResponse, error)
	CreateSubscriptionCheckoutSession(ctx context.Context, userID uuid.UUID, request dto.CreateSubscriptionCheckoutSessionRequest) (*dto.CheckoutSessionResponse, error)
	SyncPendingStripeSubscriptions(ctx context.Context) error
	SyncPendingStripeCoursePurchases(ctx context.Context) error
	GetSubscribers(ctx context.Context, limit int, offset int) ([]*dto.SubscriberResponse, int64, error)
	CreateCourseCheckoutSession(ctx context.Context, userID, courseID uuid.UUID, request dto.CreateCoursePurchaseCheckoutSessionRequest) (*dto.CheckoutSessionResponse, error)

	PublishCourseAndCreateProductCatalog(ctx context.Context, courseID uuid.UUID) (*dto.CourseResponse, error)
	SyncPublishedCourseCatalog(ctx context.Context, courseID uuid.UUID) (*dto.CourseResponse, error)

	HandleStripeWebhook(ctx context.Context, payload []byte, signature string) error
	GetMySubscription(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error)
	CancelAtPeriodEnd(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error)
	Resume(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error)
	CreateBillingPortalSession(ctx context.Context, userID uuid.UUID) (*dto.BillingPortalResponse, error)
	HasActiveSubscription(ctx context.Context, userID uuid.UUID) (bool, error)

	// Get subscription retention statistics for admin dashboard
	GetMemberRetention(ctx context.Context) (*dto.MemberRetentionResponse, error)
}

type subscriptionService struct {
	userRepository                       repository.UserRepository
	planRepository                       repository.PlanRepository
	courseRepository                     repository.CourseRepository
	enrollmentRepository                 repository.EnrollmentRepository
	subscriptionRepository               repository.SubscriptionRepository
	paymentRepository                    repository.PaymentRepository
	subscriptionRevenueShareRepository   repository.SubscriptionRevenueShareRepository
	coursePurchaseRevenueShareRepository repository.CoursePurchaseRevenueShareRepository
	coursePurchaseRepository             repository.CoursePurchaseRepository
	coursePurchaseDetailRepository       repository.CoursePurchaseDetailRepository
	courseCouponRepository               repository.CouponRepository
	stripeEventRepository                repository.StripeEventRepository
	stripeConfig                         *config.StripeConfig
}

func NewSubscriptionService(
	userRepository repository.UserRepository,
	planRepository repository.PlanRepository,
	courseRepository repository.CourseRepository,
	enrollmentRepository repository.EnrollmentRepository,
	subscriptionRepository repository.SubscriptionRepository,
	paymentRepository repository.PaymentRepository,
	subscriptionRevenueShareRepository repository.SubscriptionRevenueShareRepository,
	coursePurchaseRevenueShareRepository repository.CoursePurchaseRevenueShareRepository,
	coursePurchaseRepository repository.CoursePurchaseRepository,
	coursePurchaseDetailRepository repository.CoursePurchaseDetailRepository,
	courseCouponRepository repository.CouponRepository,
	stripeEventRepository repository.StripeEventRepository,
	stripeConfig *config.StripeConfig,
) SubscriptionService {
	stripe.Key = stripeConfig.SecretKey

	return &subscriptionService{
		userRepository:                       userRepository,
		planRepository:                       planRepository,
		courseRepository:                     courseRepository,
		enrollmentRepository:                 enrollmentRepository,
		subscriptionRepository:               subscriptionRepository,
		paymentRepository:                    paymentRepository,
		subscriptionRevenueShareRepository:   subscriptionRevenueShareRepository,
		coursePurchaseRevenueShareRepository: coursePurchaseRevenueShareRepository,
		coursePurchaseRepository:             coursePurchaseRepository,
		coursePurchaseDetailRepository:       coursePurchaseDetailRepository,
		courseCouponRepository:               courseCouponRepository,
		stripeEventRepository:                stripeEventRepository,
		stripeConfig:                         stripeConfig,
	}
}

func (s *subscriptionService) GetTransactionsHistory(ctx context.Context, userID uuid.UUID) (*dto.ListTransactionsHistoryResponse, error) {
	coursePurchases, err := s.coursePurchaseRepository.ListPaidByUserID(ctx, userID, []repository.Preload{repository.Preload("Details.Course")})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get course transaction history")
	}
	subscriptions, err := s.subscriptionRepository.ListByUserID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get subscription transaction history")
	}
	var totalAmount float64
	transactions := make([]*dto.TransactionHistoryResponse, 0, len(coursePurchases)+len(subscriptions))
	for _, purchase := range coursePurchases {
		for _, detail := range purchase.Details {
			name := strings.TrimSpace(detail.Course.Title)
			if name == "" {
				name = "Course Purchase"
			}
			amount := float64(detail.Price) / 100
			transactions = append(transactions, &dto.TransactionHistoryResponse{
				ID:        purchase.ID.String(),
				Type:      "course_purchase",
				Name:      name,
				Amount:    amount,
				Currency:  strings.ToUpper(detail.Currency),
				Status:    purchase.Status,
				CreatedAt: purchase.CreatedAt,
			})
			totalAmount += amount
		}
	}
	for _, sub := range subscriptions {
		name := strings.TrimSpace(sub.PlanName)
		if name == "" {
			name = "Subscription"
		}
		transactions = append(transactions, &dto.TransactionHistoryResponse{
			ID:        sub.ID.String(),
			Type:      "subscription",
			Name:      name,
			Amount:    sub.PlanPrice,
			Currency:  strings.ToUpper(sub.PlanCurrency),
			Status:    sub.Status,
			CreatedAt: sub.StartedAt,
		})
		totalAmount += sub.PlanPrice
	}
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.After(transactions[j].CreatedAt)
	})

	return &dto.ListTransactionsHistoryResponse{
		Transactions: transactions,
		TotalAmount:  totalAmount,
	}, nil
}

func (s *subscriptionService) CreateSubscriptionCheckoutSession(ctx context.Context, userID uuid.UUID, request dto.CreateSubscriptionCheckoutSessionRequest) (*dto.CheckoutSessionResponse, error) {
	activeSub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to check current subscription")
	}
	if activeSub != nil {
		return nil, apperror.NewBadRequestError("You already have an active subscription")
	}

	planID, err := uuid.Parse(request.PlanID)
	if err != nil {
		return nil, apperror.NewBadRequestError("Invalid plan ID")
	}

	plan, err := s.planRepository.FindByID(ctx, planID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get plan")
	}
	if plan == nil || !plan.IsActive {
		return nil, apperror.NewNotFoundError("Plan not found")
	}

	priceID, billingCycle, err := s.resolvePlanPriceID(ctx, plan)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	customerID, err := s.ensureStripeCustomer(ctx, user)
	if err != nil {
		return nil, err
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(s.stripeConfig.SuccessURL),
		CancelURL:  stripe.String(s.stripeConfig.CancelURL),
		Customer:   stripe.String(customerID),
		ExpiresAt:  stripe.Int64(time.Now().UTC().Add(time.Duration(s.stripeConfig.CheckoutSessionExpirationSeconds) * time.Second).Unix()),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"user_id":       userID.String(),
			"plan_id":       plan.ID.String(),
			"billing_cycle": billingCycle,
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id":       userID.String(),
				"plan_id":       plan.ID.String(),
				"billing_cycle": billingCycle,
			},
		},
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe checkout session")
	}

	if _, err := s.subscriptionRepository.Create(ctx, &model.Subscription{
		UserID:                  userID,
		PlanID:                  plan.ID,
		PlanName:                plan.Name,
		PlanDescription:         plan.Description,
		PlanPrice:               plan.Price,
		PlanCurrency:            strings.ToLower(plan.Currency),
		PlanStripePriceID:       plan.StripePriceID,
		StripeCheckoutSessionID: session.ID,
		StripeCustomerID:        customerID,
		BillingCycle:            billingCycle,
		Status:                  string(stripe.SubscriptionStatusIncomplete),
		StartedAt:               time.Now().UTC(),
	}); err != nil {
		return nil, apperror.NewInternalServerError("Failed to create pending subscription")
	}

	return &dto.CheckoutSessionResponse{
		SessionID: session.ID,
		URL:       session.URL,
	}, nil
}

func (s *subscriptionService) CreateCourseCheckoutSession(ctx context.Context, userID, courseID uuid.UUID, request dto.CreateCoursePurchaseCheckoutSessionRequest) (*dto.CheckoutSessionResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.Status != consts.CoursePublished {
		return nil, apperror.NewBadRequestError("Course is not available for purchase")
	}
	if course.UserID == userID {
		return nil, apperror.NewBadRequestError("Cannot purchase your own course")
	}
	if course.Price <= 0 {
		return nil, apperror.NewBadRequestError("Course is free, no purchase is required")
	}

	if err := s.ensureCourseStripeProductCatalog(ctx, course); err != nil {
		return nil, err
	}

	isPaid, err := s.coursePurchaseRepository.ExistsPaidByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to verify previous purchase")
	}
	if isPaid {
		return nil, apperror.NewBadRequestError("Course already purchased")
	}

	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	customerID, err := s.ensureStripeCustomer(ctx, user)
	if err != nil {
		return nil, err
	}

	amountCents := course.StripeAmount
	if amountCents <= 0 {
		amountCents = int64(math.Round(course.Price * 100))
	}
	if amountCents <= 0 {
		return nil, apperror.NewBadRequestError("Invalid course price")
	}

	currency := strings.ToLower(course.StripeCurrency)
	if currency == "" {
		currency = strings.ToLower(s.stripeConfig.DefaultCurrency)
	}
	if currency == "" {
		currency = "usd"
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.stripeConfig.SuccessURL),
		CancelURL:  stripe.String(s.stripeConfig.CancelURL),
		Customer:   stripe.String(customerID),
		ExpiresAt:  stripe.Int64(time.Now().UTC().Add(time.Duration(s.stripeConfig.CheckoutSessionExpirationSeconds) * time.Second).Unix()),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(course.StripePriceID),
				Quantity: stripe.Int64(1),
			},
		},
		Metadata: map[string]string{
			"purchase_type": "course_purchase",
			"user_id":       userID.String(),
		},
	}

	var appliedCouponID *uuid.UUID
	finalAmountCents := amountCents

	if request.CouponCode != "" {
		courseCoupon, err := s.courseCouponRepository.GetByCode(ctx, strings.TrimSpace(request.CouponCode), nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate coupon")
		}
		if courseCoupon == nil || !courseCoupon.IsActive {
			return nil, apperror.NewBadRequestError("Coupon code is invalid or inactive")
		}

		promotion, err := promotioncode.Get(courseCoupon.StripePromotionCodeID, nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to verify coupon on Stripe")
		}
		if promotion == nil || !promotion.Active {
			return nil, apperror.NewBadRequestError("Coupon code is inactive")
		}
		if promotion.ExpiresAt > 0 && promotion.ExpiresAt < time.Now().UTC().Unix() {
			return nil, apperror.NewBadRequestError("Coupon code has expired")
		}
		if promotion.MaxRedemptions > 0 && promotion.TimesRedeemed >= promotion.MaxRedemptions {
			return nil, apperror.NewBadRequestError("Coupon code has reached redemption limit")
		}

		params.Discounts = []*stripe.CheckoutSessionDiscountParams{
			{PromotionCode: stripe.String(courseCoupon.StripePromotionCodeID)},
		}
		appliedCouponID = &courseCoupon.ID
		finalAmountCents = applyCouponDiscountAmount(amountCents, courseCoupon)
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe checkout session")
	}

	if err := s.createCoursePurchasePending(ctx, userID, appliedCouponID, []uuid.UUID{courseID}, session.ID, finalAmountCents, currency, map[uuid.UUID]int64{courseID: finalAmountCents}); err != nil {
		return nil, err
	}

	return &dto.CheckoutSessionResponse{SessionID: session.ID, URL: session.URL}, nil
}

func (s *subscriptionService) PublishCourseAndCreateProductCatalog(ctx context.Context, courseID uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.Price <= 0 {
		return nil, apperror.NewBadRequestError("Course price must be greater than zero")
	}
	if course.Status == consts.CoursePublished {
		return nil, apperror.NewBadRequestError("Course is already published")
	}

	course.Status = consts.CoursePublished
	if err := s.ensureCourseStripeProductCatalog(ctx, course); err != nil {
		return nil, err
	}

	updated, err := s.courseRepository.Updates(ctx, course)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to publish course")
	}

	return dto.NewCourseDetailResponse(updated), nil
}

func (s *subscriptionService) SyncPublishedCourseCatalog(ctx context.Context, courseID uuid.UUID) (*dto.CourseResponse, error) {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get course")
	}
	if course == nil {
		return nil, apperror.NewNotFoundError("Course not found")
	}
	if course.Status != consts.CoursePublished {
		return nil, apperror.NewBadRequestError("Course is not published")
	}

	if err := s.ensureCourseStripeProductCatalog(ctx, course); err != nil {
		return nil, err
	}

	return dto.NewCourseDetailResponse(course), nil
}

func (s *subscriptionService) HandleStripeWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.stripeConfig.WebhookSecret)
	if err != nil {
		return apperror.NewUnauthorizedError("Invalid Stripe webhook signature")
	}

	existingEvent, err := s.stripeEventRepository.GetByEventID(ctx, event.ID)
	if err != nil {
		return apperror.NewInternalServerError("Failed to verify webhook idempotency")
	}
	if existingEvent != nil {
		return nil
	}

	_, err = s.stripeEventRepository.Create(ctx, &model.StripeEvent{
		EventID:   event.ID,
		EventType: string(event.Type),
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil
		}
		return apperror.NewInternalServerError("Failed to persist webhook event")
	}

	switch event.Type {
	case "checkout.session.completed":
		checkoutSession, err := decodeCheckoutSession(event.Data.Raw)
		if err != nil {
			return err
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModePayment {
			return s.handleCoursePurchaseCheckoutCompleted(ctx, &checkoutSession)
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModeSubscription {
			return s.handleCheckoutSessionCompleted(ctx, &checkoutSession)
		}
		return nil
	case "checkout.session.expired":
		checkoutSession, err := decodeCheckoutSession(event.Data.Raw)
		if err != nil {
			return err
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModePayment {
			return s.handleCoursePurchaseCheckoutFailed(ctx, &checkoutSession)
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModeSubscription {
			return s.handleSubscriptionCheckoutExpired(ctx, checkoutSession.ID)
		}
		return nil
	case "checkout.session.async_payment_succeeded":
		checkoutSession, err := decodeCheckoutSession(event.Data.Raw)
		if err != nil {
			return err
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModePayment {
			return s.handleCoursePurchaseCheckoutCompleted(ctx, &checkoutSession)
		}
		return nil
	case "checkout.session.async_payment_failed":
		checkoutSession, err := decodeCheckoutSession(event.Data.Raw)
		if err != nil {
			return err
		}
		if checkoutSession.Mode == stripe.CheckoutSessionModePayment {
			return s.handleCoursePurchaseCheckoutFailed(ctx, &checkoutSession)
		}
		return nil
	case "customer.subscription.created", "customer.subscription.updated", "customer.subscription.deleted":
		sub, err := decodeStripeSubscription(event.Data.Raw)
		if err != nil {
			return err
		}
		return s.syncSubscription(ctx, "", &sub)
	case "invoice.paid":
		inv, err := decodeStripeInvoice(event.Data.Raw)
		if err != nil {
			return err
		}
		return s.syncPayment(ctx, &inv, string(consts.PaymentStatusSucceeded))
	case "invoice.payment_failed":
		inv, err := decodeStripeInvoice(event.Data.Raw)
		if err != nil {
			return err
		}
		return s.syncPayment(ctx, &inv, string(consts.PaymentStatusFailed))
	default:
		return nil
	}
}

func (s *subscriptionService) GetMySubscription(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error) {
	sub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, []repository.Preload{repository.Plan})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get active subscription")
	}
	if sub == nil {
		sub, err = s.subscriptionRepository.GetLatestByUserID(ctx, userID, []repository.Preload{repository.Plan})
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to get subscription")
		}
	}
	if sub == nil {
		return nil, apperror.NewNotFoundError("Subscription not found")
	}

	return dto.NewMySubscriptionResponse(sub), nil
}

func (s *subscriptionService) CancelAtPeriodEnd(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error) {
	sub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, []repository.Preload{repository.Plan})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get active subscription")
	}

	if sub == nil || sub.StripeSubscriptionID == "" {
		return nil, apperror.NewNotFoundError("Active subscription not found")
	}
	if sub.Status == string(stripe.SubscriptionStatusCanceled) {
		return nil, apperror.NewBadRequestError("Subscription already canceled")
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(true),
	}
	stripeSub, err := stripesubscription.Update(sub.StripeSubscriptionID, params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to cancel subscription")
	}
	if err := s.syncSubscription(ctx, "", stripeSub); err != nil {
		return nil, err
	}

	updated, err := s.subscriptionRepository.FindByID(ctx, sub.ID, []repository.Preload{repository.Plan})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to load updated subscription")
	}

	return dto.NewMySubscriptionResponse(updated), nil
}

func (s *subscriptionService) Resume(ctx context.Context, userID uuid.UUID) (*dto.MySubscriptionResponse, error) {
	sub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, []repository.Preload{repository.Plan})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get active subscription")
	}
	if sub == nil || sub.StripeSubscriptionID == "" {
		return nil, apperror.NewNotFoundError("Active subscription not found")
	}
	if !sub.CancelAtPeriodEnd {
		return nil, apperror.NewBadRequestError("Subscription is not scheduled to cancel")
	}

	params := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
	}
	stripeSub, err := stripesubscription.Update(sub.StripeSubscriptionID, params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to resume subscription")
	}
	if err := s.syncSubscription(ctx, "", stripeSub); err != nil {
		return nil, err
	}

	updated, err := s.subscriptionRepository.FindByID(ctx, sub.ID, []repository.Preload{repository.Plan})
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to load updated subscription")
	}

	return dto.NewMySubscriptionResponse(updated), nil
}

func (s *subscriptionService) CreateBillingPortalSession(ctx context.Context, userID uuid.UUID) (*dto.BillingPortalResponse, error) {
	user, err := s.userRepository.FindByID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get user")
	}
	if user == nil {
		return nil, apperror.NewNotFoundError("User not found")
	}

	customerID, err := s.ensureStripeCustomer(ctx, user)
	if err != nil {
		return nil, err
	}

	params := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(customerID),
		ReturnURL: stripe.String(s.stripeConfig.BillingPortal),
	}

	portalSession, err := billingportalsession.New(params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create billing portal session")
	}

	return &dto.BillingPortalResponse{URL: portalSession.URL}, nil
}

func (s *subscriptionService) HasActiveSubscription(ctx context.Context, userID uuid.UUID) (bool, error) {
	sub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, nil)
	if err != nil {
		return false, apperror.NewInternalServerError("Failed to get active subscription")
	}
	return sub != nil, nil
}

func (s *subscriptionService) SyncPendingStripeSubscriptions(ctx context.Context) error {
	subscriptions, err := s.subscriptionRepository.ListPendingByCheckoutSessionID(ctx, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to list subscriptions")
	}

	for _, sub := range subscriptions {
		if sub == nil || strings.TrimSpace(sub.StripeCheckoutSessionID) == "" {
			continue
		}

		if err := s.syncStripeCheckoutSession(ctx, sub.StripeCheckoutSessionID); err != nil {
			return err
		}
	}

	return nil
}

func (s *subscriptionService) SyncPendingStripeCoursePurchases(ctx context.Context) error {
	purchases, err := s.coursePurchaseRepository.ListPendingByCheckoutSessionID(ctx, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to list pending purchases")
	}

	for _, purchase := range purchases {
		if purchase == nil || strings.TrimSpace(purchase.StripeCheckoutSessionID) == "" {
			continue
		}

		if err := s.syncCoursePurchaseCheckoutSession(ctx, purchase.StripeCheckoutSessionID); err != nil {
			return err
		}
	}

	return nil
}

func (s *subscriptionService) GetSubscribers(ctx context.Context, limit int, offset int) ([]*dto.SubscriberResponse, int64, error) {
	subs, total, err := s.subscriptionRepository.List(ctx, limit, offset, "created_at DESC", "", []repository.Preload{repository.User})
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to get subscribers")
	}
	result := dto.NewListSubscribersResponse(subs)
	return result, total, nil
}

func (s *subscriptionService) GetMemberRetention(ctx context.Context) (*dto.MemberRetentionResponse, error) {
	now := time.Now().UTC()
	previousMonth := now.AddDate(0, -1, 0)

	activeMemberships, err := s.subscriptionRepository.CountActiveAsOf(ctx, now)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count active memberships")
	}

	previousActiveMemberships, err := s.subscriptionRepository.CountActiveAsOf(ctx, previousMonth)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count previous active memberships")
	}

	retainedMemberships, err := s.subscriptionRepository.CountRetainedAsOf(ctx, previousMonth, now)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to count retained memberships")
	}

	retentionPct := 0.0
	if previousActiveMemberships > 0 {
		retentionPct = (float64(retainedMemberships) / float64(previousActiveMemberships)) * 100
	}

	return &dto.MemberRetentionResponse{
		ActiveMemberships: activeMemberships,
		RetentionPct:      math.Round(retentionPct*10) / 10,
	}, nil
}

func (s *subscriptionService) ensureStripeCustomer(ctx context.Context, user *model.User) (string, error) {
	if user.StripeCustomerID != nil && *user.StripeCustomerID != "" {
		return *user.StripeCustomerID, nil
	}

	params := &stripe.CustomerParams{}
	params.Email = stripe.String(user.Email)
	params.Name = stripe.String(user.Name)
	params.Metadata = map[string]string{"user_id": user.ID.String()}

	cust, err := customer.New(params)
	if err != nil {
		return "", apperror.NewInternalServerError("Failed to create Stripe customer")
	}

	user.StripeCustomerID = &cust.ID
	if _, err := s.userRepository.Updates(ctx, user); err != nil {
		return "", apperror.NewInternalServerError("Failed to persist Stripe customer")
	}

	return *user.StripeCustomerID, nil
}

func (s *subscriptionService) handleCheckoutSessionCompleted(ctx context.Context, checkoutSession *stripe.CheckoutSession) error {
	if checkoutSession == nil {
		return nil
	}

	if checkoutSession.Mode == stripe.CheckoutSessionModePayment {
		return s.handleCoursePurchaseCheckoutCompleted(ctx, checkoutSession)
	}

	if checkoutSession.Subscription == nil {
		return nil
	}

	stripeSub, err := stripesubscription.Get(checkoutSession.Subscription.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to fetch Stripe subscription")
	}

	return s.syncSubscription(ctx, checkoutSession.ID, stripeSub)
}

func (s *subscriptionService) handleCoursePurchaseCheckoutCompleted(ctx context.Context, checkoutSession *stripe.CheckoutSession) error {
	if checkoutSession == nil {
		return nil
	}

	userID := extractUserIDFromCheckoutSession(checkoutSession)

	purchase, err := s.coursePurchaseRepository.GetByCheckoutSessionID(ctx, checkoutSession.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get purchase")
	}

	if purchase == nil {
		if userID == uuid.Nil {
			return nil
		}
		purchase = &model.CoursePurchase{
			UserID:                  userID,
			StripeCheckoutSessionID: checkoutSession.ID,
		}
	}

	if userID == uuid.Nil {
		userID = purchase.UserID
	}

	// If the purchase is already paid, we don't need to update it
	wasPaid := purchase.Status == string(consts.CoursePurchaseStatusPaid)

	purchase.StripePaymentIntentID = extractCheckoutPaymentIntentID(checkoutSession)
	if purchase.Amount <= 0 {
		purchase.Amount = checkoutSession.AmountTotal
	}
	if strings.TrimSpace(purchase.Currency) == "" {
		purchase.Currency = normalizeCurrency(string(checkoutSession.Currency), s.stripeConfig.DefaultCurrency)
	}
	if purchase.StripeFee <= 0 && purchase.Amount > 0 {
		purchase.StripeFee = calculateStripeFeeFromAmount(purchase.Amount)
	}

	updatePurchaseStatusFromCheckoutPayment(purchase, checkoutSession)

	if purchase.ID == uuid.Nil {
		if _, err := s.coursePurchaseRepository.Create(ctx, purchase); err != nil {
			return apperror.NewInternalServerError("Failed to create purchase")
		}
	} else {
		if _, err := s.coursePurchaseRepository.Updates(ctx, purchase); err != nil {
			return apperror.NewInternalServerError("Failed to update purchase")
		}
	}

	if purchase.Status == string(consts.CoursePurchaseStatusPaid) && !wasPaid && purchase.CouponID != nil {
		if err := s.syncCouponCurrentRedemptions(ctx, *purchase.CouponID); err != nil {
			return err
		}
	}

	if purchase.Status == string(consts.CoursePurchaseStatusPaid) {
		if err := s.enrollUserToPurchasedCourses(ctx, userID, purchase.ID); err != nil {
			return err
		}
	}

	if purchase.Status == string(consts.CoursePurchaseStatusPaid) && !wasPaid {
		s.cancelOtherPendingCoursePurchases(ctx, userID, purchase.ID)
		if err := s.coursePurchaseRevenueShareRepository.RebuildByCoursePurchaseID(ctx, purchase.ID); err != nil {
			return apperror.NewInternalServerError("Failed to persist course purchase revenue shares")
		}
	}

	return nil
}

func (s *subscriptionService) syncCouponCurrentRedemptions(ctx context.Context, couponID uuid.UUID) error {
	couponData, err := s.courseCouponRepository.FindByID(ctx, couponID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to load coupon")
	}
	if couponData == nil || strings.TrimSpace(couponData.StripePromotionCodeID) == "" {
		return nil
	}

	promo, err := promotioncode.Get(couponData.StripePromotionCodeID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to sync coupon redemptions from Stripe")
	}

	updates := map[string]any{
		"current_redemptions": int64(promo.TimesRedeemed),
		"is_active":           promo.Active,
	}

	if _, err := s.courseCouponRepository.Update(ctx, couponID, updates); err != nil {
		return apperror.NewInternalServerError("Failed to update coupon redemptions")
	}

	return nil
}

func (s *subscriptionService) handleCoursePurchaseCheckoutFailed(ctx context.Context, checkoutSession *stripe.CheckoutSession) error {
	if checkoutSession == nil {
		return nil
	}

	purchase, err := s.coursePurchaseRepository.GetByCheckoutSessionID(ctx, checkoutSession.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get purchase")
	}
	if purchase == nil {
		return nil
	}

	purchase.StripePaymentIntentID = extractCheckoutPaymentIntentID(checkoutSession)

	purchase.Status = string(consts.CoursePurchaseStatusFailed)
	purchase.PurchasedAt = nil
	if purchase.Amount <= 0 {
		purchase.Amount = checkoutSession.AmountTotal
	}
	if strings.TrimSpace(purchase.Currency) == "" {
		purchase.Currency = normalizeCurrency(string(checkoutSession.Currency), s.stripeConfig.DefaultCurrency)
	}
	if purchase.StripeFee <= 0 && purchase.Amount > 0 {
		purchase.StripeFee = calculateStripeFeeFromAmount(purchase.Amount)
	}

	if _, err := s.coursePurchaseRepository.Updates(ctx, purchase); err != nil {
		return apperror.NewInternalServerError("Failed to update purchase")
	}

	return nil
}

func (s *subscriptionService) syncStripeCheckoutSession(ctx context.Context, checkoutSessionID string) error {
	checkoutSession, err := checkoutsession.Get(checkoutSessionID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to fetch Stripe checkout session")
	}
	if checkoutSession == nil {
		return nil
	}

	// Handle expired session
	if checkoutSession.Status == stripe.CheckoutSessionStatusExpired {
		return s.handleSubscriptionCheckoutExpired(ctx, checkoutSessionID)
	}

	if checkoutSession.Subscription == nil {
		return nil
	}

	stripeSub, err := stripesubscription.Get(checkoutSession.Subscription.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to fetch Stripe subscription")
	}

	return s.syncSubscription(ctx, checkoutSessionID, stripeSub)
}

func (s *subscriptionService) handleSubscriptionCheckoutExpired(ctx context.Context, checkoutSessionID string) error {
	subscription, err := s.subscriptionRepository.GetByCheckoutSessionID(ctx, checkoutSessionID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get subscription")
	}
	if subscription == nil {
		return nil
	}

	// Mark as expired if still incomplete
	if subscription.Status == string(stripe.SubscriptionStatusIncomplete) {
		updates := map[string]any{
			"status": string(stripe.SubscriptionStatusIncompleteExpired),
		}
		_, _ = s.subscriptionRepository.Update(ctx, subscription.ID, updates)
	}

	return nil
}

func (s *subscriptionService) syncCoursePurchaseCheckoutSession(ctx context.Context, checkoutSessionID string) error {
	checkoutSession, err := checkoutsession.Get(checkoutSessionID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to fetch Stripe checkout session")
	}
	if checkoutSession == nil || checkoutSession.Mode != stripe.CheckoutSessionModePayment {
		return nil
	}

	if checkoutSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		return s.handleCoursePurchaseCheckoutCompleted(ctx, checkoutSession)
	}

	if checkoutSession.Status == stripe.CheckoutSessionStatusExpired {
		return s.handleCoursePurchaseCheckoutFailed(ctx, checkoutSession)
	}

	return nil
}

func (s *subscriptionService) syncSubscription(ctx context.Context, checkoutSessionID string, stripeSub *stripe.Subscription) error {
	if stripeSub == nil {
		return nil
	}

	existing, err := s.subscriptionRepository.GetByStripeSubscriptionID(ctx, stripeSub.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to load subscription")
	}
	if existing == nil && strings.TrimSpace(checkoutSessionID) != "" {
		existing, err = s.subscriptionRepository.GetByCheckoutSessionID(ctx, checkoutSessionID, nil)
		if err != nil {
			return apperror.NewInternalServerError("Failed to load pending subscription")
		}
	}

	userID, planID, billingCycle := s.extractSubscriptionMetadata(ctx, stripeSub)
	customerID := ""
	if stripeSub.Customer != nil {
		customerID = stripeSub.Customer.ID
	}

	if existing == nil && userID != uuid.Nil {
		existing, err = s.findReusablePendingSubscription(ctx, userID, planID, billingCycle, customerID)
		if err != nil {
			return err
		}
	}

	if existing != nil {
		if userID == uuid.Nil {
			userID = existing.UserID
		}
		if planID == uuid.Nil {
			planID = existing.PlanID
		}
		if billingCycle == "" {
			billingCycle = existing.BillingCycle
		}
	}
	if userID == uuid.Nil {
		return nil
	}

	status := string(stripeSub.Status)
	startedAt := stripeTimestamp(stripeSub.Created)
	if startedAt == nil {
		now := time.Now().UTC()
		startedAt = &now
	}

	planName, planDescription, planPrice, planCurrency, planStripePriceID := s.resolveSubscriptionPlanSnapshot(ctx, stripeSub, planID)
	planName, planDescription, planPrice, planCurrency, planStripePriceID = s.mergeSubscriptionPlanSnapshotFallback(
		existing,
		planName,
		planDescription,
		planPrice,
		planCurrency,
		planStripePriceID,
	)

	var currentPeriodStart *time.Time
	var currentPeriodEnd *time.Time
	if stripeSub.Items != nil && len(stripeSub.Items.Data) > 0 {
		currentPeriodStart = stripeTimestamp(stripeSub.Items.Data[0].CurrentPeriodStart)
		currentPeriodEnd = stripeTimestamp(stripeSub.Items.Data[0].CurrentPeriodEnd)
	}

	if existing == nil {
		return s.createSubscriptionFromStripe(
			ctx,
			userID,
			planID,
			checkoutSessionID,
			customerID,
			billingCycle,
			status,
			startedAt,
			currentPeriodStart,
			currentPeriodEnd,
			stripeSub,
			planName,
			planDescription,
			planPrice,
			planCurrency,
			planStripePriceID,
		)
	}

	var stripeCanceledAt *time.Time
	if stripeSub.CanceledAt > 0 {
		stripeCanceledAt = stripeTimestamp(stripeSub.CanceledAt)
	}

	updates := buildSubscriptionUpdates(
		checkoutSessionID,
		planID,
		billingCycle,
		status,
		customerID,
		stripeSub.ID,
		currentPeriodStart,
		currentPeriodEnd,
		stripeSub.CancelAtPeriodEnd,
		planName,
		planDescription,
		planPrice,
		planCurrency,
		planStripePriceID,
		stripeCanceledAt,
	)

	if _, err := s.subscriptionRepository.Update(ctx, existing.ID, updates); err != nil {
		return apperror.NewInternalServerError("Failed to update subscription")
	}

	if status == string(stripe.SubscriptionStatusActive) || status == string(stripe.SubscriptionStatusTrialing) {
		s.cancelOtherPendingSubscriptions(ctx, userID, existing.ID)
	}

	// Payment rows are created from invoice.paid / invoice.payment_failed only.
	// Do not call syncLatestInvoice here: checkout.session.completed and invoice.paid
	// both run syncSubscription + syncPayment and would insert duplicate payments for the same invoice.

	if err := s.reconcileSubscriptionEnrollments(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) syncPayment(ctx context.Context, inv *stripe.Invoice, status string) error {
	if inv == nil {
		return nil
	}

	stripeSubscriptionID := ""
	if inv.Parent != nil && inv.Parent.SubscriptionDetails != nil && inv.Parent.SubscriptionDetails.Subscription != nil {
		stripeSubscriptionID = inv.Parent.SubscriptionDetails.Subscription.ID
	}
	if stripeSubscriptionID == "" {
		return nil
	}
	sub, err := s.subscriptionRepository.GetByStripeSubscriptionID(ctx, stripeSubscriptionID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get local subscription")
	}
	if sub == nil {
		return nil
	}

	payment, err := s.paymentRepository.GetByStripeInvoiceID(ctx, inv.ID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to load payment")
	}
	var persistedPayment *model.Payment

	periodStart, periodEnd := subscriptionBillingPeriodFromStripeAPI(stripeSubscriptionID)
	if periodStart == nil || periodEnd == nil {
		ps, pe := billingPeriodFromInvoiceLineItems(inv)
		if ps != nil && pe != nil {
			periodStart, periodEnd = ps, pe
		}
	}
	if periodStart == nil || periodEnd == nil {
		ps, pe := invoiceLevelBillingPeriodFallback(inv)
		if ps != nil && pe != nil {
			periodStart, periodEnd = ps, pe
		}
	}

	paidAt := stripeTimestamp(inv.StatusTransitions.PaidAt)
	billingPeriodStart := periodStart
	billingPeriodEnd := periodEnd
	paymentIntentID := extractInvoicePaymentIntentID(inv)
	failureReason := s.resolveInvoiceFailureReason(inv, paymentIntentID)
	amount := inv.AmountPaid
	if amount == 0 {
		amount = inv.AmountDue
	}
	stripeFee := calculateStripeFeeFromAmount(amount)

	if payment == nil {
		newPayment := &model.Payment{
			SubscriptionID:      sub.ID,
			StripeInvoiceID:     inv.ID,
			StripePaymentIntent: paymentIntentID,
			Status:              status,
			Amount:              amount,
			Currency:            strings.ToLower(string(inv.Currency)),
			StripeFee:           stripeFee,
			FailureReason:       failureReason,
			AttemptCount:        int64(inv.AttemptCount),
			BillingPeriodStart:  billingPeriodStart,
			BillingPeriodEnd:    billingPeriodEnd,
			PaidAt:              paidAt,
		}
		_, err := s.paymentRepository.Create(ctx, newPayment)
		if err != nil {
			if isDuplicateKeyError(err) {
				payment, err = s.paymentRepository.GetByStripeInvoiceID(ctx, inv.ID, nil)
				if err != nil {
					return apperror.NewInternalServerError("Failed to load payment")
				}
				if payment == nil {
					return apperror.NewInternalServerError("Failed to create payment")
				}
			} else {
				return apperror.NewInternalServerError("Failed to create payment")
			}
		} else {
			persistedPayment = newPayment
		}
	}
	if payment != nil && persistedPayment == nil {
		payment.Status = status
		payment.StripePaymentIntent = paymentIntentID
		payment.Amount = amount
		payment.Currency = strings.ToLower(string(inv.Currency))
		payment.StripeFee = stripeFee
		payment.FailureReason = failureReason
		payment.AttemptCount = int64(inv.AttemptCount)
		payment.BillingPeriodStart = billingPeriodStart
		payment.BillingPeriodEnd = billingPeriodEnd
		payment.PaidAt = paidAt
		if _, err := s.paymentRepository.Updates(ctx, payment); err != nil {
			return apperror.NewInternalServerError("Failed to update payment")
		}
		persistedPayment = payment
	}

	switch status {
	case string(consts.PaymentStatusSucceeded):
		sub.Status = string(stripe.SubscriptionStatusActive)
		if periodStart != nil && periodEnd != nil {
			sub.CurrentPeriodStart = periodStart
			sub.CurrentPeriodEnd = periodEnd
		}
		// Recovery from past_due / prior terminal markers — subscription is active again for this period.
		sub.CanceledAt = nil
		sub.EndedAt = nil
		s.cancelOtherPendingSubscriptions(ctx, sub.UserID, sub.ID)
	case string(consts.PaymentStatusFailed):
		sub.Status = string(stripe.SubscriptionStatusPastDue)
	}

	if _, err := s.subscriptionRepository.Updates(ctx, sub); err != nil {
		return apperror.NewInternalServerError("Failed to update subscription state")
	}
	if status == string(consts.PaymentStatusSucceeded) && persistedPayment != nil {
		if err := s.subscriptionRevenueShareRepository.RebuildByPaymentID(ctx, persistedPayment.ID); err != nil {
			return apperror.NewInternalServerError("Failed to persist subscription revenue shares")
		}
	}
	if err := s.reconcileSubscriptionEnrollments(ctx, sub.UserID); err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) cancelOtherPendingCoursePurchases(ctx context.Context, userID uuid.UUID, keepPurchaseID uuid.UUID) {
	if userID == uuid.Nil || keepPurchaseID == uuid.Nil {
		return
	}

	// Get courses from the successful purchase
	details, err := s.coursePurchaseDetailRepository.ListByPurchaseID(ctx, keepPurchaseID, nil)
	if err != nil {
		return
	}

	successCourseIDs := make(map[uuid.UUID]bool)
	for _, detail := range details {
		if detail != nil {
			successCourseIDs[detail.CourseID] = true
		}
	}
	if len(successCourseIDs) == 0 {
		return
	}

	// Get all purchases of the user with details
	purchases, err := s.coursePurchaseRepository.ListByUserID(ctx, userID, []repository.Preload{repository.Details})
	if err != nil {
		return
	}

	// Cancel only those with overlapping course IDs
	for _, purchase := range purchases {
		if purchase == nil || purchase.ID == keepPurchaseID {
			continue
		}
		if purchase.Status != string(consts.CoursePurchaseStatusPending) {
			continue
		}

		// Check if this purchase has any of the same courses
		hasOverlap := false
		if purchase.Details != nil {
			for _, detail := range purchase.Details {
				if detail != nil && successCourseIDs[detail.CourseID] {
					hasOverlap = true
					break
				}
			}
		}
		if !hasOverlap {
			continue
		}

		// Cancel this purchase
		updates := map[string]any{
			"status": string(consts.CoursePurchaseStatusFailed),
		}
		_, _ = s.coursePurchaseRepository.Update(ctx, purchase.ID, updates)

		checkoutSessionID := strings.TrimSpace(purchase.StripeCheckoutSessionID)
		if checkoutSessionID == "" {
			continue
		}

		_, _ = checkoutsession.Expire(checkoutSessionID, nil)
	}
}

func (s *subscriptionService) cancelOtherPendingSubscriptions(ctx context.Context, userID uuid.UUID, keepSubscriptionID uuid.UUID) {
	if userID == uuid.Nil {
		return
	}

	subscriptions, err := s.subscriptionRepository.ListByUserID(ctx, userID, nil)
	if err != nil {
		return
	}

	for _, sub := range subscriptions {
		if sub == nil || sub.ID == keepSubscriptionID {
			continue
		}
		if !isPendingSubscriptionStatus(sub.Status) {
			continue
		}

		updates := map[string]any{
			"status":               string(stripe.SubscriptionStatusCanceled),
			"canceled_at":          time.Now().UTC(),
			"ended_at":             time.Now().UTC(),
			"cancel_at_period_end": false,
		}
		_, _ = s.subscriptionRepository.Update(ctx, sub.ID, updates)

		stripeSubscriptionID := strings.TrimSpace(sub.StripeSubscriptionID)
		if stripeSubscriptionID != "" {
			_, _ = stripesubscription.Cancel(stripeSubscriptionID, nil)
			continue
		}

		checkoutSessionID := strings.TrimSpace(sub.StripeCheckoutSessionID)
		if checkoutSessionID == "" {
			continue
		}

		_, _ = checkoutsession.Expire(checkoutSessionID, nil)
	}
}

func isPendingSubscriptionStatus(status string) bool {
	switch status {
	case string(stripe.SubscriptionStatusIncomplete),
		string(stripe.SubscriptionStatusPastDue),
		string(stripe.SubscriptionStatusIncompleteExpired),
		string(stripe.SubscriptionStatusUnpaid):
		return true
	default:
		return false
	}
}

func (s *subscriptionService) extractSubscriptionMetadata(ctx context.Context, sub *stripe.Subscription) (uuid.UUID, uuid.UUID, string) {
	if sub == nil {
		return uuid.Nil, uuid.Nil, ""
	}

	userID := uuid.Nil
	planID := uuid.Nil
	billingCycle := ""

	metadata := sub.Metadata
	if metadata != nil {
		if uid, err := uuid.Parse(metadata["user_id"]); err == nil {
			userID = uid
		}
		if pid, err := uuid.Parse(metadata["plan_id"]); err == nil {
			planID = pid
		}
		billingCycle = metadata["billing_cycle"]
	}

	if planID == uuid.Nil {
		priceID := ""
		if len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil {
			priceID = sub.Items.Data[0].Price.ID
		}
		if priceID != "" {
			plan, err := s.planRepository.Find(ctx, "stripe_price_id = ?", nil, priceID)
			if err == nil && plan != nil {
				planID = plan.ID
				if billingCycle == "" {
					billingCycle = string(plan.BillingCycle)
				}
			}
		}
	}

	if billingCycle == "" && len(sub.Items.Data) > 0 && sub.Items.Data[0].Price != nil {
		if sub.Items.Data[0].Price.Recurring != nil {
			interval := string(sub.Items.Data[0].Price.Recurring.Interval)
			switch interval {
			case "month":
				billingCycle = string(consts.BillingCycleMonthly)
			case "year":
				billingCycle = string(consts.BillingCycleYearly)
			default:
				billingCycle = interval
			}
		}
	}

	if userID == uuid.Nil {
		if sub.Customer != nil {
			u, err := s.userRepository.Find(ctx, "stripe_customer_id = ?", nil, sub.Customer.ID)
			if err == nil && u != nil {
				userID = u.ID
			}
		}
	}

	return userID, planID, billingCycle
}

func (s *subscriptionService) resolvePlanPriceID(ctx context.Context, plan *model.Plan) (string, string, error) {
	if plan == nil {
		return "", "", apperror.NewBadRequestError("Invalid plan")
	}

	billingCycle := strings.TrimSpace(strings.ToLower(string(plan.BillingCycle)))
	if billingCycle != string(consts.BillingCycleMonthly) && billingCycle != string(consts.BillingCycleYearly) {
		return "", "", apperror.NewBadRequestError("Selected plan has invalid billing cycle")
	}

	interval := "month"
	if billingCycle == string(consts.BillingCycleYearly) {
		interval = "year"
	}

	currency := strings.ToLower(strings.TrimSpace(plan.Currency))
	if currency == "" {
		currency = strings.ToLower(strings.TrimSpace(s.stripeConfig.DefaultCurrency))
	}
	if currency == "" {
		currency = "usd"
	}

	amountCents := int64(math.Round(plan.Price * 100))
	if amountCents <= 0 {
		return "", "", apperror.NewBadRequestError("Selected plan has invalid price")
	}

	productID := strings.TrimSpace(plan.StripeProductID)
	if productID != "" {
		if _, err := product.Get(productID, nil); err != nil {
			if isStripeResourceMissing(err) {
				productID = ""
			} else {
				return "", "", apperror.NewInternalServerError("Failed to load Stripe plan product")
			}
		}
	}
	if productID == "" {
		createdProduct, err := product.New(&stripe.ProductParams{
			Name:        stripe.String(plan.Name),
			Description: stripe.String(plan.Description),
			Metadata: map[string]string{
				"plan_id":       plan.ID.String(),
				"billing_cycle": billingCycle,
				"catalog_type":  "subscription_plan",
			},
		})
		if err != nil {
			return "", "", apperror.NewInternalServerError("Failed to create Stripe plan product")
		}
		productID = createdProduct.ID
	}

	priceID := strings.TrimSpace(plan.StripePriceID)
	needNewPrice := priceID == ""
	if !needNewPrice {
		stripePriceObj, err := price.Get(priceID, nil)
		if err != nil {
			if isStripeResourceMissing(err) {
				needNewPrice = true
			} else {
				return "", "", apperror.NewInternalServerError("Failed to load Stripe plan price")
			}
		} else if stripePriceObj == nil || stripePriceObj.Recurring == nil ||
			stripePriceObj.Recurring.Interval != stripe.PriceRecurringInterval(interval) ||
			strings.ToLower(string(stripePriceObj.Currency)) != currency ||
			stripePriceObj.UnitAmount != amountCents {
			needNewPrice = true
		}
	}

	if needNewPrice {
		oldPriceID := priceID
		newPrice, err := price.New(&stripe.PriceParams{
			Currency:   stripe.String(currency),
			UnitAmount: stripe.Int64(amountCents),
			Product:    stripe.String(productID),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(interval),
			},
			Metadata: map[string]string{
				"plan_id":      plan.ID.String(),
				"catalog_type": "subscription_plan",
			},
			Active: stripe.Bool(plan.IsActive),
		})
		if err != nil {
			return "", "", apperror.NewInternalServerError("Failed to create Stripe plan price")
		}
		priceID = newPrice.ID
		if oldPriceID != "" {
			_, _ = price.Update(oldPriceID, &stripe.PriceParams{Active: stripe.Bool(false)})
		}
	}

	if plan.StripeProductID != productID || plan.StripePriceID != priceID || strings.ToLower(plan.Currency) != currency {
		updatedPlan, err := s.planRepository.Update(ctx, plan.ID, map[string]any{
			"stripe_product_id": productID,
			"stripe_price_id":   priceID,
			"currency":          currency,
		})
		if err != nil {
			return "", "", apperror.NewInternalServerError("Failed to persist Stripe plan catalog")
		}
		if updatedPlan != nil {
			plan.StripeProductID = updatedPlan.StripeProductID
			plan.StripePriceID = updatedPlan.StripePriceID
			plan.Currency = updatedPlan.Currency
		}
	}

	return priceID, billingCycle, nil
}

func (s *subscriptionService) resolveSubscriptionPlanSnapshot(ctx context.Context, sub *stripe.Subscription, planID uuid.UUID) (string, string, float64, string, string) {
	name := ""
	description := ""
	planPrice := float64(0)
	currency := ""
	stripePriceID := ""

	if planID != uuid.Nil {
		plan, err := s.planRepository.FindByID(ctx, planID, nil)
		if err == nil && plan != nil {
			name = plan.Name
			description = plan.Description
			planPrice = plan.Price
			currency = strings.ToLower(plan.Currency)
			stripePriceID = plan.StripePriceID
		}
	}

	if sub != nil && sub.Items != nil && len(sub.Items.Data) > 0 && sub.Items.Data[0] != nil && sub.Items.Data[0].Price != nil {
		priceObj := sub.Items.Data[0].Price
		if stripePriceID == "" {
			stripePriceID = priceObj.ID
		}
		if planPrice <= 0 {
			planPrice = float64(priceObj.UnitAmount) / 100
		}
		if currency == "" {
			currency = strings.ToLower(string(priceObj.Currency))
		}
		if name == "" && priceObj.Product != nil {
			name = priceObj.Product.Name
			description = priceObj.Product.Description
		}
	}

	return name, description, planPrice, currency, stripePriceID
}

func (s *subscriptionService) createCoursePurchasePending(
	ctx context.Context,
	userID uuid.UUID,
	couponID *uuid.UUID,
	courseIDs []uuid.UUID,
	checkoutSessionID string,
	amount int64,
	currency string,
	priceByCourseID map[uuid.UUID]int64,
) error {
	purchase, err := s.coursePurchaseRepository.Create(ctx, &model.CoursePurchase{
		UserID:                  userID,
		CouponID:                couponID,
		StripeCheckoutSessionID: checkoutSessionID,
		Amount:                  amount,
		Currency:                currency,
		StripeFee:               calculateStripeFeeFromAmount(amount),
		Status:                  string(consts.CoursePurchaseStatusPending),
	})
	if err != nil {
		return apperror.NewInternalServerError("Failed to create pending purchase")
	}

	details := make([]*model.CoursePurchaseDetail, 0, len(courseIDs))
	for _, id := range courseIDs {
		price := priceByCourseID[id]
		details = append(details, &model.CoursePurchaseDetail{
			CoursePurchaseID: purchase.ID,
			CourseID:         id,
			Price:            price,
			Currency:         currency,
		})
	}

	if err := s.coursePurchaseDetailRepository.CreateBatch(ctx, details); err != nil {
		return apperror.NewInternalServerError("Failed to create pending purchase details")
	}

	return nil
}

func (s *subscriptionService) ensureCourseStripeProductCatalog(ctx context.Context, course *model.Course) error {
	if course == nil {
		return apperror.NewBadRequestError("Invalid course")
	}

	amountCents := int64(math.Round(course.Price * 100))
	if amountCents <= 0 {
		return apperror.NewBadRequestError("Course price must be greater than zero")
	}

	currency := strings.ToLower(strings.TrimSpace(course.StripeCurrency))
	if currency == "" {
		currency = strings.ToLower(strings.TrimSpace(s.stripeConfig.DefaultCurrency))
	}
	if currency == "" {
		currency = "usd"
	}

	productID := strings.TrimSpace(course.StripeProductID)
	if productID != "" {
		if _, err := product.Get(productID, nil); err != nil {
			if isStripeResourceMissing(err) {
				productID = ""
			} else {
				return apperror.NewInternalServerError("Failed to load Stripe product")
			}
		}
	}
	if productID == "" {
		createdProduct, err := product.New(&stripe.ProductParams{
			Name:        stripe.String(course.Title),
			Description: stripe.String(course.Description),
			Metadata: map[string]string{
				"course_id":      course.ID.String(),
				"course_slug":    course.Slug,
				"catalog_type":   "lifetime_course",
				"price_amount":   strconv.FormatInt(amountCents, 10),
				"price_currency": currency,
			},
		})
		if err != nil {
			return apperror.NewInternalServerError("Failed to create Stripe product")
		}
		productID = createdProduct.ID
		course.StripeProductID = productID
	} else {
		_, _ = product.Update(productID, &stripe.ProductParams{
			Name:        stripe.String(course.Title),
			Description: stripe.String(course.Description),
		})
	}

	currentPriceID := strings.TrimSpace(course.StripePriceID)
	needNewPrice := currentPriceID == "" || course.StripeAmount != amountCents || strings.ToLower(course.StripeCurrency) != currency
	if !needNewPrice {
		if _, err := price.Get(currentPriceID, nil); err != nil {
			if isStripeResourceMissing(err) {
				needNewPrice = true
			} else {
				return apperror.NewInternalServerError("Failed to load Stripe price")
			}
		}
	}
	if needNewPrice {
		oldPriceID := currentPriceID
		if oldPriceID != "" {
			_, _ = price.Update(oldPriceID, &stripe.PriceParams{Active: stripe.Bool(false)})
		}

		createdPrice, err := price.New(&stripe.PriceParams{
			Currency:   stripe.String(currency),
			UnitAmount: stripe.Int64(amountCents),
			Product:    stripe.String(productID),
			Metadata: map[string]string{
				"course_id":    course.ID.String(),
				"catalog_type": "lifetime_course",
			},
		})
		if err != nil {
			return apperror.NewInternalServerError("Failed to create Stripe price")
		}

		course.StripePriceID = createdPrice.ID
	}

	now := time.Now().UTC()
	course.StripeCurrency = currency
	course.StripeAmount = amountCents
	course.StripeSyncedAt = &now

	if _, err := s.courseRepository.Updates(ctx, course); err != nil {
		return apperror.NewInternalServerError("Failed to update course Stripe catalog")
	}

	return nil
}

func extractInvoicePaymentIntentID(inv *stripe.Invoice) string {
	if inv == nil || inv.Payments == nil || len(inv.Payments.Data) == 0 {
		return ""
	}

	for _, invoicePayment := range inv.Payments.Data {
		if invoicePayment == nil || invoicePayment.Payment == nil {
			continue
		}
		if invoicePayment.Payment.PaymentIntent != nil {
			return invoicePayment.Payment.PaymentIntent.ID
		}
	}

	return ""
}

func (s *subscriptionService) resolveInvoiceFailureReason(inv *stripe.Invoice, paymentIntentID string) string {
	if inv != nil && inv.LastFinalizationError != nil && inv.LastFinalizationError.Msg != "" {
		return inv.LastFinalizationError.Msg
	}

	if paymentIntentID == "" {
		return ""
	}

	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil || pi == nil || pi.LastPaymentError == nil {
		return ""
	}

	return pi.LastPaymentError.Msg
}

func stripeTimestamp(v int64) *time.Time {
	if v <= 0 {
		return nil
	}
	t := time.Unix(v, 0).UTC()
	return &t
}

// subscriptionBillingPeriodFromStripeAPI returns the subscription's current billing window from
// SubscriptionItem.current_period_* — the canonical source. Do not use Invoice.period_start/end
// for this: Stripe documents those as the invoice "usage" window and they are often equal.
func subscriptionBillingPeriodFromStripeAPI(subscriptionID string) (start, end *time.Time) {
	id := strings.TrimSpace(subscriptionID)
	if id == "" {
		return nil, nil
	}
	stripeSub, err := stripesubscription.Get(id, nil)
	if err != nil || stripeSub == nil || stripeSub.Items == nil || len(stripeSub.Items.Data) == 0 {
		return nil, nil
	}
	item := stripeSub.Items.Data[0]
	if item == nil {
		return nil, nil
	}
	start = stripeTimestamp(item.CurrentPeriodStart)
	end = stripeTimestamp(item.CurrentPeriodEnd)
	if start == nil || end == nil || !end.After(*start) {
		return nil, nil
	}
	return start, end
}

// billingPeriodFromInvoiceLineItems uses each line's period (service interval for that price).
// Stripe recommends this over invoice-level period for subscription invoices.
func billingPeriodFromInvoiceLineItems(inv *stripe.Invoice) (start, end *time.Time) {
	if inv == nil || inv.Lines == nil {
		return nil, nil
	}
	var best *stripe.Period
	for _, line := range inv.Lines.Data {
		if line == nil || line.Period == nil {
			continue
		}
		p := line.Period
		if p.End <= p.Start {
			continue
		}
		if best == nil || (p.End-p.Start) > (best.End-best.Start) {
			best = p
		}
	}
	if best == nil {
		return nil, nil
	}
	return stripeTimestamp(best.Start), stripeTimestamp(best.End)
}

func invoiceLevelBillingPeriodFallback(inv *stripe.Invoice) (start, end *time.Time) {
	if inv == nil || inv.PeriodEnd <= inv.PeriodStart {
		return nil, nil
	}
	return stripeTimestamp(inv.PeriodStart), stripeTimestamp(inv.PeriodEnd)
}

func isStripeResourceMissing(err error) bool {
	stripeErr, ok := err.(*stripe.Error)
	if !ok || stripeErr == nil {
		return false
	}
	return stripeErr.Code == stripe.ErrorCodeResourceMissing
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "idx_payments_stripe_invoice_id_unique")
}

// calculateStripeFeeFromAmount estimates Stripe processing fee in the smallest currency unit.
// Formula: round(amount * 2.9%) + 30 cents.
func calculateStripeFeeFromAmount(amount int64) int64 {
	if amount <= 0 {
		return 0
	}

	fee := int64(math.Round(float64(amount)*0.029)) + 30
	if fee < 0 {
		return 0
	}

	return fee
}

func applyCouponDiscountAmount(amount int64, coupon *model.Coupon) int64 {
	if amount <= 0 || coupon == nil {
		return amount
	}

	switch strings.ToLower(strings.TrimSpace(coupon.DiscountType)) {
	case "percent":
		discount := int64(math.Round(float64(amount) * (float64(coupon.DiscountValue) / 100)))
		if discount <= 0 {
			return amount
		}
		if discount >= amount {
			return 0
		}
		return amount - discount
	case "amount":
		if coupon.DiscountValue <= 0 {
			return amount
		}
		if coupon.DiscountValue >= amount {
			return 0
		}
		return amount - coupon.DiscountValue
	default:
		return amount
	}
}

func decodeCheckoutSession(raw json.RawMessage) (stripe.CheckoutSession, error) {
	var checkoutSession stripe.CheckoutSession
	if err := json.Unmarshal(raw, &checkoutSession); err != nil {
		return stripe.CheckoutSession{}, apperror.NewBadRequestError("Invalid checkout session payload")
	}
	return checkoutSession, nil
}

func decodeStripeSubscription(raw json.RawMessage) (stripe.Subscription, error) {
	var sub stripe.Subscription
	if err := json.Unmarshal(raw, &sub); err != nil {
		return stripe.Subscription{}, apperror.NewBadRequestError("Invalid subscription payload")
	}
	return sub, nil
}

func decodeStripeInvoice(raw json.RawMessage) (stripe.Invoice, error) {
	var inv stripe.Invoice
	if err := json.Unmarshal(raw, &inv); err != nil {
		return stripe.Invoice{}, apperror.NewBadRequestError("Invalid invoice payload")
	}
	return inv, nil
}

func extractUserIDFromCheckoutSession(checkoutSession *stripe.CheckoutSession) uuid.UUID {
	if checkoutSession == nil || checkoutSession.Metadata == nil {
		return uuid.Nil
	}
	parsedUserID, err := uuid.Parse(checkoutSession.Metadata["user_id"])
	if err != nil {
		return uuid.Nil
	}
	return parsedUserID
}

func extractCheckoutPaymentIntentID(checkoutSession *stripe.CheckoutSession) string {
	if checkoutSession == nil || checkoutSession.PaymentIntent == nil {
		return ""
	}
	return checkoutSession.PaymentIntent.ID
}

func updatePurchaseStatusFromCheckoutPayment(purchase *model.CoursePurchase, checkoutSession *stripe.CheckoutSession) {
	if purchase == nil || checkoutSession == nil {
		return
	}

	if checkoutSession.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
		now := time.Now().UTC()
		purchase.Status = string(consts.CoursePurchaseStatusPaid)
		purchase.PurchasedAt = &now
		return
	}

	purchase.Status = string(consts.CoursePurchaseStatusPending)
}

func (s *subscriptionService) enrollUserToPurchasedCourses(ctx context.Context, userID, purchaseID uuid.UUID) error {
	details, err := s.coursePurchaseDetailRepository.ListByPurchaseID(ctx, purchaseID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to load purchase details")
	}

	for _, detail := range details {
		if detail == nil {
			continue
		}
		isEnrolled, err := s.enrollmentRepository.CheckEnrollment(ctx, userID, detail.CourseID)
		if err != nil {
			return apperror.NewInternalServerError("Failed to check enrollment")
		}
		if isEnrolled {
			continue
		}

		if _, err := s.enrollmentRepository.EnrollIfNotExists(ctx, userID, detail.CourseID, consts.EnrollmentTypeCoursePurchase); err != nil {
			return apperror.NewInternalServerError("Failed to enroll purchased course")
		}
		if _, err := s.courseRepository.Update(ctx, detail.CourseID, map[string]any{
			"total_student": gorm.Expr("total_student + ?", 1),
		}); err != nil {
			return apperror.NewInternalServerError("Failed to update course total students")
		}
	}

	return nil
}

func normalizeCurrency(currency string, defaultCurrency string) string {
	normalized := strings.ToLower(strings.TrimSpace(currency))
	if normalized != "" {
		return normalized
	}

	normalizedDefault := strings.ToLower(strings.TrimSpace(defaultCurrency))
	if normalizedDefault != "" {
		return normalizedDefault
	}

	return "usd"
}

func (s *subscriptionService) findReusablePendingSubscription(
	ctx context.Context,
	userID uuid.UUID,
	planID uuid.UUID,
	billingCycle string,
	customerID string,
) (*model.Subscription, error) {
	// Stripe can send customer.subscription.* before checkout.session.completed.
	// Reuse the latest local pending row instead of creating a duplicate subscription.
	userSubs, err := s.subscriptionRepository.ListByUserID(ctx, userID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to list user subscriptions")
	}

	for _, cand := range userSubs {
		if cand == nil {
			continue
		}
		if strings.TrimSpace(cand.StripeSubscriptionID) != "" {
			continue
		}
		if strings.TrimSpace(cand.Status) != string(stripe.SubscriptionStatusIncomplete) {
			continue
		}
		if planID != uuid.Nil && cand.PlanID != planID {
			continue
		}
		if billingCycle != "" && cand.BillingCycle != billingCycle {
			continue
		}
		if customerID != "" && strings.TrimSpace(cand.StripeCustomerID) != "" && cand.StripeCustomerID != customerID {
			continue
		}
		return cand, nil
	}

	return nil, nil
}

func (s *subscriptionService) mergeSubscriptionPlanSnapshotFallback(
	existing *model.Subscription,
	planName string,
	planDescription string,
	planPrice float64,
	planCurrency string,
	planStripePriceID string,
) (string, string, float64, string, string) {
	if existing != nil {
		if planName == "" {
			planName = existing.PlanName
		}
		if planDescription == "" {
			planDescription = existing.PlanDescription
		}
		if planPrice <= 0 {
			planPrice = existing.PlanPrice
		}
		if planCurrency == "" {
			planCurrency = existing.PlanCurrency
		}
		if planStripePriceID == "" {
			planStripePriceID = existing.PlanStripePriceID
		}
	}

	planCurrency = normalizeCurrency(planCurrency, s.stripeConfig.DefaultCurrency)
	return planName, planDescription, planPrice, planCurrency, planStripePriceID
}

func (s *subscriptionService) createSubscriptionFromStripe(
	ctx context.Context,
	userID uuid.UUID,
	planID uuid.UUID,
	checkoutSessionID string,
	customerID string,
	billingCycle string,
	status string,
	startedAt *time.Time,
	currentPeriodStart *time.Time,
	currentPeriodEnd *time.Time,
	stripeSub *stripe.Subscription,
	planName string,
	planDescription string,
	planPrice float64,
	planCurrency string,
	planStripePriceID string,
) error {
	newSub := &model.Subscription{
		UserID:                  userID,
		PlanID:                  planID,
		PlanName:                planName,
		PlanDescription:         planDescription,
		PlanPrice:               planPrice,
		PlanCurrency:            planCurrency,
		PlanStripePriceID:       planStripePriceID,
		StripeCheckoutSessionID: checkoutSessionID,
		StripeSubscriptionID:    stripeSub.ID,
		StripeCustomerID:        customerID,
		BillingCycle:            billingCycle,
		Status:                  status,
		CurrentPeriodStart:      currentPeriodStart,
		CurrentPeriodEnd:        currentPeriodEnd,
		CancelAtPeriodEnd:       stripeSub.CancelAtPeriodEnd,
		StartedAt:               *startedAt,
	}

	if status == string(stripe.SubscriptionStatusCanceled) && stripeSub.CanceledAt > 0 {
		canceledAt := *stripeTimestamp(stripeSub.CanceledAt)
		newSub.CanceledAt = &canceledAt
		newSub.EndedAt = currentPeriodEnd
	}

	_, err := s.subscriptionRepository.Create(ctx, newSub)
	if err != nil {
		return apperror.NewInternalServerError("Failed to create subscription")
	}
	if status == string(stripe.SubscriptionStatusActive) || status == string(stripe.SubscriptionStatusTrialing) {
		s.cancelOtherPendingSubscriptions(ctx, userID, newSub.ID)
	}
	if err := s.reconcileSubscriptionEnrollments(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *subscriptionService) reconcileSubscriptionEnrollments(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return nil
	}

	activeSub, err := s.subscriptionRepository.GetActiveByUserID(ctx, userID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to check active subscription")
	}
	if activeSub != nil {
		return nil
	}

	if err := s.enrollmentRepository.CancelActiveSubscriptionEnrollmentsByUser(ctx, userID); err != nil {
		return apperror.NewInternalServerError("Failed to cancel subscription enrollments")
	}
	return nil
}

func buildSubscriptionUpdates(
	checkoutSessionID string,
	planID uuid.UUID,
	billingCycle string,
	status string,
	customerID string,
	stripeSubscriptionID string,
	currentPeriodStart *time.Time,
	currentPeriodEnd *time.Time,
	cancelAtPeriodEnd bool,
	planName string,
	planDescription string,
	planPrice float64,
	planCurrency string,
	planStripePriceID string,
	stripeCanceledAt *time.Time,
) map[string]any {
	updates := map[string]any{
		"status":                 status,
		"current_period_start":   currentPeriodStart,
		"current_period_end":     currentPeriodEnd,
		"cancel_at_period_end":   cancelAtPeriodEnd,
		"stripe_customer_id":     customerID,
		"stripe_subscription_id": stripeSubscriptionID,
	}
	if strings.TrimSpace(checkoutSessionID) != "" {
		updates["stripe_checkout_session_id"] = checkoutSessionID
	}
	if billingCycle != "" {
		updates["billing_cycle"] = billingCycle
	}
	if planID != uuid.Nil {
		updates["plan_id"] = planID
	}
	if planName != "" {
		updates["plan_name"] = planName
	}
	if planDescription != "" {
		updates["plan_description"] = planDescription
	}
	if planPrice > 0 {
		updates["plan_price"] = planPrice
	}
	if planCurrency != "" {
		updates["plan_currency"] = planCurrency
	}
	if planStripePriceID != "" {
		updates["plan_stripe_price_id"] = planStripePriceID
	}

	// Active / trialing: subscription is entitled; clear end markers if Stripe moved state back.
	if status == string(stripe.SubscriptionStatusActive) || status == string(stripe.SubscriptionStatusTrialing) {
		updates["ended_at"] = nil
		updates["canceled_at"] = nil
	}

	if status == string(stripe.SubscriptionStatusCanceled) {
		if stripeCanceledAt != nil {
			updates["canceled_at"] = stripeCanceledAt
		} else {
			now := time.Now().UTC()
			updates["canceled_at"] = &now
		}
		// Access ends at period end (Stripe) when subscription is fully canceled.
		updates["ended_at"] = currentPeriodEnd
	}
	return updates
}
