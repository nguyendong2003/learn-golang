package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
	checkoutsession "github.com/stripe/stripe-go/v83/checkout/session"
	"github.com/stripe/stripe-go/v83/customer"
	"github.com/stripe/stripe-go/v83/price"
	"github.com/stripe/stripe-go/v83/product"
	"github.com/stripe/stripe-go/v83/promotioncode"
)

type CartService interface {
	AddCourse(ctx context.Context, userID, courseID uuid.UUID) error
	RemoveCourse(ctx context.Context, userID, courseID uuid.UUID) error
	GetMyCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error)
	Checkout(ctx context.Context, userID uuid.UUID, request dto.CartCheckoutRequest) (*dto.CheckoutSessionResponse, error)
}

type cartService struct {
	cartRepository                 repository.CartRepository
	userRepository                 repository.UserRepository
	courseRepository               repository.CourseRepository
	couponRepository               repository.CouponRepository
	coursePurchaseRepository       repository.CoursePurchaseRepository
	coursePurchaseDetailRepository repository.CoursePurchaseDetailRepository
	stripeConfig                   *config.StripeConfig
}

func NewCartService(
	cartRepository repository.CartRepository,
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	couponRepository repository.CouponRepository,
	coursePurchaseRepository repository.CoursePurchaseRepository,
	coursePurchaseDetailRepository repository.CoursePurchaseDetailRepository,
	stripeConfig *config.StripeConfig,
) CartService {
	stripe.Key = stripeConfig.SecretKey

	return &cartService{
		cartRepository:                 cartRepository,
		userRepository:                 userRepository,
		courseRepository:               courseRepository,
		couponRepository:               couponRepository,
		coursePurchaseRepository:       coursePurchaseRepository,
		coursePurchaseDetailRepository: coursePurchaseDetailRepository,
		stripeConfig:                   stripeConfig,
	}
}

func (s *cartService) AddCourse(ctx context.Context, userID, courseID uuid.UUID) error {
	course, err := s.courseRepository.FindByID(ctx, courseID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get course")
	}
	if course == nil {
		return apperror.NewNotFoundError("Course not found")
	}
	if course.Status != consts.CoursePublished {
		return apperror.NewBadRequestError("Course is not available for purchase")
	}
	if course.UserID == userID {
		return apperror.NewBadRequestError("Cannot add your own course")
	}
	if course.Price <= 0 {
		return apperror.NewBadRequestError("Course is free")
	}

	isPaid, err := s.coursePurchaseRepository.ExistsPaidByUserAndCourse(ctx, userID, courseID)
	if err != nil {
		return apperror.NewInternalServerError("Failed to validate purchase history")
	}
	if isPaid {
		return apperror.NewBadRequestError("Course already purchased")
	}

	isInCart, err := s.cartRepository.Exists(ctx, userID, courseID)
	if err != nil {
		return apperror.NewInternalServerError("Failed to check cart item")
	}
	if isInCart {
		return apperror.NewBadRequestError("Course already exists in cart")
	}

	if err := s.cartRepository.AddCourse(ctx, userID, courseID); err != nil {
		return apperror.NewInternalServerError("Failed to add course to cart")
	}
	return nil
}

func (s *cartService) RemoveCourse(ctx context.Context, userID, courseID uuid.UUID) error {
	isInCart, err := s.cartRepository.Exists(ctx, userID, courseID)
	if err != nil {
		return apperror.NewInternalServerError("Failed to check cart item")
	}
	if !isInCart {
		return apperror.NewBadRequestError("Course is not in cart")
	}

	if err := s.cartRepository.RemoveCourse(ctx, userID, courseID); err != nil {
		return apperror.NewInternalServerError("Failed to remove course from cart")
	}
	return nil
}

func (s *cartService) GetMyCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error) {
	items, err := s.cartRepository.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get cart")
	}

	for _, item := range items {
		if item == nil || item.Course == nil {
			continue
		}
		if err := s.ensureCourseStripeProductCatalog(ctx, item.Course); err != nil {
			return nil, err
		}
	}

	return dto.NewCartResponse(items), nil
}

