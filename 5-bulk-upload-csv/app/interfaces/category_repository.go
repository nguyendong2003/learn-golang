package interfaces

import (
	"bulk-upload-csv/model"
	"context"
)

type CategoryRepositoryInterface interface {
	GetList(ctx context.Context) ([]model.Category, []error)
	GetDetail(ctx context.Context, params model.GetDetailCategoryParams) (*model.Category, []error)
	Create(ctx context.Context, data model.Category) (*model.Category, []error)
}
