package repository

import (
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"

	"github.com/google/uuid"
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

func (r InventoryTransactionRepository) BatchInsert(ctx context.Context, data []model.InventoryTransaction) error {
	const batchSize = 2000
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(data); i += batchSize {
			end := i + batchSize
			if end > len(data) {
				end = len(data)
			}
			if err := tx.CreateInBatches(data[i:end], batchSize).Error; err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (r InventoryTransactionRepository) InventoryWorker(jobs <-chan []model.InventoryTransaction, errCh chan<- error) {
	for batch := range jobs {
		err := r.db.Transaction(func(tx *gorm.DB) error {
			return tx.Create(&batch).Error
		})
		if err != nil {
			errCh <- err
			return
		}
	}
	errCh <- nil
}

func (r InventoryTransactionRepository) LoadCategoryCache(ctx context.Context) (map[string]uuid.UUID, error) {
	var categories []model.Category
	if err := r.db.WithContext(ctx).Model(&model.Category{}).Find(&categories).Error; err != nil {
		return nil, err
	}

	categoryCache := make(map[string]uuid.UUID)
	for _, category := range categories {
		categoryCache[category.Code] = category.ID
	}

	return categoryCache, nil
}

func (r InventoryTransactionRepository) LoadWarehouseCache(ctx context.Context) (map[string]uuid.UUID, error) {
	var warehouses []model.Warehouse
	if err := r.db.WithContext(ctx).Model(&model.Warehouse{}).Find(&warehouses).Error; err != nil {
		return nil, err
	}

	warehouseCache := make(map[string]uuid.UUID)
	for _, warehouse := range warehouses {
		warehouseCache[warehouse.Code] = warehouse.ID
	}

	return warehouseCache, nil
}

func (r InventoryTransactionRepository) LoadProductCache(ctx context.Context) (map[string]uuid.UUID, error) {
	var products []model.Product
	if err := r.db.WithContext(ctx).Model(&model.Product{}).Find(&products).Error; err != nil {
		return nil, err
	}

	productCache := make(map[string]uuid.UUID)
	for _, product := range products {
		productCache[product.Sku] = product.ID
	}

	return productCache, nil
}
