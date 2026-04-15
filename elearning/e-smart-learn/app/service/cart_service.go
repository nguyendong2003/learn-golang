package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"errors"
	"fmt"
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
	"gorm.io/gorm"
)

type CartService interface {
	AddCourse(ctx context.Context, userID, courseID uuid.UUID) error
	RemoveCourse(ctx context.Context, userID, courseID uuid.UUID) error
	GetMyCart(ctx context.Context, userID uuid.UUID) (*dto.CartResponse, error)
	PreviewCheckout(ctx context.Context, userID uuid.UUID, request dto.CartCheckoutRequest) (*dto.CartCheckoutPreviewResponse, error)
	Checkout(ctx context.Context, userID uuid.UUID, request dto.CartCheckoutRequest) (*dto.CheckoutSessionResponse, error)
}

type cartService struct {
	cartRepository                 repository.CartRepository
	userRepository                 repository.UserRepository
	courseRepository               repository.CourseRepository
	couponRepository               repository.CouponRepository
	courseCouponRepository         repository.CourseCouponRepository
	coursePurchaseRepository       repository.CoursePurchaseRepository
	coursePurchaseDetailRepository repository.CoursePurchaseDetailRepository
	stripeConfig                   *config.StripeConfig
}

func NewCartService(
	cartRepository repository.CartRepository,
	userRepository repository.UserRepository,
	courseRepository repository.CourseRepository,
	couponRepository repository.CouponRepository,
	courseCouponRepository repository.CourseCouponRepository,
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
		courseCouponRepository:         courseCouponRepository,
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

func (s *cartService) PreviewCheckout(ctx context.Context, userID uuid.UUID, request dto.CartCheckoutRequest) (*dto.CartCheckoutPreviewResponse, error) {
	items, err := s.cartRepository.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get cart")
	}
	if len(items) == 0 {
		return nil, apperror.NewBadRequestError("Cart is empty")
	}

	enteredCouponCode := strings.TrimSpace(request.CouponCode)
	var enteredCoupon *model.Coupon
	if enteredCouponCode != "" {
		couponData, err := s.couponRepository.GetByCode(ctx, enteredCouponCode, nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate coupon")
		}
		if err := validateCouponAvailability(couponData); err != nil {
			return nil, err
		}
		enteredCoupon = couponData
	}

	remainingCouponQuota := make(map[uuid.UUID]int64)
	checkoutCurrency := ""
	enteredCouponAppliedCount := 0
	totalAmount := int64(0)

	response := &dto.CartCheckoutPreviewResponse{Items: make([]*dto.CartCheckoutPreviewItemResponse, 0, len(items))}

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

		courseCurrency := strings.ToLower(strings.TrimSpace(course.StripeCurrency))
		if courseCurrency == "" {
			courseCurrency = strings.ToLower(strings.TrimSpace(s.stripeConfig.DefaultCurrency))
		}
		if courseCurrency == "" {
			courseCurrency = "usd"
		}
		if checkoutCurrency == "" {
			checkoutCurrency = courseCurrency
		} else if checkoutCurrency != courseCurrency {
			return nil, apperror.NewBadRequestError("Cart contains courses with different currencies")
		}

		var appliedCoupon *model.Coupon

		enteredApplicableCoupon, err := s.resolveEnteredCouponForCourse(ctx, course.ID, enteredCoupon)
		if err != nil {
			return nil, err
		}

		if consumeCouponQuota(enteredApplicableCoupon, remainingCouponQuota) {
			appliedCoupon = enteredApplicableCoupon
			if enteredCoupon != nil && enteredApplicableCoupon != nil && enteredApplicableCoupon.ID == enteredCoupon.ID {
				enteredCouponAppliedCount++
			}
		} else {
			defaultCoupon, err := s.resolveDefaultCouponForCourse(ctx, course.ID)
			if err != nil {
				return nil, err
			}
			if consumeCouponQuota(defaultCoupon, remainingCouponQuota) {
				appliedCoupon = defaultCoupon
			}
		}

		finalAmount := applyCouponDiscountAmount(course.StripeAmount, appliedCoupon)
		totalAmount += finalAmount

		previewItem := &dto.CartCheckoutPreviewItemResponse{
			Course: dto.NewCourseDetailResponse(course),
		}
		if appliedCoupon != nil {
			previewItem.Coupon = dto.NewCartCheckoutCouponResponse(appliedCoupon)
		}

		response.Items = append(response.Items, previewItem)
	}

	if enteredCoupon != nil && enteredCouponAppliedCount == 0 {
		return nil, apperror.NewBadRequestError("Coupon is not applicable to any course in cart")
	}

	if checkoutCurrency == "" {
		checkoutCurrency = "usd"
	}
	response.TotalAmount = float64(totalAmount) / 100
	response.Currency = checkoutCurrency

	return response, nil
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
	checkoutCurrency := ""
	remainingCouponQuota := make(map[uuid.UUID]int64)

	enteredCouponCode := strings.TrimSpace(request.CouponCode)
	var enteredCoupon *model.Coupon
	if enteredCouponCode != "" {
		couponData, err := s.couponRepository.GetByCode(ctx, enteredCouponCode, nil)
		if err != nil {
			return nil, apperror.NewInternalServerError("Failed to validate coupon")
		}
		if err := validateCouponAvailability(couponData); err != nil {
			return nil, err
		}
		enteredCoupon = couponData
	}

	enteredCouponAppliedCount := 0

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

		courseCurrency := strings.ToLower(strings.TrimSpace(course.StripeCurrency))
		if courseCurrency == "" {
			courseCurrency = strings.ToLower(strings.TrimSpace(s.stripeConfig.DefaultCurrency))
		}
		if courseCurrency == "" {
			courseCurrency = "usd"
		}
		if checkoutCurrency == "" {
			checkoutCurrency = courseCurrency
		} else if checkoutCurrency != courseCurrency {
			return nil, apperror.NewBadRequestError("Cart contains courses with different currencies")
		}

		var appliedCoupon *model.Coupon

		enteredApplicableCoupon, err := s.resolveEnteredCouponForCourse(ctx, course.ID, enteredCoupon)
		if err != nil {
			return nil, err
		}

		if consumeCouponQuota(enteredApplicableCoupon, remainingCouponQuota) {
			appliedCoupon = enteredApplicableCoupon
			if enteredCoupon != nil && enteredApplicableCoupon != nil && enteredApplicableCoupon.ID == enteredCoupon.ID {
				enteredCouponAppliedCount++
			}
		} else {
			defaultCoupon, err := s.resolveDefaultCouponForCourse(ctx, course.ID)
			if err != nil {
				return nil, err
			}
			if consumeCouponQuota(defaultCoupon, remainingCouponQuota) {
				appliedCoupon = defaultCoupon
			}
		}

		finalAmount := applyCouponDiscountAmount(course.StripeAmount, appliedCoupon)
		totalAmount += finalAmount

		lineItemMetadata := map[string]string{
			"course_id":         course.ID.String(),
			"course_slug":       course.Slug,
			"price_original":    strconv.FormatInt(course.StripeAmount, 10),
			"price_final":       strconv.FormatInt(finalAmount, 10),
			"price_currency":    checkoutCurrency,
			"applied_coupon_id": "",
			"applied_coupon":    "",
		}
		if appliedCoupon != nil {
			lineItemMetadata["applied_coupon_id"] = appliedCoupon.ID.String()
			lineItemMetadata["applied_coupon"] = strings.TrimSpace(appliedCoupon.Code)
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency:   stripe.String(checkoutCurrency),
				UnitAmount: stripe.Int64(finalAmount),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name:        stripe.String(course.Title),
					Description: stripe.String(buildCartCoursePriceDescription(checkoutCurrency, course.StripeAmount, finalAmount, appliedCoupon)),
					Metadata:    lineItemMetadata,
				},
			},
			Quantity: stripe.Int64(1),
		})

		var couponID *uuid.UUID
		if appliedCoupon != nil {
			couponID = &appliedCoupon.ID
		}
		details = append(details, &model.CoursePurchaseDetail{
			CourseID:      course.ID,
			CouponID:      couponID,
			PriceOriginal: course.StripeAmount,
			PriceFinal:    finalAmount,
			Currency:      checkoutCurrency,
		})
	}

	if enteredCoupon != nil && enteredCouponAppliedCount == 0 {
		return nil, apperror.NewBadRequestError("Coupon is not applicable to any course in cart")
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

	session, err := checkoutsession.New(params)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe checkout session")
	}

	purchase := &model.CoursePurchase{
		UserID:                  userID,
		StripeCheckoutSessionID: session.ID,
		TotalAmount:             totalAmount,
		StripeFee:               calculateStripeFeeFromAmount(totalAmount),
		Currency:                checkoutCurrency,
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

func (s *cartService) resolveEnteredCouponForCourse(ctx context.Context, courseID uuid.UUID, enteredCoupon *model.Coupon) (*model.Coupon, error) {
	if enteredCoupon == nil {
		return nil, nil
	}

	courseCoupon, err := s.courseCouponRepository.GetByCourseAndCouponCode(ctx, courseID, enteredCoupon.Code, []repository.Preload{repository.Coupon})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewInternalServerError("Failed to validate course coupon")
	}
	if err != nil || courseCoupon == nil || courseCoupon.Coupon == nil {
		return nil, nil
	}

	if err := validateCouponAvailability(courseCoupon.Coupon); err != nil {
		return nil, nil
	}

	return courseCoupon.Coupon, nil
}

