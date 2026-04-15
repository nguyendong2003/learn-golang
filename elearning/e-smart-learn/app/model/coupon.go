package model

import (
	"time"

	"github.com/google/uuid"
)

type Coupon struct {
	BaseModel

	UserID                *uuid.UUID `gorm:"type:uuid;index"`
	Code                  string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	StripeCouponID        string     `gorm:"type:varchar(255);not null;index"`
	StripePromotionCodeID string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	DiscountType          string     `gorm:"type:discount_type_enum;not null"`
	DiscountValue         int64      `gorm:"not null;default:0"`
	Currency              string     `gorm:"type:varchar(10);default:'usd'"`
	MaxRedemptions        *int64     `gorm:"default:null"`
	CurrentRedemptions    int64      `gorm:"not null;default:0"`
	IsActive              bool       `gorm:"not null;default:true"`
	ExpiresAt             *time.Time `gorm:"type:timestamptz"`

	User          *User           `gorm:"foreignKey:UserID;references:ID"`
	CourseCoupons []*CourseCoupon `gorm:"foreignKey:CouponID"`
}

func (Coupon) TableName() string {
	return "coupons"
}
