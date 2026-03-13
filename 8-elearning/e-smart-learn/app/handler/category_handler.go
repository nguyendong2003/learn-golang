package handler

import (
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CategoryHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
	GetAll() gin.HandlerFunc
	GetList() gin.HandlerFunc
}

type categoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(
	categoryService service.CategoryService,
) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}

// Create godoc
// @Summary Create a new category
// @Description Create a new category with the provided name and description. Category name must be unique.
// @Tags categories
// @Accept json
// @Produce json
// @Param payload body dto.CreateCategoryRequest true "Category create request"
// @Success 201 {object} dto.ApiResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request payload or validation failed"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories [post]
// @Security BearerAuth
func (h *categoryHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateCategoryRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		data, err := h.categoryService.Create(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusCreated, res)
	}
}

// Update godoc
// @Summary Update an existing category
// @Description Update a category by ID with the provided name and description. Category name must be unique.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Param payload body dto.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} dto.ApiResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format or validation failed"
// @Failure 404 {object} dto.ApiResponse "Category not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories/{id} [put]
// @Security BearerAuth
func (h *categoryHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Bind & validate JSON body
		var request dto.UpdateCategoryRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		// Call service
		data, err := h.categoryService.Update(c.Request.Context(), id, request)

		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// Delete godoc
// @Summary Delete a category
// @Description Delete a category by ID. Cannot delete a category that has associated courses.
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "Category not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories/{id} [delete]
// @Security BearerAuth
func (h *categoryHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get category by ID
// @Description Retrieve a single category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CategoryResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "Category not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories/{id} [get]
func (h *categoryHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		data, err := h.categoryService.GetByID(c.Request.Context(), id)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetAll godoc
// @Summary Get all categories
// @Description Retrieve all categories without pagination
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse{data=[]dto.CategoryResponse}
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories/all [get]
func (h *categoryHandler) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := h.categoryService.GetAll(c.Request.Context())

		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetList godoc
// @Summary Get paginated list of categories
// @Description Retrieve a paginated list of categories with filtering, sorting, and pagination support
// @Tags categories
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10, max: 100)" default(10)
// @Param offset query int false "Number of items to skip (default: 0)" default(0)
// @Param sortBy query string false "Field to sort by (default: created_at)" default(created_at)
// @Param sortOrder query string false "Sort order: asc or desc (default: desc)" default(desc) Enums(asc,desc)
// @Param name query string false "Filter categories by name (partial match)"
// @Success 200 {object} dto.ApiResponse{data=[]dto.CategoryResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse "Invalid query parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/categories [get]
func (h *categoryHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListCategoryRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		// Process default pagination
		request.Process()

		limit := request.Limit
		offset := request.Offset
		sortBy := request.SortBy
		sortOrder := request.SortOrder

		// Call service
		data, total, err := h.categoryService.GetList(c.Request.Context(), request)

		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}
