package model

import (
	"time"

	"github.com/google/uuid"
)

type InventoryTransaction struct {
	ModelUUID
	ProductID       uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID      uuid.UUID `gorm:"type:uuid;not null;index"`
	WarehouseID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Quantity        int       `gorm:"not null"`
	TransactionType string    `gorm:"size:10;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
