package service

import (
	"context"
	"elearning-api/apperror"
	"elearning-api/config"
	"elearning-api/consts"
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
	Create(ctx context.Context, userID uuid.UUID, request dto.CreateCouponRequest) (*dto.CouponResponse, error)
	GetByID(ctx context.Context, userID, couponID uuid.UUID) (*dto.CouponResponse, error)
	Deactivate(ctx context.Context, userID, couponID uuid.UUID) error
	GetList(ctx context.Context, userID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error)
	GetAssignableList(ctx context.Context, userID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error)
	GetAssignableListForCourse(ctx context.Context, userID, courseID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponseForCourse, int64, error)
}

type couponService struct {
	couponRepository       repository.CouponRepository
	courseRepository       repository.CourseRepository
	courseCouponRepository repository.CourseCouponRepository
	stripeConfig           *config.StripeConfig
}

func NewCouponService(
	couponRepository repository.CouponRepository,
	courseRepository repository.CourseRepository,
	courseCouponRepository repository.CourseCouponRepository,
	stripeConfig *config.StripeConfig,
) CouponService {
	stripe.Key = stripeConfig.SecretKey

	return &couponService{
		couponRepository:       couponRepository,
		courseRepository:       courseRepository,
		courseCouponRepository: courseCouponRepository,
		stripeConfig:           stripeConfig,
	}
}

