package dto

import (
	"elearning-api/consts"
	"elearning-api/model"
	"strings"
	"time"
)

type CouponResponse struct {
	ID                    string     `json:"id"`
	Code                  string     `json:"code"`
	StripeCouponID        string     `json:"stripe_coupon_id"`
	StripePromotionCodeID string     `json:"stripe_promotion_code_id"`
	DiscountType          string     `json:"discount_type"`
	DiscountValue         float64    `json:"discount_value"`
	Currency              string     `json:"currency"`
	MaxRedemptions        *int64     `json:"max_redemptions,omitempty"`
	CurrentRedemptions    int64      `json:"current_redemptions"`
	IsActive              bool       `json:"is_active"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
}

type CouponResponseForCourse struct {
	*CouponResponse
	IsAssigned bool `json:"is_assigned"`
	IsDefault  bool `json:"is_default"`
}

func NewCouponResponse(coupon *model.Coupon) *CouponResponse {
	if coupon == nil {
		return nil
	}

	discountType := strings.TrimSpace(coupon.DiscountType)
	discountValue := float64(coupon.DiscountValue)
	if strings.EqualFold(discountType, consts.DiscountTypeAmount) {
		discountValue = discountValue / 100
	}

	return &CouponResponse{
		ID:                    coupon.ID.String(),
		Code:                  coupon.Code,
		StripeCouponID:        coupon.StripeCouponID,
		StripePromotionCodeID: coupon.StripePromotionCodeID,
		DiscountType:          discountType,
		DiscountValue:         discountValue,
		Currency:              coupon.Currency,
		MaxRedemptions:        coupon.MaxRedemptions,
		CurrentRedemptions:    coupon.CurrentRedemptions,
		IsActive:              coupon.IsActive,
		ExpiresAt:             coupon.ExpiresAt,
	}
}

func NewListCouponResponse(coupons []*model.Coupon) []*CouponResponse {
	res := make([]*CouponResponse, len(coupons))
	for i, coupon := range coupons {
		res[i] = NewCouponResponse(coupon)
	}
	return res
}

type CreateCoursePurchaseCheckoutSessionRequest struct {
	CouponCode string `json:"coupon_code" binding:"omitempty,max=100"`
}

type CreateCouponRequest struct {
	Code           string     `json:"code" binding:"required,min=3,max=100"`
	DiscountType   string     `json:"discount_type" binding:"required,oneof=percent amount"`
	DiscountValue  int64      `json:"discount_value" binding:"required,gt=0"`
	Currency       *string    `json:"currency" binding:"omitempty,len=3"`
	MaxRedemptions *int64     `json:"max_redemptions" binding:"omitempty,gt=0"`
	ExpiresAt      *time.Time `json:"expires_at"`
}

type ListCouponRequest struct {
	PagingRequest

	Code     *string `form:"code,omitempty"`
	CourseID *string `form:"course_id,omitempty"`
}
