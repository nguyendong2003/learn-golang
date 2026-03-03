package repository

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) interfaces.ProductRepositoryInterface {
	return &ProductRepository{
		db: db,
	}
}

func (r ProductRepository) GetList(ctx context.Context, params dto.GetListProductRequest) ([]model.Product, []error) {
	limit := params.Limit
	offset := params.GetOffset()

	var result []model.Product
	var errs []error

	query := r.db.WithContext(ctx).
		Model(model.Product{}).
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

func (r ProductRepository) GetDetail(ctx context.Context, params model.GetDetailProductParams) (*model.Product, []error) {
	var result *model.Product
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

func (r ProductRepository) Create(ctx context.Context, data model.Product) (*model.Product, []error) {
	var errs []error

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return &data, nil
}
