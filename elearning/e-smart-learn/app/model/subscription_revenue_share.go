package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionRevenueShare is a persistent ledger row for one instructor's
// revenue split from one successful subscription payment cycle.
type SubscriptionRevenueShare struct {
	BaseModel

	PaymentID             uuid.UUID `gorm:"type:uuid;not null;index"`
	SubscriptionID        uuid.UUID `gorm:"type:uuid;not null;index"`
	SubscriberUserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	InstructorUserID      uuid.UUID `gorm:"type:uuid;not null;index"`
	TotalActiveCourses    int       `gorm:"not null;default:0"`
	InstructorActiveCount int       `gorm:"not null;default:0"`
	PaymentAmount         int64     `gorm:"not null;default:0"`
	PaymentStripeFee      int64     `gorm:"not null;default:0"`
	AllocatedAmount       int64     `gorm:"not null;default:0"`
	InstructorGross       int64     `gorm:"not null;default:0"`
	PlatformGross         int64     `gorm:"not null;default:0"`
	AllocatedStripeFee    int64     `gorm:"not null;default:0"`
	InstructorNet         int64     `gorm:"not null;default:0"`
	PlatformNet           int64     `gorm:"not null;default:0"`
	BillingPeriodStart    *time.Time
	BillingPeriodEnd      *time.Time
	PaidAt                *time.Time

	Payment      *Payment      `gorm:"foreignKey:PaymentID;references:ID"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;references:ID"`
	Subscriber   *User         `gorm:"foreignKey:SubscriberUserID;references:ID"`
	Instructor   *User         `gorm:"foreignKey:InstructorUserID;references:ID"`
}

func (SubscriptionRevenueShare) TableName() string {
	return "subscription_revenue_shares"
}
