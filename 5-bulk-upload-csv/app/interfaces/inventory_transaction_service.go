package interfaces

import (
	"bulk-upload-csv/dto"
	"context"
	"io"
)

type InventoryTransactionServiceInterface interface {
	ProcessBulkUpload(ctx context.Context, file io.Reader) (*dto.BulkUploadResponse, error)
	ProcessBulkUploadGoroutine(ctx context.Context, file io.Reader) (*dto.BulkUploadResponse, error)
}
