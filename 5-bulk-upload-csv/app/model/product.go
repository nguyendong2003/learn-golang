package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ModelUUID
	Sku        string    `gorm:"size:50;uniqueIndex;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
