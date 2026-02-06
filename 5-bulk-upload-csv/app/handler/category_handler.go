package handler

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService interfaces.CategoryServiceInterface
}

func NewCategoryHandler(
	categoryService interfaces.CategoryServiceInterface,
) interfaces.CategoryHandlerInterface {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetListCategoryRequest
		var errs []error

		if err := c.ShouldBind(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.categoryService.GetList(c.Request.Context())
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *CategoryHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.GetCategoryDetailRequest
		var errs []error

		if err := c.ShouldBind(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		response.Request = request

		data, errs := h.categoryService.GetDetail(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}

func (h *CategoryHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c.FullPath())

		var request dto.CreateCategoryRequest
		var errs []error
		if err := c.ShouldBind(&request); err != nil {
			errs = append(errs, err)
			errorResponse(c, response, errs)
			return
		}

		data, errs := h.categoryService.Create(c.Request.Context(), request)
		if errs != nil {
			errorResponse(c, response, errs)
			return
		}

		response.Data = data

		c.JSON(http.StatusOK, response)
	}
}
