package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Warehouse struct {
	ModelUUID
	Code      string `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type WarehouseIDAndCode struct {
	ModelUUID
	Code string `gorm:"size:50;uniqueIndex;not null"`
}

type GetDetailWarehouseParams struct {
	Id *string
}

func (p GetDetailWarehouseParams) Map() (Warehouse, error) {
	var warehouse Warehouse
	if p.Id != nil {
		id, err := uuid.Parse(*p.Id)
		if err != nil {
			return warehouse, err
		}

		warehouse.ID = id
	}

	return warehouse, nil
}
