package handler

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WarehouseHandler struct {
	warehouseService interfaces.WarehouseServiceInterface
}

func NewWarehouseHandler(
	warehouseService interfaces.WarehouseServiceInterface,
) interfaces.WarehouseHandlerInterface {
	return &WarehouseHandler{
		warehouseService: warehouseService,
	}
}

func (h *WarehouseHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetListWarehouseRequest
		var errs []error

		// ShouldBindQuery để bind dữ liệu từ URL Query String (dùng với tag 'form' trong struct)
		if err := c.ShouldBindQuery(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.warehouseService.GetList(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *WarehouseHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetWarehouseDetailRequest
		var errs []error

		// ShouldBind để bind dữ liệu từ URL Path Param (dùng với tag 'uri' trong struct)
		if err := c.ShouldBindUri(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.warehouseService.GetDetail(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *WarehouseHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.CreateWarehouseRequest
		var errs []error
		if err := c.ShouldBind(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		data, errs := h.warehouseService.Create(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}
