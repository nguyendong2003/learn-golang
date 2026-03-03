package service

import (
	"bufio"
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"strconv"
)

type InventoryTransactionService struct {
	inventoryTransactionRepository interfaces.InventoryTransactionRepositoryInterface
}

func NewInventoryTransactionService(
	inventoryTransactionRepository interfaces.InventoryTransactionRepositoryInterface,
) interfaces.InventoryTransactionServiceInterface {
	return &InventoryTransactionService{
		inventoryTransactionRepository: inventoryTransactionRepository,
	}
}

func (s *InventoryTransactionService) ProcessBulkUpload(ctx context.Context, file io.Reader) (*dto.BulkUploadResponse, error) {
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	// read header
	header, err := reader.Read()
	if err != nil {
		return nil, errors.New("invalid csv header")
	}

	expectedHeader := []string{"product_sku", "category_code", "warehouse_code", "quantity", "transaction_type"}
	if !reflect.DeepEqual(header, expectedHeader) {
		return nil, errors.New("csv header not match")
	}

	categoryCache, err := s.inventoryTransactionRepository.LoadCategoryCache(ctx)
	if err != nil {
		return nil, err
	}

	warehouseCache, err := s.inventoryTransactionRepository.LoadWarehouseCache(ctx)
	if err != nil {
		return nil, err
	}

	productCache, err := s.inventoryTransactionRepository.LoadProductCache(ctx)
	if err != nil {
		return nil, err
	}

	var errors []dto.RowError
	var inventories []model.InventoryTransaction
	rowIndex := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowIndex++

		productSku, ok := productCache[record[0]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "product_sku",
				Value:   record[0],
				Message: "product_sku not found",
			})
			continue
		}

		categoryCode, ok := categoryCache[record[1]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "category_code",
				Value:   record[1],
				Message: "category_code not found",
			})
			continue
		}

		warehouseCode, ok := warehouseCache[record[2]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "warehouse_code",
				Value:   record[2],
				Message: "warehouse_code not found",
			})
			continue
		}

		transactionType := record[4]
		if transactionType != "IN" && transactionType != "OUT" {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "transaction_type",
				Value:   record[4],
				Message: "transaction_type must be 'IN' or 'OUT'",
			})
			continue
		}

		quantity, err := strconv.Atoi(record[3])
		if err != nil {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be a number",
			})
			continue
		}

		if transactionType == "IN" && quantity <= 0 {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be > 0 for IN transaction",
			})
			continue
		}

		if transactionType == "OUT" && quantity >= 0 {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be < 0 for OUT transaction",
			})
			continue
		}

		inventories = append(inventories, model.InventoryTransaction{
			ProductID:       productSku,
			CategoryID:      categoryCode,
			WarehouseID:     warehouseCode,
			Quantity:        quantity,
			TransactionType: transactionType,
		})
	}

	if len(errors) > 0 {
		return &dto.BulkUploadResponse{
			TotalProcessed: rowIndex - 1,
			TotalSuccess:   0,
			Errors:         errors,
		}, nil
	}

	// TRANSACTION + BATCH INSERT
	err = s.inventoryTransactionRepository.BatchInsert(ctx, inventories)
	if err != nil {
		return nil, err
	}

	return &dto.BulkUploadResponse{
		TotalProcessed: rowIndex - 1,
		TotalSuccess:   len(inventories),
		Errors:         nil,
	}, nil

}

// Goroutine + Channel + Batch Insert
func (s *InventoryTransactionService) ProcessBulkUploadGoroutine(ctx context.Context, file io.Reader) (*dto.BulkUploadResponse, error) {
	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	// read header
	header, err := reader.Read()
	if err != nil {
		return nil, errors.New("invalid csv header")
	}

	expectedHeader := []string{"product_sku", "category_code", "warehouse_code", "quantity", "transaction_type"}
	if !reflect.DeepEqual(header, expectedHeader) {
		return nil, errors.New("csv header not match")
	}

	categoryCache, err := s.inventoryTransactionRepository.LoadCategoryCache(ctx)
	if err != nil {
		return nil, err
	}

	warehouseCache, err := s.inventoryTransactionRepository.LoadWarehouseCache(ctx)
	if err != nil {
		return nil, err
	}

	productCache, err := s.inventoryTransactionRepository.LoadProductCache(ctx)
	if err != nil {
		return nil, err
	}

	const (
		BatchSize   = 2000
		WorkerCount = 4
	)

	jobs := make(chan []model.InventoryTransaction, WorkerCount)
	errCh := make(chan error, WorkerCount)

	// start workers
	for i := 0; i < WorkerCount; i++ {
		go s.inventoryTransactionRepository.InventoryWorker(jobs, errCh)
	}

	var errors []dto.RowError
	batch := make([]model.InventoryTransaction, 0, BatchSize)
	rowIndex := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowIndex++

		productSku, ok := productCache[record[0]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "product_sku",
				Value:   record[0],
				Message: "product_sku not found",
			})
			continue
		}

		categoryCode, ok := categoryCache[record[1]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "category_code",
				Value:   record[1],
				Message: "category_code not found",
			})
			continue
		}

		warehouseCode, ok := warehouseCache[record[2]]
		if !ok {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "warehouse_code",
				Value:   record[2],
				Message: "warehouse_code not found",
			})
			continue
		}

		transactionType := record[4]
		if transactionType != "IN" && transactionType != "OUT" {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "transaction_type",
				Value:   record[4],
				Message: "transaction_type must be 'IN' or 'OUT'",
			})
			continue
		}

		quantity, err := strconv.Atoi(record[3])
		if err != nil {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be a number",
			})
			continue
		}

		if transactionType == "IN" && quantity <= 0 {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be > 0 for IN transaction",
			})
			continue
		}

		if transactionType == "OUT" && quantity >= 0 {
			errors = append(errors, dto.RowError{
				Row:     rowIndex,
				Field:   "quantity",
				Value:   record[3],
				Message: "quantity must be < 0 for OUT transaction",
			})
			continue
		}

		batch = append(batch, model.InventoryTransaction{
			ProductID:       productSku,
			CategoryID:      categoryCode,
			WarehouseID:     warehouseCode,
			Quantity:        quantity,
			TransactionType: transactionType,
		})

		if len(batch) == BatchSize {
			jobs <- batch
			batch = make([]model.InventoryTransaction, 0, BatchSize)
		}
	}

	if len(errors) > 0 {
		close(jobs)
		return &dto.BulkUploadResponse{
			TotalProcessed: rowIndex - 1,
			TotalSuccess:   0,
			Errors:         errors,
		}, nil
	}

	if len(batch) > 0 {
		jobs <- batch
	}

	close(jobs)

	// wait workers
	for i := 0; i < WorkerCount; i++ {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}

	return &dto.BulkUploadResponse{
		TotalProcessed: rowIndex - 1,
		TotalSuccess:   rowIndex - 1,
		Errors:         nil,
	}, nil
}