func (s *couponService) Create(ctx context.Context, userID uuid.UUID, request dto.CreateCouponRequest) (*dto.CouponResponse, error) {
	if userID == uuid.Nil {
		return nil, apperror.NewUnauthorizedError("User not found")
	}

	discountType := strings.ToLower(strings.TrimSpace(request.DiscountType))
	if discountType != string(consts.DiscountTypePercent) && discountType != string(consts.DiscountTypeAmount) {
		return nil, apperror.NewBadRequestError("discount type must be either percent or amount")
	}
	if request.DiscountValue <= 0 {
		return nil, apperror.NewBadRequestError("discount value must be greater than zero")
	}
	if discountType == string(consts.DiscountTypePercent) && request.DiscountValue > 100 {
		return nil, apperror.NewBadRequestError("discount value for percent must be between 1 and 100")
	}

	if request.ExpiresAt != nil && request.ExpiresAt.Before(time.Now().UTC()) {
		return nil, apperror.NewBadRequestError("expires at must be in the future")
	}

	code := strings.TrimSpace(request.Code)
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

	if discountType == string(consts.DiscountTypePercent) {
		couponParams.PercentOff = stripe.Float64(float64(discountValue))
	} else {
		couponParams.AmountOff = stripe.Int64(discountValue)
		if request.Currency != nil {
			currency = strings.ToLower(strings.TrimSpace(*request.Currency))
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

	newCoupon := &model.Coupon{
		UserID:                &userID,
		Code:                  code,
		StripeCouponID:        stripeCoupon.ID,
		StripePromotionCodeID: stripePromotionCode.ID,
		DiscountType:          discountType,
		DiscountValue:         discountValue,
		Currency:              currency,
		MaxRedemptions:        request.MaxRedemptions,
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

func (s *couponService) GetByID(ctx context.Context, userID, couponID uuid.UUID) (*dto.CouponResponse, error) {
	if userID == uuid.Nil {
		return nil, apperror.NewUnauthorizedError("User not found")
	}

	couponData, err := s.couponRepository.Find(ctx, "id = ? AND user_id = ?", nil, couponID, userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("Failed to get coupon")
	}
	if couponData == nil {
		return nil, apperror.NewNotFoundError("Coupon not found")
	}

	return dto.NewCouponResponse(couponData), nil
}

func (s *couponService) Deactivate(ctx context.Context, userID, couponID uuid.UUID) error {
	if userID == uuid.Nil {
		return apperror.NewUnauthorizedError("User not found")
	}

	couponData, err := s.couponRepository.Find(ctx, "id = ? AND user_id = ?", nil, couponID, userID)
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

func (s *couponService) GetList(ctx context.Context, userID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, apperror.NewUnauthorizedError("User not found")
	}

	// Build sort query
	orderQuery := buildCouponSortQuery(request.SortBy, request.SortOrder)

	// Build filter query
	query, args := buildCouponQuery(userID, request)

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

func (s *couponService) GetAssignableList(ctx context.Context, userID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponse, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, apperror.NewUnauthorizedError("User not found")
	}

	orderQuery := buildCouponSortQuery(request.SortBy, request.SortOrder)
	query, args := buildAssignableCouponQuery(userID, request)

	coupons, total, err := s.couponRepository.List(
		ctx,
		request.Limit,
		request.Offset,
		orderQuery,
		query,
		nil,
		args...,
	)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve assignable coupons")
	}

	return dto.NewListCouponResponse(coupons), total, nil
}

func (s *couponService) GetAssignableListForCourse(ctx context.Context, userID, courseID uuid.UUID, request dto.ListCouponRequest) ([]*dto.CouponResponseForCourse, int64, error) {
	if userID == uuid.Nil {
		return nil, 0, apperror.NewUnauthorizedError("User not found")
	}

	if courseID == uuid.Nil {
		return nil, 0, apperror.NewBadRequestError("Course ID is required")
	}

	courseData, err := s.courseRepository.Find(ctx, "id = ? AND user_id = ?", nil, courseID, userID)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to validate course")
	}
	if courseData == nil {
		return nil, 0, apperror.NewNotFoundError("Course not found or you do not have permission to access this course")
	}

	courseCoupons, err := s.courseCouponRepository.ListByCourseID(ctx, courseID, nil)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to retrieve course coupons")
	}

	assignedMap := make(map[uuid.UUID]bool, len(courseCoupons))
	defaultMap := make(map[uuid.UUID]bool, len(courseCoupons))
	assignedIDs := make([]uuid.UUID, 0, len(courseCoupons))
	for _, cc := range courseCoupons {
		if !assignedMap[cc.CouponID] {
			assignedIDs = append(assignedIDs, cc.CouponID)
		}
		assignedMap[cc.CouponID] = true
		if cc.IsDefault {
			defaultMap[cc.CouponID] = true
		}
	}

	baseQuery, baseArgs := buildAssignableCouponQuery(userID, request)
	orderQuery := buildCouponSortQuery(request.SortBy, request.SortOrder)

	assignedQuery := baseQuery + " AND 1 = 0"
	assignedArgs := append([]any{}, baseArgs...)
	if len(assignedIDs) > 0 {
		assignedQuery = baseQuery + " AND id IN ?"
		assignedArgs = append(assignedArgs, assignedIDs)
	}

	unassignedQuery := baseQuery
	unassignedArgs := append([]any{}, baseArgs...)
	if len(assignedIDs) > 0 {
		unassignedQuery = baseQuery + " AND id NOT IN ?"
		unassignedArgs = append(unassignedArgs, assignedIDs)
	}

	assignedTotal, err := s.couponRepository.Count(ctx, assignedQuery, assignedArgs...)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to count assigned coupons")
	}

	unassignedTotal, err := s.couponRepository.Count(ctx, unassignedQuery, unassignedArgs...)
	if err != nil {
		return nil, 0, apperror.NewInternalServerError("Failed to count unassigned coupons")
	}

	total := assignedTotal + unassignedTotal
	if total == 0 {
		return []*dto.CouponResponseForCourse{}, 0, nil
	}

	offset := request.Offset
	remaining := request.Limit
	resultCoupons := make([]*model.Coupon, 0, request.Limit)

	if remaining > 0 && offset < int(assignedTotal) {
		assignedOffset := offset
		assignedLimit := remaining
		availableAssigned := int(assignedTotal) - assignedOffset
		if assignedLimit > availableAssigned {
			assignedLimit = availableAssigned
		}

		assignedCoupons, _, err := s.couponRepository.List(
			ctx,
			assignedLimit,
			assignedOffset,
			orderQuery,
			assignedQuery,
			nil,
			assignedArgs...,
		)
		if err != nil {
			return nil, 0, apperror.NewInternalServerError("Failed to retrieve assigned coupons")
		}

		resultCoupons = append(resultCoupons, assignedCoupons...)
		remaining -= len(assignedCoupons)
		offset = 0
	} else {
		offset -= int(assignedTotal)
		if offset < 0 {
			offset = 0
		}
	}

	if remaining > 0 {
		unassignedCoupons, _, err := s.couponRepository.List(
			ctx,
			remaining,
			offset,
			orderQuery,
			unassignedQuery,
			nil,
			unassignedArgs...,
		)
		if err != nil {
			return nil, 0, apperror.NewInternalServerError("Failed to retrieve unassigned coupons")
		}

		resultCoupons = append(resultCoupons, unassignedCoupons...)
	}

	res := make([]*dto.CouponResponseForCourse, 0, len(resultCoupons))
	for _, cp := range resultCoupons {
		isAssigned := assignedMap[cp.ID]
		res = append(res, &dto.CouponResponseForCourse{
			CouponResponse: dto.NewCouponResponse(cp),
			IsAssigned:     isAssigned,
			IsDefault:      isAssigned && defaultMap[cp.ID],
		})
	}

	return res, total, nil
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

func buildCouponQuery(userID uuid.UUID, request dto.ListCouponRequest) (string, []any) {
	conditions := []string{"user_id = ?"}
	args := []any{userID}

	filters := map[string]*string{
		"code": request.Code,
	}

	for column, value := range filters {
		util.AddILIKECondition(&conditions, &args, column, value)
	}

	query := strings.Join(conditions, " AND ")

	return query, args
}

func buildAssignableCouponQuery(userID uuid.UUID, request dto.ListCouponRequest) (string, []any) {
	conditions := []string{
		"user_id = ?",
		"is_active = ?",
		"(expires_at IS NULL OR expires_at >= ?)",
		"(max_redemptions IS NULL OR current_redemptions < max_redemptions)",
	}
	args := []any{userID, true, time.Now().UTC()}

	util.AddILIKECondition(&conditions, &args, "code", request.Code)

	query := strings.Join(conditions, " AND ")

	return query, args
}

func validateCouponAvailability(coupon *model.Coupon) error {
	if coupon == nil {
		return apperror.NewBadRequestError("Coupon not found")
	}
	if !coupon.IsActive {
		return apperror.NewBadRequestError("Coupon is inactive")
	}

	now := time.Now().UTC()
	if coupon.ExpiresAt != nil && coupon.ExpiresAt.Before(now) {
		return apperror.NewBadRequestError("Coupon has expired")
	}
	if coupon.MaxRedemptions != nil && *coupon.MaxRedemptions > 0 && coupon.CurrentRedemptions >= *coupon.MaxRedemptions {
		return apperror.NewBadRequestError("Coupon has reached redemption limit")
	}

	return nil
}
