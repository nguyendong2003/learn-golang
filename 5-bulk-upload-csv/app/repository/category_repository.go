package repository

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) interfaces.CategoryRepositoryInterface {
	return &CategoryRepository{
		db: db,
	}
}

func (r CategoryRepository) GetList(ctx context.Context, params dto.GetListCategoryRequest) ([]model.Category, []error) {
	limit := params.Limit
	offset := params.GetOffset()

	var result []model.Category
	var errs []error

	query := r.db.WithContext(ctx).
		Model(model.Category{}).
		Order("updated_at desc").
		Limit(limit).
		Offset(offset).
		Find(&result)
	if err := query.Error; err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return result, nil
}

func (r CategoryRepository) GetDetail(ctx context.Context, params model.GetDetailCategoryParams) (*model.Category, []error) {
	var result *model.Category
	var errs []error

	paramsMap, err := params.Map()
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	if err := r.db.WithContext(ctx).Where(paramsMap).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errs = append(errs, err)
			return nil, errs
		}

		errs = append(errs, err)
		return nil, errs
	}

	return result, nil

}

func (r CategoryRepository) Create(ctx context.Context, data model.Category) (*model.Category, []error) {
	var errs []error

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return &data, nil
}
