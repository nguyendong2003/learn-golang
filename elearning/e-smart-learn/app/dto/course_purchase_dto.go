package dto

import (
	"elearning-api/model"
	"strings"
	"time"
)

type CoursePurchaseResponse struct {
	ID                string                          `json:"id"`
	UserID            string                          `json:"user_id"`
	TotalAmount       int64                           `json:"total_amount"`
	Currency          string                          `json:"currency"`
	StripeFee         int64                           `json:"stripe_fee"`
	CheckoutSessionID string                          `json:"checkout_session_id"`
	Status            string                          `json:"status"`
	CreatedAt         time.Time                       `json:"created_at"`
	PurchasedAt       *time.Time                      `json:"purchased_at"`
	Details           []*CoursePurchaseDetailResponse `json:"details,omitempty"`
}

type CoursePurchaseDetailResponse struct {
	ID            string    `json:"id"`
	CourseID      string    `json:"course_id"`
	PriceOriginal int64     `json:"price_original"`
	PriceFinal    int64     `json:"price_final"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

type CoursePurchaseCheckoutPreviewResponse struct {
	Course        *CourseResponse             `json:"course"`
	Coupon        *CartCheckoutCouponResponse `json:"coupon,omitempty"`
	PriceOriginal float64                     `json:"price_original"`
	PriceFinal    float64                     `json:"price_final"`
	Currency      string                      `json:"currency"`
}

func NewCoursePurchaseResponse(purchase *model.CoursePurchase) *CoursePurchaseResponse {
	resp := &CoursePurchaseResponse{
		ID:                purchase.ID.String(),
		UserID:            purchase.UserID.String(),
		TotalAmount:       purchase.TotalAmount,
		Currency:          purchase.Currency,
		StripeFee:         purchase.StripeFee,
		CheckoutSessionID: purchase.StripeCheckoutSessionID,
		Status:            purchase.Status,
		CreatedAt:         purchase.CreatedAt,
		PurchasedAt:       purchase.PurchasedAt,
	}

	if len(purchase.Details) > 0 {
		details := make([]*CoursePurchaseDetailResponse, 0, len(purchase.Details))
		for _, d := range purchase.Details {
			if d == nil {
				continue
			}
			details = append(details, &CoursePurchaseDetailResponse{
				ID:            d.ID.String(),
				CourseID:      d.CourseID.String(),
				PriceOriginal: d.PriceOriginal,
				PriceFinal:    d.PriceFinal,
				Currency:      d.Currency,
				CreatedAt:     d.CreatedAt,
			})
		}
		resp.Details = details
	}

	return resp
}

func NewCoursePurchaseCheckoutPreviewResponse(course *model.Course, coupon *model.Coupon, originalAmount, totalAmount float64, currency string) *CoursePurchaseCheckoutPreviewResponse {
	if course == nil {
		return nil
	}

	res := &CoursePurchaseCheckoutPreviewResponse{
		Course:        NewCourseDetailResponse(course),
		PriceOriginal: originalAmount,
		PriceFinal:    totalAmount,
		Currency:      strings.ToLower(strings.TrimSpace(currency)),
	}

	if res.Currency == "" {
		res.Currency = "usd"
	}

	if coupon != nil {
		res.Coupon = NewCartCheckoutCouponResponse(coupon)
	}

	return res
}
