package model

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	BaseModel

	Name            string  `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description     string  `gorm:"type:text"`
	BillingCycle    string  `gorm:"type:billing_cycle_enum;not null;index"` // monthly, yearly
	Price           float64 `gorm:"type:decimal(10,2);not null;default:0"`
	StripePriceID   string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	StripeProductID string  `gorm:"type:varchar(255);not null;uniqueIndex"`
	Currency        string  `gorm:"type:varchar(10);default:'usd'"`
	Tag             string  `gorm:"type:text"`
	IsRecommend     bool    `gorm:"default:false"`
	AccessFeatures  string  `gorm:"type:jsonb;default:'[]'"` // JSON array of feature strings
	IsActive        bool    `gorm:"default:true"`

	// Relations
	Subscriptions []*Subscription `gorm:"foreignKey:PlanID;references:ID"`
}

func (Plan) TableName() string {
	return "plans"
}

type Subscription struct {
	BaseModel

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	PlanID            uuid.UUID `gorm:"type:uuid;not null;index"`
	PlanName          string    `gorm:"type:varchar(100);not null;default:''"`
	PlanDescription   string    `gorm:"type:text"`
	PlanPrice         float64   `gorm:"type:decimal(10,2);not null;default:0"`
	PlanCurrency      string    `gorm:"type:varchar(10);not null;default:'usd'"`
	PlanStripePriceID string    `gorm:"type:varchar(255);index"`

	StripeCheckoutSessionID string `gorm:"type:varchar(255);index"`
	StripeSubscriptionID    string `gorm:"type:varchar(255);uniqueIndex"`
	StripeCustomerID        string `gorm:"type:varchar(255);index"`
	BillingCycle            string `gorm:"type:billing_cycle_enum;not null"`
	Status                  string `gorm:"type:subscription_status_enum;default:'incomplete'"`
	CurrentPeriodStart      *time.Time
	CurrentPeriodEnd        *time.Time
	CancelAtPeriodEnd       bool      `gorm:"default:false"`
	StartedAt               time.Time `gorm:"not null"`
	EndedAt                 *time.Time
	CanceledAt              *time.Time

	// Relations
	User     *User      `gorm:"foreignKey:UserID;references:ID"`
	Plan     *Plan      `gorm:"foreignKey:PlanID;references:ID"`
	Payments []*Payment `gorm:"foreignKey:SubscriptionID;references:ID"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}

type Payment struct {
	BaseModel

	SubscriptionID      uuid.UUID `gorm:"type:uuid;not null;index"`
	StripeInvoiceID     string    `gorm:"type:varchar(255);index"`
	StripePaymentIntent string    `gorm:"type:varchar(255);index"`
	Status              string    `gorm:"type:payment_status_enum;default:'pending'"` // pending, succeeded, failed, refunded
	Amount              int64     `gorm:"default:0"`
	Currency            string    `gorm:"type:varchar(10);default:'usd'"`
	StripeFee           int64     `gorm:"default:0"`
	FailureReason       string    `gorm:"type:text"`
	AttemptCount        int64     `gorm:"default:0"`
	PaidAt              *time.Time

	// Relations
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;references:ID"`
}

func (Payment) TableName() string {
	return "payments"
}
