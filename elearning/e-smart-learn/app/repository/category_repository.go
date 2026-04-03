package repository

import (
	"context"
	"elearning-api/model"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	Repository[model.Category]

	GetByName(ctx context.Context, name string) (*model.Category, error)
}

type categoryRepository struct {
	*repository[model.Category]
}

func NewCategoryRepository(db DbRepository) CategoryRepository {
	return &categoryRepository{
		repository: NewBaseRepository[model.Category](db),
	}
}

func (r *categoryRepository) GetByName(
	ctx context.Context,
	name string,
) (*model.Category, error) {
	var category model.Category
	err := r.baseQuery(ctx).
		Where("name = ?", name).
		First(&category).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &category, nil
}
