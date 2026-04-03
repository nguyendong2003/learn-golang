package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/repository"
	"elearning-api/util"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/coupon"
	"github.com/stripe/stripe-go/v83/promotioncode"
)

type CouponService interface {
	Create(ctx context.Context, request dto.CreateCouponRequest) (*dto.CouponResponse, error)
	GetByID(ctx context.Context, couponID uuid.UUID) (*dto.CouponResponse, error)
	Deactivate(ctx context.Context, couponID uuid.UUID) error
	GetList(ctx context.Context, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error)
	GetAvailableCoupons(ctx context.Context, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error)
}

type couponService struct {
	couponRepository repository.CouponRepository
	stripeConfig     *config.StripeConfig
}

func NewCouponService(couponRepository repository.CouponRepository, stripeConfig *config.StripeConfig) CouponService {
	stripe.Key = stripeConfig.SecretKey

	return &couponService{
		couponRepository: couponRepository,
		stripeConfig:     stripeConfig,
	}
}

func (s *couponService) Create(ctx context.Context, request dto.CreateCouponRequest) (*dto.CouponResponse, error) {
	discountType := strings.ToLower(strings.TrimSpace(request.DiscountType))
	if discountType != "percent" && discountType != "amount" {
		return nil, apperror.NewBadRequestError("discount_type must be either percent or amount")
	}
	if request.DiscountValue <= 0 {
		return nil, apperror.NewBadRequestError("discount_value must be greater than zero")
	}
	if discountType == "percent" && request.DiscountValue > 100 {
		return nil, apperror.NewBadRequestError("discount_value for percent must be between 1 and 100")
	}

	if request.ExpiresAt != nil && request.ExpiresAt.Before(time.Now().UTC()) {
		return nil, apperror.NewBadRequestError("expires_at must be in the future")
	}

	code := strings.ToUpper(strings.TrimSpace(request.Code))
	existingCoupon, err := s.couponRepository.GetByCode(ctx, code, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to validate coupon code")
	}
	if existingCoupon != nil {
		return nil, apperror.NewBadRequestError("Coupon code already exists")
	}

	couponParams := &stripe.CouponParams{
		Duration: stripe.String(string(stripe.CouponDurationOnce)),
	}

	discountValue := request.DiscountValue
	currency := strings.ToLower(s.stripeConfig.DefaultCurrency)
	if currency == "" {
		currency = "usd"
	}

	if discountType == "percent" {
		couponParams.PercentOff = stripe.Float64(float64(discountValue))
	} else {
		couponParams.AmountOff = stripe.Int64(discountValue)
		if request.Currency != "" {
			currency = strings.ToLower(strings.TrimSpace(request.Currency))
		}
		couponParams.Currency = stripe.String(currency)
	}

	if request.MaxRedemptions != nil {
		couponParams.MaxRedemptions = stripe.Int64(*request.MaxRedemptions)
	}
	if request.ExpiresAt != nil {
		expiresAt := request.ExpiresAt.UTC().Unix()
		couponParams.RedeemBy = stripe.Int64(expiresAt)
	}

	stripeCoupon, err := coupon.New(couponParams)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to create Stripe coupon")
	}

	promotionParams := &stripe.PromotionCodeParams{
		Code: stripe.String(code),
		Promotion: &stripe.PromotionCodePromotionParams{
			Type:   stripe.String(string(stripe.PromotionCodePromotionTypeCoupon)),
			Coupon: stripe.String(stripeCoupon.ID),
		},
	}
	if request.MaxRedemptions != nil {
		promotionParams.MaxRedemptions = stripe.Int64(*request.MaxRedemptions)
	}
	if request.ExpiresAt != nil {
		expiresAt := request.ExpiresAt.UTC().Unix()
		promotionParams.ExpiresAt = stripe.Int64(expiresAt)
	}

	stripePromotionCode, err := promotioncode.New(promotionParams)
	if err != nil {
		if stripeErr, ok := err.(*stripe.Error); ok {
			if stripeErr.Type == stripe.ErrorTypeInvalidRequest {
				if strings.Contains(strings.ToLower(stripeErr.Msg), "already") || strings.Contains(strings.ToLower(stripeErr.Msg), "exists") {
					return nil, apperror.NewBadRequestError("Coupon code already exists on Stripe, please use another code")
				}
				if stripeErr.Msg != "" {
					return nil, apperror.NewBadRequestError(stripeErr.Msg)
				}
				return nil, apperror.NewBadRequestError("Invalid Stripe promotion code request")
			}
		}

		return nil, apperror.NewInternalServerError("Failed to create Stripe promotion code")
	}

	maxRedemptions := int64(0)
	if request.MaxRedemptions != nil {
		maxRedemptions = *request.MaxRedemptions
	}

	newCoupon := &model.Coupon{
		Code:                  code,
		StripeCouponID:        stripeCoupon.ID,
		StripePromotionCodeID: stripePromotionCode.ID,
		DiscountType:          discountType,
		DiscountValue:         discountValue,
		Currency:              currency,
		MaxRedemptions:        maxRedemptions,
		CurrentRedemptions:    int64(stripePromotionCode.TimesRedeemed),
		IsActive:              stripePromotionCode.Active,
		ExpiresAt:             request.ExpiresAt,
	}

	created, err := s.couponRepository.Create(ctx, newCoupon)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to persist coupon")
	}

	return dto.NewCouponResponse(created), nil
}

