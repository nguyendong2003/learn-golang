package interfaces

import (
	"bulk-upload-csv/model"
	"context"

	"gorm.io/gorm"
)

type InventoryTransactionRepositoryInterface interface {
	BatchInsertTransactions(ctx context.Context, data []model.InventoryTransaction) error
	BatchInsertTransactionsWithTx(ctx context.Context, tx *gorm.DB, data []model.InventoryTransaction) error
	GetAllProducts(ctx context.Context) ([]model.ProductIDAndSku, error)
	GetAllCategories(ctx context.Context) ([]model.CategoryIDAndCode, error)
	GetAllWarehouses(ctx context.Context) ([]model.WarehouseIDAndCode, error)
	BeginTx(ctx context.Context) *gorm.DB
}
