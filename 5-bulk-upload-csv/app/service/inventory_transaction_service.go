package service

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"bulk-upload-csv/model"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

type CSVRow struct {
	LineNumber      int
	ProductSku      string
	CategoryCode    string
	WarehouseCode   string
	Quantity        string
	TransactionType string
}

type ProcessedRow struct {
	Transaction *model.InventoryTransaction
	Errors      []dto.RowError
}

func (s *InventoryTransactionService) ProcessBulkUpload(ctx context.Context, fileReader io.Reader) (*dto.BulkUploadResponse, error) {
	reader := csv.NewReader(fileReader)

	// 1. Đọc và validate header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	expectedHeaders := []string{"product_sku", "category_code", "warehouse_code", "quantity", "transaction_type"}
	if len(header) != len(expectedHeaders) {
		return nil, fmt.Errorf("invalid header: expected %d columns, got %d", len(expectedHeaders), len(header))
	}

	for i, h := range header {
		if strings.TrimSpace(h) != expectedHeaders[i] {
			return nil, fmt.Errorf("invalid header at column %d: expected '%s', got '%s'", i+1, expectedHeaders[i], h)
		}
	}

	// 2. Load Master Data vào memory (parallel)
	var (
		productMap   map[string]uuid.UUID
		categoryMap  map[string]uuid.UUID
		warehouseMap map[string]uuid.UUID
		wg           sync.WaitGroup
		mu           sync.Mutex
		loadErrors   []error
	)

	productMap = make(map[string]uuid.UUID)
	categoryMap = make(map[string]uuid.UUID)
	warehouseMap = make(map[string]uuid.UUID)

	wg.Add(3)

	// Load products
	go func() {
		defer wg.Done()
		products, err := s.inventoryTransactionRepository.GetAllProducts(ctx)
		if err != nil {
			mu.Lock()
			loadErrors = append(loadErrors, fmt.Errorf("failed to load products: %w", err))
			mu.Unlock()
			return
		}
		for _, p := range products {
			productMap[p.Sku] = p.ID
		}
	}()

	// Load categories
	go func() {
		defer wg.Done()
		categories, err := s.inventoryTransactionRepository.GetAllCategories(ctx)
		if err != nil {
			mu.Lock()
			loadErrors = append(loadErrors, fmt.Errorf("failed to load categories: %w", err))
			mu.Unlock()
			return
		}
		for _, c := range categories {
			categoryMap[c.Code] = c.ID
		}
	}()

	// Load warehouses
	go func() {
		defer wg.Done()
		warehouses, err := s.inventoryTransactionRepository.GetAllWarehouses(ctx)
		if err != nil {
			mu.Lock()
			loadErrors = append(loadErrors, fmt.Errorf("failed to load warehouses: %w", err))
			mu.Unlock()
			return
		}
		for _, w := range warehouses {
			warehouseMap[w.Code] = w.ID
		}
	}()

	wg.Wait()

	if len(loadErrors) > 0 {
		return nil, fmt.Errorf("failed to load master data: %v", loadErrors)
	}

	// 3. Đọc tất cả CSV rows vào memory
	var csvRows []CSVRow
	lineNumber := 1 // Header là line 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV at line %d: %w", lineNumber+1, err)
		}

		lineNumber++

		if len(record) != 5 {
			return nil, fmt.Errorf("invalid row at line %d: expected 5 columns, got %d", lineNumber, len(record))
		}

		csvRows = append(csvRows, CSVRow{
			LineNumber:      lineNumber,
			ProductSku:      strings.TrimSpace(record[0]),
			CategoryCode:    strings.TrimSpace(record[1]),
			WarehouseCode:   strings.TrimSpace(record[2]),
			Quantity:        strings.TrimSpace(record[3]),
			TransactionType: strings.TrimSpace(record[4]),
		})
	}

	// 4. Process rows với goroutines
	numWorkers := 50 // Tăng số workers
	if len(csvRows) < numWorkers {
		numWorkers = len(csvRows)
	}

	if numWorkers == 0 {
		return &dto.BulkUploadResponse{
			TotalProcessed: 0,
			TotalSuccess:   0,
			Errors:         []dto.RowError{},
		}, nil
	}

	rowChan := make(chan CSVRow, 100)
	resultChan := make(chan ProcessedRow, 100)

	// Start workers
	var workerWg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		workerWg.Add(1)
		go func() {
			defer workerWg.Done()
			s.processWorker(rowChan, resultChan, productMap, categoryMap, warehouseMap)
		}()
	}

	// Send rows to workers trong goroutine riêng
	go func() {
		for _, row := range csvRows {
			rowChan <- row
		}
		close(rowChan)
	}()

	// Collect results trong goroutine riêng
	var (
		txsToInsert []model.InventoryTransaction
		allErrors   []dto.RowError
		resultMu    sync.Mutex
	)

	go func() {
		for result := range resultChan {
			resultMu.Lock()
			if len(result.Errors) > 0 {
				allErrors = append(allErrors, result.Errors...)
			} else if result.Transaction != nil {
				txsToInsert = append(txsToInsert, *result.Transaction)
			}
			resultMu.Unlock()
		}
	}()

	// Wait for all workers to finish
	workerWg.Wait()
	close(resultChan)

	// Đợi result collector goroutine hoàn thành
	// Sleep ngắn để đảm bảo tất cả results đã được collect
	time.Sleep(100 * time.Millisecond)

	// 5. Batch insert parallel với transactions
	totalSuccess := 0
	if len(txsToInsert) > 0 {
		batchSize := 5000 // Tăng batch size
		numBatches := (len(txsToInsert) + batchSize - 1) / batchSize

		type BatchResult struct {
			Success int
			Error   error
		}

		batchResults := make(chan BatchResult, numBatches)
		var batchWg sync.WaitGroup

		// Insert batches parallel
		for i := 0; i < len(txsToInsert); i += batchSize {
			end := i + batchSize
			if end > len(txsToInsert) {
				end = len(txsToInsert)
			}

			batch := txsToInsert[i:end]
			batchWg.Add(1)

			go func(batchData []model.InventoryTransaction) {
				defer batchWg.Done()

				// Sử dụng transaction
				tx := s.inventoryTransactionRepository.BeginTx(ctx)
				err := s.inventoryTransactionRepository.BatchInsertTransactionsWithTx(ctx, tx, batchData)
				if err != nil {
					tx.Rollback()
					batchResults <- BatchResult{Success: 0, Error: err}
					return
				}

				if err := tx.Commit().Error; err != nil {
					tx.Rollback()
					batchResults <- BatchResult{Success: 0, Error: err}
					return
				}

				batchResults <- BatchResult{Success: len(batchData), Error: nil}
			}(batch)
		}

		batchWg.Wait()
		close(batchResults)

		// Collect batch results
		var insertErrors []error
		for result := range batchResults {
			if result.Error != nil {
				insertErrors = append(insertErrors, result.Error)
			} else {
				totalSuccess += result.Success
			}
		}

		if len(insertErrors) > 0 {
			return nil, fmt.Errorf("failed to insert some batches: %v", insertErrors)
		}
	}

	return &dto.BulkUploadResponse{
		TotalProcessed: len(csvRows),
		TotalSuccess:   totalSuccess,
		Errors:         allErrors,
	}, nil
}

