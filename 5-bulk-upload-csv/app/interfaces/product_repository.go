package interfaces

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/model"
	"context"
)

type ProductRepositoryInterface interface {
	GetList(ctx context.Context, params dto.GetListProductRequest) ([]model.Product, []error)
	GetDetail(ctx context.Context, params model.GetDetailProductParams) (*model.Product, []error)
	Create(ctx context.Context, data model.Product) (*model.Product, []error)
}
