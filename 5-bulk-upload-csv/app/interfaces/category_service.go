package interfaces

import (
	"bulk-upload-csv/dto"
	"context"
)

type CategoryServiceInterface interface {
	GetList(ctx context.Context) ([]dto.CategoryResponse, []error)
	GetDetail(ctx context.Context, params dto.GetCategoryDetailRequest) (*dto.CategoryResponse, []error)
	Create(ctx context.Context, data dto.CreateCategoryRequest) (*dto.CategoryResponse, []error)
}
