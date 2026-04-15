package model

import "github.com/google/uuid"

type CourseCoupon struct {
	BaseModel

	CourseID  uuid.UUID `gorm:"type:uuid;not null;index"`
	CouponID  uuid.UUID `gorm:"type:uuid;not null;index"`
	IsDefault bool      `gorm:"not null;default:false"`

	Course *Course `gorm:"foreignKey:CourseID;references:ID"`
	Coupon *Coupon `gorm:"foreignKey:CouponID;references:ID"`
}

func (CourseCoupon) TableName() string {
	return "course_coupons"
}
