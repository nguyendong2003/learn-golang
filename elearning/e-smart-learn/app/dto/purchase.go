package dto

import (
	"elearning-api/model"
	"time"
)

type CoursePurchaseDetailResponse struct {
	ID        string    `json:"id"`
	CourseID  string    `json:"course_id"`
	Price     int64     `json:"price"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type CoursePurchaseResponse struct {
	ID                string                          `json:"id"`
	UserID            string                          `json:"user_id"`
	Amount            int64                           `json:"amount"`
	Currency          string                          `json:"currency"`
	StripeFee         int64                           `json:"stripe_fee"`
	CheckoutSessionID string                          `json:"checkout_session_id"`
	Status            string                          `json:"status"`
	CreatedAt         time.Time                       `json:"created_at"`
	PurchasedAt       *time.Time                      `json:"purchased_at"`
	Details           []*CoursePurchaseDetailResponse `json:"details,omitempty"`
}

func NewCoursePurchaseResponse(purchase *model.CoursePurchase) *CoursePurchaseResponse {
	resp := &CoursePurchaseResponse{
		ID:                purchase.ID.String(),
		UserID:            purchase.UserID.String(),
		Amount:            purchase.Amount,
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
				ID:        d.ID.String(),
				CourseID:  d.CourseID.String(),
				Price:     d.Price,
				Currency:  d.Currency,
				CreatedAt: d.CreatedAt,
			})
		}
		resp.Details = details
	}

	return resp
}
