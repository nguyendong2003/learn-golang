package interfaces

import (
	"bulk-upload-csv/dto"
	"context"
)

type ProductServiceInterface interface {
	GetList(ctx context.Context, params dto.GetListProductRequest) ([]dto.ProductResponse, []error)
	GetDetail(ctx context.Context, params dto.GetProductDetailRequest) (*dto.ProductResponse, []error)
	Create(ctx context.Context, data dto.CreateProductRequest) (*dto.ProductResponse, []error)
}
