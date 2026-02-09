package interfaces

import (
	"bulk-upload-csv/dto"
	"context"
	"io"
)

type InventoryTransactionServiceInterface interface {
	ProcessBulkUpload(ctx context.Context, fileReader io.Reader) (*dto.BulkUploadResponse, error)
}