func (s *cartService) resolveDefaultCouponForCourse(ctx context.Context, courseID uuid.UUID) (*model.Coupon, error) {
	defaultCoupon, err := s.courseCouponRepository.GetDefaultByCourseID(ctx, courseID, []repository.Preload{repository.Coupon})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, apperror.NewInternalServerError("Failed to load default coupon")
	}
	if defaultCoupon == nil || defaultCoupon.Coupon == nil {
		return nil, nil
	}

	if err := validateCouponAvailability(defaultCoupon.Coupon); err != nil {
		return nil, nil
	}

	return defaultCoupon.Coupon, nil
}

func consumeCouponQuota(coupon *model.Coupon, remainingCouponQuota map[uuid.UUID]int64) bool {
	if coupon == nil {
		return false
	}
	if coupon.MaxRedemptions == nil || *coupon.MaxRedemptions <= 0 {
		return true
	}

	remaining, exists := remainingCouponQuota[coupon.ID]
	if !exists {
		remaining = *coupon.MaxRedemptions - coupon.CurrentRedemptions
		if remaining < 0 {
			remaining = 0
		}
	}

	if remaining <= 0 {
		remainingCouponQuota[coupon.ID] = 0
		return false
	}

	remainingCouponQuota[coupon.ID] = remaining - 1
	return true
}

