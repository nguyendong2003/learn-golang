package model

import (
	"time"

	"github.com/google/uuid"
)

type CoursePurchase struct {
	BaseModel

	UserID                  uuid.UUID  `gorm:"type:uuid;not null;index"`
	CouponID                *uuid.UUID `gorm:"type:uuid;index"`
	StripeCheckoutSessionID string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	StripePaymentIntentID   string     `gorm:"type:varchar(255);index"`
	Amount                  int64      `gorm:"not null;default:0"`
	Currency                string     `gorm:"type:varchar(10);not null;default:'usd'"`
	StripeFee               int64      `gorm:"not null;default:0"`
	Status                  string     `gorm:"type:course_purchase_status_enum;not null;default:'pending'"`
	PurchasedAt             *time.Time

	User    *User                   `gorm:"foreignKey:UserID;references:ID"`
	Coupon  *Coupon                 `gorm:"foreignKey:CouponID;references:ID"`
	Details []*CoursePurchaseDetail `gorm:"foreignKey:CoursePurchaseID;references:ID"`
}

func (CoursePurchase) TableName() string {
	return "course_purchases"
}

type CoursePurchaseDetail struct {
	BaseModel

	CoursePurchaseID uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID         uuid.UUID `gorm:"type:uuid;not null;index"`
	Price            int64     `gorm:"not null;default:0"`
	Currency         string    `gorm:"type:varchar(10);not null;default:'usd'"`

	CoursePurchase *CoursePurchase `gorm:"foreignKey:CoursePurchaseID;references:ID"`
	Course         *Course         `gorm:"foreignKey:CourseID;references:ID"`
}

func (CoursePurchaseDetail) TableName() string {
	return "course_purchase_details"
}
