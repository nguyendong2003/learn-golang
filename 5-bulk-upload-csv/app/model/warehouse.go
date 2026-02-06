package model

import "time"

type Warehouse struct {
	ModelUUID
	Code      string `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
