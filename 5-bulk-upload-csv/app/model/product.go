package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ModelUUID
	Sku        string    `gorm:"size:50;uniqueIndex;not null"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type ProductIDAndSku struct {
	ModelUUID
	Sku string `gorm:"size:50;uniqueIndex;not null"`
}

type GetDetailProductParams struct {
	Id *string
}

func (p GetDetailProductParams) Map() (Product, error) {
	var product Product
	if p.Id != nil {
		id, err := uuid.Parse(*p.Id)
		if err != nil {
			return product, err
		}

		product.ID = id
	}

	return product, nil
}
