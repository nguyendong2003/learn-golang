package interfaces

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/model"
	"context"
)

type ProductRepositoryInterface interface {
	GetList(ctx context.Context, params dto.GetListProductRequest) ([]model.Product, []error)
	GetProductsSkus(ctx context.Context, skus []string) ([]model.ProductIDAndSku, []error)
	GetAllProducts(ctx context.Context) ([]model.ProductIDAndSku, error)
	GetDetail(ctx context.Context, params model.GetDetailProductParams) (*model.Product, []error)
	Create(ctx context.Context, data model.Product) (*model.Product, []error)
}
