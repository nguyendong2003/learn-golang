package interfaces

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/model"
	"context"
)

type WarehouseRepositoryInterface interface {
	GetList(ctx context.Context, params dto.GetListWarehouseRequest) ([]model.Warehouse, []error)
	GetDetail(ctx context.Context, params model.GetDetailWarehouseParams) (*model.Warehouse, []error)
	Create(ctx context.Context, data model.Warehouse) (*model.Warehouse, []error)
}
