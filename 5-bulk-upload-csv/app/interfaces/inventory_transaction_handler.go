package interfaces

import "github.com/gin-gonic/gin"

type InventoryTransactionHandlerInterface interface {
	ProcessBulkUpload() gin.HandlerFunc
	ProcessBulkUploadGoroutine() gin.HandlerFunc
}
