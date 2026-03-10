package model

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	BaseModel

	Name         string  `gorm:"type:varchar(100);not null;uniqueIndex"`
	Description  string  `gorm:"type:text"`
	MonthlyPrice float64 `gorm:"type:decimal(10,2);default:0"`
	YearlyPrice  float64 `gorm:"type:decimal(10,2);default:0"`
	IsDefault    bool    `gorm:"default:false"`
	IsActive     bool    `gorm:"default:true"`

	// Relations
	Subscriptions []*Subscription `gorm:"foreignKey:PlanID;references:ID"`
}

func (Plan) TableName() string {
	return "plans"
}

type Subscription struct {
	BaseModel

	UserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	PlanID       uuid.UUID `gorm:"type:uuid;not null;index"`
	BillingCycle string    `gorm:"type:varchar(50);not null"`         // monthly, yearly
	Status       string    `gorm:"type:varchar(50);default:'active'"` // active, canceled, expired, trailing
	StartedAt    time.Time `gorm:"not null"`
	EndedAt      *time.Time
	CanceledAt   *time.Time

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

	SubscriptionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Status         string    `gorm:"type:varchar(50);default:'pending'"` // pending, succeeded, failed, refunded
	PaidAt         *time.Time

	// Relations
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;references:ID"`
}

func (Payment) TableName() string {
	return "payments"
}
