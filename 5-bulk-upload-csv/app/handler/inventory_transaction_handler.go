package handler

import (
	"bulk-upload-csv/interfaces"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type InventoryTransactionHandler struct {
	inventoryTransactionService interfaces.InventoryTransactionServiceInterface
}

func NewInventoryTransactionHandler(
	inventoryTransactionService interfaces.InventoryTransactionServiceInterface,
) interfaces.InventoryTransactionHandlerInterface {
	return &InventoryTransactionHandler{
		inventoryTransactionService: inventoryTransactionService,
	}
}

func (h *InventoryTransactionHandler) ProcessBulkUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Nhận file từ form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Giới hạn định dạng file
		if !strings.HasSuffix(file.Filename, ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer openedFile.Close()

		// Gọi service xử lý
		result, err := h.inventoryTransactionService.ProcessBulkUpload(c.Request.Context(), openedFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func (h *InventoryTransactionHandler) ProcessBulkUploadGoroutine() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Nhận file từ form-data
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
			return
		}

		// Giới hạn định dạng file
		if !strings.HasSuffix(file.Filename, ".csv") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files are allowed"})
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
			return
		}
		defer openedFile.Close()

		// Gọi service xử lý
		result, err := h.inventoryTransactionService.ProcessBulkUploadGoroutine(c.Request.Context(), openedFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
