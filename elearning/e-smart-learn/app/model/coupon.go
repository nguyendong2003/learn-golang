package model

import (
	"time"
)

type Coupon struct {
	BaseModel

	Code                  string `gorm:"type:varchar(100);not null;uniqueIndex"`
	StripeCouponID        string `gorm:"type:varchar(255);not null;index"`
	StripePromotionCodeID string `gorm:"type:varchar(255);not null;uniqueIndex"`
	DiscountType          string `gorm:"type:discount_type_enum;not null"`
	DiscountValue         int64  `gorm:"default:0"`
	Currency              string `gorm:"type:varchar(10);default:'usd'"`
	MaxRedemptions        int64  `gorm:"default:0"`
	CurrentRedemptions    int64  `gorm:"default:0"`
	IsActive              bool   `gorm:"default:true"`
	ExpiresAt             *time.Time
}

func (Coupon) TableName() string {
	return "coupons"
}