func buildCartCoursePriceDescription(currency string, originalAmount, finalAmount int64, appliedCoupon *model.Coupon) string {
	original := formatStripeMinorAmount(currency, originalAmount)
	final := formatStripeMinorAmount(currency, finalAmount)
	coupon := formatAppliedCouponDisplay(currency, appliedCoupon)

	return fmt.Sprintf("Original: %s | Coupon: %s | Final: %s", original, coupon, final)
}

func formatStripeMinorAmount(currency string, amount int64) string {
	if amount < 0 {
		amount = 0
	}
	normalizedCurrency := strings.ToUpper(strings.TrimSpace(currency))
	if normalizedCurrency == "" {
		normalizedCurrency = "USD"
	}

	var currencySymbols = map[string]string{
		"USD": "$",
		"EUR": "€",
		"VND": "₫",
		"JPY": "¥",
	}

	symbol := currencySymbols[normalizedCurrency]
	major := amount / 100
	minor := amount % 100
	return fmt.Sprintf("%s%d.%02d", symbol, major, minor)
}

func formatAppliedCouponDisplay(currency string, coupon *model.Coupon) string {
	if coupon == nil {
		return "None"
	}

	code := strings.TrimSpace(coupon.Code)
	if code == "" {
		code = "N/A"
	}

	switch strings.TrimSpace(coupon.DiscountType) {
	case consts.DiscountTypePercent:
		return fmt.Sprintf("%s (-%d%%)", code, coupon.DiscountValue)
	case consts.DiscountTypeAmount:
		return fmt.Sprintf("%s (-%s)", code, formatStripeMinorAmount(currency, coupon.DiscountValue))
	default:
		return code
	}
}