func (s *cartService) Checkout(ctx context.Context, userID uuid.UUID, request dto.CartCheckoutRequest) (*dto.CheckoutSessionResponse, error) {
	items, err := s.cartRepository.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get cart")
	}
	if len(items) == 0 {
		return nil, apperror.NewBadRequestError("Cart is empty")
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

	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(items))
	details := make([]*model.CoursePurchaseDetail, 0, len(items))
	totalAmount := int64(0)
	currency := strings.ToLower(s.stripeConfig.DefaultCurrency)
	if currency == "" {
		currency = "usd"
	}

	for _, item := range items {
		if item == nil || item.Course == nil {
			continue
		}

		course := item.Course
		if course.Status != consts.CoursePublished {
			return nil, apperror.NewBadRequestError("Cart contains unpublished course")
		}
		if course.UserID == userID {
			return nil, apperror.NewBadRequestError("Cart contains your own course")
		}
		if course.Price <= 0 {
			return nil, apperror.NewBadRequestError("Cart contains free course")
		}

		isPaid, err := s.coursePurchaseRepository.ExistsPaidByUserAndCourse(ctx, userID, course.ID)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate purchase history")
		}
		if isPaid {
			return nil, apperror.NewBadRequestError("Cart contains already purchased course")
		}

		if err := s.ensureCourseStripeProductCatalog(ctx, course); err != nil {
			return nil, err
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(course.StripePriceID),
			Quantity: stripe.Int64(1),
		})

		totalAmount += course.StripeAmount
		if course.StripeCurrency != "" {
			currency = strings.ToLower(course.StripeCurrency)
		}

		details = append(details, &model.CoursePurchaseDetail{
			CourseID: course.ID,
			Price:    course.StripeAmount,
			Currency: currency,
		})
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(s.stripeConfig.SuccessURL),
		CancelURL:  stripe.String(s.stripeConfig.CancelURL),
		Customer:   stripe.String(customerID),
		LineItems:  lineItems,
		Metadata: map[string]string{
			"purchase_type": "cart_purchase",
			"user_id":       userID.String(),
		},
	}

	var appliedCouponID *uuid.UUID

	if request.CouponCode != "" {
		couponCode := strings.TrimSpace(request.CouponCode)
		couponData, err := s.couponRepository.GetByCode(ctx, couponCode, nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate coupon")
		}
		if couponData == nil || !couponData.IsActive {
			return nil, apperror.NewBadRequestError("Coupon code is invalid or inactive")
		}

		promotion, err := promotioncode.Get(couponData.StripePromotionCodeID, nil)
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
			{PromotionCode: stripe.String(couponData.StripePromotionCodeID)},
		}
		appliedCouponID = &couponData.ID
	}

	session, err := checkoutsession.New(params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe checkout session")
	}

	purchase := &model.CoursePurchase{
		UserID:                  userID,
		CouponID:                appliedCouponID,
		StripeCheckoutSessionID: session.ID,
		Amount:                  totalAmount,
		Currency:                currency,
		Status:                  string(consts.CoursePurchaseStatusPending),
	}
	createdPurchase, err := s.coursePurchaseRepository.Create(ctx, purchase)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create purchase")
	}

	for _, detail := range details {
		detail.CoursePurchaseID = createdPurchase.ID
	}
	if err := s.coursePurchaseDetailRepository.CreateBatch(ctx, details); err != nil {
		return nil, apperror.NewInternalServerError("Failed to create purchase details")
	}

	return &dto.CheckoutSessionResponse{SessionID: session.ID, URL: session.URL}, nil
}

func (s *cartService) ensureStripeCustomer(ctx context.Context, user *model.User) (string, error) {
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

func (s *cartService) ensureCourseStripeProductCatalog(ctx context.Context, course *model.Course) error {
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
			if isStripeResourceMissingError(err) {
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
			if isStripeResourceMissingError(err) {
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

func isStripeResourceMissingError(err error) bool {
	stripeErr, ok := err.(*stripe.Error)
	if !ok || stripeErr == nil {
		return false
	}
	return stripeErr.Code == stripe.ErrorCodeResourceMissing
}
