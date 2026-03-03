package interfaces

import (
	"bulk-upload-csv/dto"
	"context"
)

type WarehouseServiceInterface interface {
	GetList(ctx context.Context, params dto.GetListWarehouseRequest) ([]dto.WarehouseResponse, []error)
	GetDetail(ctx context.Context, params dto.GetWarehouseDetailRequest) (*dto.WarehouseResponse, []error)
	Create(ctx context.Context, data dto.CreateWarehouseRequest) (*dto.WarehouseResponse, []error)
}
