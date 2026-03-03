package handler

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService interfaces.ProductServiceInterface
}

func NewProductHandler(
	productService interfaces.ProductServiceInterface,
) interfaces.ProductHandlerInterface {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetListProductRequest
		var errs []error

		// ShouldBindQuery để bind dữ liệu từ URL Query String (dùng với tag 'form' trong struct)
		if err := c.ShouldBindQuery(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.productService.GetList(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *ProductHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetProductDetailRequest
		var errs []error

		// ShouldBind để bind dữ liệu từ URL Path Param (dùng với tag 'uri' trong struct)
		if err := c.ShouldBindUri(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.productService.GetDetail(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *ProductHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.CreateProductRequest
		var errs []error
		if err := c.ShouldBind(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		data, errs := h.productService.Create(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}
