package dto

import (
	"encoding/json"
	"time"

	"elearning-api/model"
)

type PlanResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	AccessFeatures []string `json:"access_features"`
	BillingCycle   string   `json:"billing_cycle"`
	Price          float64  `json:"price"`
	Currency       string   `json:"currency"`
	Tag            string   `json:"tag"`
	IsRecommend    bool     `json:"is_recommend"`
	IsActive       bool     `json:"is_active"`
}

type CreateSubscriptionCheckoutSessionRequest struct {
	PlanID string `json:"plan_id" binding:"required,uuid"`
}

type CheckoutSessionResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

type SubscriptionPlanInfo struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	BillingCycle string  `json:"billing_cycle"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
}

type MySubscriptionResponse struct {
	ID                   string                `json:"id"`
	StripeSubscriptionID string                `json:"stripe_subscription_id"`
	Status               string                `json:"status"`
	BillingCycle         string                `json:"billing_cycle"`
	CancelAtPeriodEnd    bool                  `json:"cancel_at_period_end"`
	CurrentPeriodStart   *time.Time            `json:"current_period_start"`
	CurrentPeriodEnd     *time.Time            `json:"current_period_end"`
	StartedAt            time.Time             `json:"started_at"`
	EndedAt              *time.Time            `json:"ended_at"`
	CanceledAt           *time.Time            `json:"canceled_at"`
	Plan                 *SubscriptionPlanInfo `json:"plan,omitempty"`
}

type BillingPortalResponse struct {
	URL string `json:"url"`
}

type MemberRetentionResponse struct {
	ActiveMemberships int64   `json:"active_memberships"`
	RetentionPct      float64 `json:"retention_pct"`
}

type PaymentResponse struct {
	ID                  string     `json:"id"`
	StripeInvoiceID     string     `json:"stripe_invoice_id"`
	StripePaymentIntent string     `json:"stripe_payment_intent"`
	Status              string     `json:"status"`
	Amount              float64    `json:"amount"`
	Currency            string     `json:"currency"`
	StripeFee           float64    `json:"stripe_fee"`
	FailureReason       string     `json:"failure_reason"`
	AttemptCount        int64      `json:"attempt_count"`
	PaidAt              *time.Time `json:"paid_at"`
	CreatedAt           time.Time  `json:"created_at"`
}

type SubscriptionPurchaseResponse struct {
	ID                      string                `json:"id"`
	StripeSubscriptionID    string                `json:"stripe_subscription_id"`
	StripeCheckoutSessionID string                `json:"stripe_checkout_session_id"`
	Status                  string                `json:"status"`
	BillingCycle            string                `json:"billing_cycle"`
	CancelAtPeriodEnd       bool                  `json:"cancel_at_period_end"`
	CurrentPeriodStart      *time.Time            `json:"current_period_start"`
	CurrentPeriodEnd        *time.Time            `json:"current_period_end"`
	StartedAt               time.Time             `json:"started_at"`
	EndedAt                 *time.Time            `json:"ended_at"`
	CanceledAt              *time.Time            `json:"canceled_at"`
	Plan                    *SubscriptionPlanInfo `json:"plan,omitempty"`
	User                    *UserResponse         `json:"user,omitempty"`
	Payments                []*PaymentResponse    `json:"payments,omitempty"`
}

func NewListPlanResponse(plans []*model.Plan) []*PlanResponse {
	res := make([]*PlanResponse, len(plans))
	for i, p := range plans {
		res[i] = NewPlanDetailResponse(p)
	}
	return res
}

func NewPlanDetailResponse(p *model.Plan) *PlanResponse {
	var features []string
	if p.AccessFeatures != "" {
		_ = json.Unmarshal([]byte(p.AccessFeatures), &features)
	}

	return &PlanResponse{
		ID:             p.ID.String(),
		Name:           p.Name,
		Description:    p.Description,
		AccessFeatures: features,
		BillingCycle:   p.BillingCycle,
		Price:          p.Price,
		Currency:       p.Currency,
		Tag:            p.Tag,
		IsRecommend:    p.IsRecommend,
		IsActive:       p.IsActive,
	}
}

func NewMySubscriptionResponse(s *model.Subscription) *MySubscriptionResponse {
	if s == nil {
		return nil
	}

	res := &MySubscriptionResponse{
		ID:                   s.ID.String(),
		StripeSubscriptionID: s.StripeSubscriptionID,
		Status:               s.Status,
		BillingCycle:         s.BillingCycle,
		CancelAtPeriodEnd:    s.CancelAtPeriodEnd,
		CurrentPeriodStart:   s.CurrentPeriodStart,
		CurrentPeriodEnd:     s.CurrentPeriodEnd,
		StartedAt:            s.StartedAt,
		EndedAt:              s.EndedAt,
		CanceledAt:           s.CanceledAt,
		Plan: &SubscriptionPlanInfo{
			ID:           s.PlanID.String(),
			Name:         s.PlanName,
			Description:  s.PlanDescription,
			BillingCycle: s.BillingCycle,
			Price:        s.PlanPrice,
			Currency:     s.PlanCurrency,
		},
	}

	return res
}

func NewSubscriptionPurchaseResponse(s *model.Subscription) *SubscriptionPurchaseResponse {
	if s == nil {
		return nil
	}

	res := &SubscriptionPurchaseResponse{
		ID:                      s.ID.String(),
		StripeSubscriptionID:    s.StripeSubscriptionID,
		StripeCheckoutSessionID: s.StripeCheckoutSessionID,
		Status:                  s.Status,
		BillingCycle:            s.BillingCycle,
		CancelAtPeriodEnd:       s.CancelAtPeriodEnd,
		CurrentPeriodStart:      s.CurrentPeriodStart,
		CurrentPeriodEnd:        s.CurrentPeriodEnd,
		StartedAt:               s.StartedAt,
		EndedAt:                 s.EndedAt,
		CanceledAt:              s.CanceledAt,
		Plan: &SubscriptionPlanInfo{
			ID:           s.PlanID.String(),
			Name:         s.PlanName,
			Description:  s.PlanDescription,
			BillingCycle: s.BillingCycle,
			Price:        s.PlanPrice,
			Currency:     s.PlanCurrency,
		},
		User: NewUserDetailResponse(s.User),
	}

	if len(s.Payments) > 0 {
		res.Payments = make([]*PaymentResponse, len(s.Payments))
		for i, p := range s.Payments {
			res.Payments[i] = &PaymentResponse{
				ID:                  p.ID.String(),
				StripeInvoiceID:     p.StripeInvoiceID,
				StripePaymentIntent: p.StripePaymentIntent,
				Status:              p.Status,
				Amount:              float64(p.Amount) / 100,
				Currency:            p.Currency,
				StripeFee:           float64(p.StripeFee) / 100,
				FailureReason:       p.FailureReason,
				AttemptCount:        p.AttemptCount,
				PaidAt:              p.PaidAt,
				CreatedAt:           p.CreatedAt,
			}
		}
	}

	return res
}

type SubscriberResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	Name          string     `json:"name"`
	AvatarURL     string     `json:"avatar_url"`
	Tier          string     `json:"tier"`
	BillingCycle  string     `json:"billing_cycle"`
	Status        string     `json:"status"`
	JoinedAt      time.Time  `json:"joined_at"`
	ExpiryAt      *time.Time `json:"expiry_at"`
	RenewalStatus string     `json:"renewal_status"`
}

type TransactionHistoryResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type ListTransactionsHistoryResponse struct {
	Transactions []*TransactionHistoryResponse `json:"transactions"`
	TotalAmount  float64                       `json:"total_amount"`
}

func NewSubscriberResponse(s *model.Subscription) *SubscriberResponse {
	if s == nil {
		return nil
	}
	renewalStatus := "none"
	if s.CancelAtPeriodEnd {
		renewalStatus = "manual"
	} else if s.Status == "active" || s.Status == "trialing" {
		renewalStatus = "auto-renew"
	}
	return &SubscriberResponse{
		ID:            s.User.ID.String(),
		Email:         s.User.Email,
		Name:          s.User.Name,
		AvatarURL:     s.User.Avatar,
		Tier:          s.PlanName,
		BillingCycle:  s.BillingCycle,
		Status:        s.Status,
		JoinedAt:      s.StartedAt,
		ExpiryAt:      s.CurrentPeriodEnd,
		RenewalStatus: renewalStatus,
	}
}

func NewListSubscribersResponse(subscribers []*model.Subscription) []*SubscriberResponse {
	res := make([]*SubscriberResponse, len(subscribers))
	for i, s := range subscribers {
		res[i] = NewSubscriberResponse(s)
	}
	return res
}
