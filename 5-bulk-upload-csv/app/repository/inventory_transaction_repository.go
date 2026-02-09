package repository

import (
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"

	"gorm.io/gorm"
)

type InventoryTransactionRepository struct {
	db *gorm.DB
}

func NewInventoryTransactionRepository(db *gorm.DB) interfaces.InventoryTransactionRepositoryInterface {
	return &InventoryTransactionRepository{
		db: db,
	}
}

func (r *InventoryTransactionRepository) BatchInsertTransactions(ctx context.Context, data []model.InventoryTransaction) error {
	// Sử dụng CreateInBatches để GORM tự động chia nhỏ SQL nếu mảng quá lớn
	return r.db.WithContext(ctx).CreateInBatches(data, 1000).Error
}

func (r *InventoryTransactionRepository) BatchInsertTransactionsWithTx(ctx context.Context, tx *gorm.DB, data []model.InventoryTransaction) error {
	// Insert with transaction context
	return tx.WithContext(ctx).CreateInBatches(data, 1000).Error
}

func (r *InventoryTransactionRepository) BeginTx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}

func (r *InventoryTransactionRepository) GetAllProducts(ctx context.Context) ([]model.ProductIDAndSku, error) {
	var result []model.ProductIDAndSku
	if err := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Select("id, sku").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *InventoryTransactionRepository) GetAllCategories(ctx context.Context) ([]model.CategoryIDAndCode, error) {
	var result []model.CategoryIDAndCode
	if err := r.db.WithContext(ctx).
		Model(&model.Category{}).
		Select("id, code").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *InventoryTransactionRepository) GetAllWarehouses(ctx context.Context) ([]model.WarehouseIDAndCode, error) {
	var result []model.WarehouseIDAndCode
	if err := r.db.WithContext(ctx).
		Model(&model.Warehouse{}).
		Select("id, code").
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}
