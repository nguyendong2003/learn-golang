package model

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ModelUUID
	Code      string `gorm:"size:50;uniqueIndex;not null"`
	Name      string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetDetailCategoryParams struct {
	Id *string
}

func (p GetDetailCategoryParams) Map() (Category, error) {
	var category Category
	if p.Id != nil {
		id, err := uuid.Parse(*p.Id)
		if err != nil {
			return category, err
		}

		category.ID = id
	}

	return category, nil
}