func (s *InventoryTransactionService) processWorker(
	rowChan <-chan CSVRow,
	resultChan chan<- ProcessedRow,
	productMap map[string]uuid.UUID,
	categoryMap map[string]uuid.UUID,
	warehouseMap map[string]uuid.UUID,
) {
	for row := range rowChan {
		result := s.processRow(row, productMap, categoryMap, warehouseMap)
		resultChan <- result
	}
}

func (s *InventoryTransactionService) processRow(
	row CSVRow,
	productMap map[string]uuid.UUID,
	categoryMap map[string]uuid.UUID,
	warehouseMap map[string]uuid.UUID,
) ProcessedRow {
	var rowErrors []dto.RowError

	// Validate and lookup IDs
	productID, productOk := productMap[row.ProductSku]
	if !productOk {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "product_sku",
			Value:   row.ProductSku,
			Message: "product not found",
		})
	}

	categoryID, categoryOk := categoryMap[row.CategoryCode]
	if !categoryOk {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "category_code",
			Value:   row.CategoryCode,
			Message: "category not found",
		})
	}

	warehouseID, warehouseOk := warehouseMap[row.WarehouseCode]
	if !warehouseOk {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "warehouse_code",
			Value:   row.WarehouseCode,
			Message: "warehouse not found",
		})
	}

	// Validate quantity
	quantity, qtyErr := strconv.Atoi(row.Quantity)
	if qtyErr != nil {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "quantity",
			Value:   row.Quantity,
			Message: "quantity must be a valid integer",
		})
	}

	// Validate transaction type
	txType := strings.ToUpper(row.TransactionType)
	if txType != "IN" && txType != "OUT" {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "transaction_type",
			Value:   row.TransactionType,
			Message: "transaction_type must be IN or OUT",
		})
	}

	// Business rule: quantity validation based on transaction type
	if qtyErr == nil && txType == "IN" && quantity < 0 {
		rowErrors = append(rowErrors, dto.RowError{
			Row:     row.LineNumber,
			Field:   "quantity",
			Value:   row.Quantity,
			Message: "quantity must be >= 0 for IN transaction",
		})
	}

	// If there are any errors, return them
	if len(rowErrors) > 0 {
		return ProcessedRow{
			Transaction: nil,
			Errors:      rowErrors,
		}
	}

	// Create transaction
	transaction := &model.InventoryTransaction{
		ProductID:       productID,
		CategoryID:      categoryID,
		WarehouseID:     warehouseID,
		Quantity:        quantity,
		TransactionType: txType,
	}

	return ProcessedRow{
		Transaction: transaction,
		Errors:      nil,
	}
}
