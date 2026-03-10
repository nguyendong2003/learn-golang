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

type UserHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
	GetList() gin.HandlerFunc
	FilterAndPaginateAndSort() gin.HandlerFunc
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(
	userService service.UserService,
) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param payload body dto.CreateUserRequest true "Create user payload"
// @Success 201 {object} any
// @Failure 400 {object} any
// @Router /api/v1/users [post]
func (h *userHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateUserRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		user, err := h.userService.Create(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = user

		c.JSON(http.StatusCreated, res)
	}
}

// UpdateUser godoc
// @Summary Update user
// @Description Update an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param payload body dto.UpdateUserRequest true "Update payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Router /api/v1/users/{id} [put]
func (h *userHandler) Update() gin.HandlerFunc {
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
		var request dto.UpdateUserRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		// Call service
		data, err := h.userService.Update(
			c.Request.Context(),
			id,
			request,
		)

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

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Router /api/v1/users/{id} [delete]
func (h *userHandler) Delete() gin.HandlerFunc {
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
		if err := h.userService.DeleteByID(c.Request.Context(), id); err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Router /api/v1/users/{id} [get]
func (h *userHandler) GetByID() gin.HandlerFunc {
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
		data, err := h.userService.GetByID(c.Request.Context(), id)
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

// Test API: http://localhost:8080/api/v1/users?page=1&limit=5
// ListUsers godoc
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/users [get]
func (h *userHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var paginationRequest dto.PagingRequest

		// Bind query params
		if err := c.ShouldBindQuery(&paginationRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid pagination parameters"))
			return
		}

		// Process default values
		paginationRequest.Process()

		limit := paginationRequest.Limit
		offset := paginationRequest.Offset
		sortBy := paginationRequest.SortBy
		sortOrder := paginationRequest.SortOrder

		// Call service
		users, total, err := h.userService.GetList(
			c.Request.Context(),
			limit,
			offset,
		)

		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = users
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// Test API: http://localhost:8080/api/v1/users/filter?page=1&limit=5&name=dong2&sort=username:desc,created_at:asc
// FilterUsers godoc
// @Summary Filter and paginate users
// @Description Filter users by username/name and paginate/sort
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Param username query string false "Username filter"
// @Param name query string false "Name filter"
// @Param sort query string false "Sort e.g. username:desc,created_at:asc"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/users/filter [get]
func (h *userHandler) FilterAndPaginateAndSort() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.FilterUserRequest

		// Bind query params (pagination + filter + sort)
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
		data, total, err := h.userService.FilterAndPaginateAndSort(
			c.Request.Context(),
			request,
		)

		if err != nil {
			c.Error(err)
			return
		}

		// Create request object with filter params
		requestObj := map[string]any{
			"limit":     limit,
			"offset":    offset,
			"sortBy":    sortBy,
			"sortOrder": sortOrder,
		}
		if request.Username != nil {
			requestObj["username"] = *request.Username
		}
		if request.Name != nil {
			requestObj["name"] = *request.Name
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}
