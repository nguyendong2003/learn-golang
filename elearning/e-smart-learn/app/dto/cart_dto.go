package dto

import (
	"elearning-api/consts"
	"elearning-api/model"
	"strings"
)

type CartCheckoutRequest struct {
	CouponCode string `json:"coupon_code" binding:"omitempty,max=100"`
}

type CartCheckoutCouponResponse struct {
	ID            string  `json:"id"`
	Code          string  `json:"code"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
}

func NewCartCheckoutCouponResponse(coupon *model.Coupon) *CartCheckoutCouponResponse {
	if coupon == nil {
		return nil
	}

	discountType := strings.TrimSpace(coupon.DiscountType)
	discountValue := float64(coupon.DiscountValue)
	if strings.EqualFold(discountType, consts.DiscountTypeAmount) {
		discountValue = discountValue / 100
	}

	return &CartCheckoutCouponResponse{
		ID:            coupon.ID.String(),
		Code:          strings.TrimSpace(coupon.Code),
		DiscountType:  discountType,
		DiscountValue: discountValue,
	}
}

type CartCheckoutPreviewItemResponse struct {
	Course *CourseResponse             `json:"course"`
	Coupon *CartCheckoutCouponResponse `json:"coupon,omitempty"`
}

type CartCheckoutPreviewResponse struct {
	Items       []*CartCheckoutPreviewItemResponse `json:"items"`
	TotalAmount float64                            `json:"total_amount"`
	Currency    string                             `json:"currency"`
}

type CartResponse struct {
	Items       []*CourseResponse `json:"items"`
	TotalAmount float64           `json:"total_amount"`
	Currency    string            `json:"currency"`
}

func NewCartResponse(items []*model.CartItem) *CartResponse {
	res := &CartResponse{Items: make([]*CourseResponse, 0, len(items))}

	for _, item := range items {
		if item == nil || item.Course == nil {
			continue
		}
		courseResp := NewCourseDetailResponse(item.Course)
		res.Items = append(res.Items, courseResp)
		res.TotalAmount += item.Course.Price
		if item.Course.StripeCurrency != "" {
			res.Currency = item.Course.StripeCurrency
		}
	}

	if res.Currency == "" {
		res.Currency = "usd"
	}

	return res
}
