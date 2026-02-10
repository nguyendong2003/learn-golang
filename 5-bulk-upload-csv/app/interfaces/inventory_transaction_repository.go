package interfaces

import (
	"bulk-upload-csv/model"
	"context"

	"github.com/google/uuid"
)

type InventoryTransactionRepositoryInterface interface {
	BatchInsert(ctx context.Context, data []model.InventoryTransaction) error
	InventoryWorker(jobs <-chan []model.InventoryTransaction, errCh chan<- error)
	LoadCategoryCache(ctx context.Context) (map[string]uuid.UUID, error)
	LoadWarehouseCache(ctx context.Context) (map[string]uuid.UUID, error)
	LoadProductCache(ctx context.Context) (map[string]uuid.UUID, error)
}