func (s *couponService) GetByID(ctx context.Context, couponID uuid.UUID) (*dto.CouponResponse, error) {
	couponData, err := s.couponRepository.FindByID(ctx, couponID, nil)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get coupon")
	}
	if couponData == nil {
		return nil, apperror.NewNotFoundError("Coupon not found")
	}

	return dto.NewCouponResponse(couponData), nil
}

func (s *couponService) Deactivate(ctx context.Context, couponID uuid.UUID) error {
	couponData, err := s.couponRepository.FindByID(ctx, couponID, nil)
	if err != nil {
		return apperror.NewInternalServerError("Failed to get coupon")
	}
	if couponData == nil {
		return apperror.NewNotFoundError("Coupon not found")
	}

	if !couponData.IsActive {
		return apperror.NewBadRequestError("Coupon is already inactive")
	}

	if couponData.IsActive {
		_, err = promotioncode.Update(couponData.StripePromotionCodeID, &stripe.PromotionCodeParams{Active: stripe.Bool(false)})
		if err != nil {
			return apperror.NewInternalServerError("Failed to deactivate Stripe promotion code")
		}
	}

	if _, err := s.couponRepository.Update(ctx, couponData.ID, map[string]any{"is_active": false}); err != nil {
		return apperror.NewInternalServerError("Failed to update coupon")
	}

	return nil
}

func (s *couponService) GetList(ctx context.Context, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error) {
	// Build sort query
	orderQuery := buildCouponSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildCouponQuery(request)

	categories, total, err := s.couponRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		nil,
		args...,
	)

	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve coupons")
	}

	return dto.NewListCouponResponse(categories), total, nil
}

func buildCouponSortQuery(sortBy string, sortOrder string) string {
	defaultSort := "created_at DESC"

	if sortBy == "" {
		return defaultSort
	}

	if sortOrder == "" {
		sortOrder = "DESC"
	}

	allowedSort := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"expires_at": true,
	}

	if !allowedSort[sortBy] {
		return defaultSort
	}

	return sortBy + " " + strings.ToUpper(sortOrder)
}

func buildCouponQuery(request dto.ListCouponRequest) (string, []any) {
	var conditions []string
	var args []any

	filters := map[string]*string{
		"code": request.Code,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}

func (s *couponService) GetAvailableCoupons(ctx context.Context, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error) {
	orderQuery := buildCouponSortQuery(request.SortBy, request.SortOrder)

	coupons, total, err := s.couponRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		"is_active = ?",
		nil,
		true)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve available coupons")
	}

	var response []*dto.CouponResponse
	for _, c := range coupons {
		response = append(response, dto.NewCouponResponse(c))
	}

	return response, total, nil
}
