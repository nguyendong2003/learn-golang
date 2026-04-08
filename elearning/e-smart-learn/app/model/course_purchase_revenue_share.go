package model

import (
	"time"

	"github.com/google/uuid"
)

type CoursePurchaseRevenueShare struct {
	BaseModel

	CoursePurchaseID       uuid.UUID `gorm:"type:uuid;not null;index"`
	CoursePurchaseDetailID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	PurchaserUserID        uuid.UUID `gorm:"type:uuid;not null;index"`
	InstructorUserID       uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID               uuid.UUID `gorm:"type:uuid;not null;index"`
	PurchaseAmount         int64     `gorm:"not null;default:0"`
	PurchaseStripeFee      int64     `gorm:"not null;default:0"`
	DetailAmount           int64     `gorm:"not null;default:0"`
	DetailCouponDiscount   int64     `gorm:"not null;default:0"`
	DetailNetAmount        int64     `gorm:"not null;default:0"`
	AllocatedStripeFee     int64     `gorm:"not null;default:0"`
	InstructorGross        int64     `gorm:"not null;default:0"`
	PlatformGross          int64     `gorm:"not null;default:0"`
	InstructorNet          int64     `gorm:"not null;default:0"`
	PlatformNet            int64     `gorm:"not null;default:0"`
	PurchasedAt            *time.Time
}

func (CoursePurchaseRevenueShare) TableName() string {
	return "course_purchase_revenue_shares"
}
