package repository

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
	"errors"

	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) interfaces.WarehouseRepositoryInterface {
	return &WarehouseRepository{
		db: db,
	}
}

func (r WarehouseRepository) GetList(ctx context.Context, params dto.GetListWarehouseRequest) ([]model.Warehouse, []error) {
	limit := params.Limit
	offset := params.GetOffset()

	var result []model.Warehouse
	var errs []error

	query := r.db.WithContext(ctx).
		Model(model.Warehouse{}).
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

func (r WarehouseRepository) GetDetail(ctx context.Context, params model.GetDetailWarehouseParams) (*model.Warehouse, []error) {
	var result *model.Warehouse
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

func (r WarehouseRepository) Create(ctx context.Context, data model.Warehouse) (*model.Warehouse, []error) {
	var errs []error

	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return &data, nil
}
