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

func (h *userHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateUserRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		data, err := h.userService.Create(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		resp := dto.Success(data, "User created successfully")

		c.JSON(http.StatusCreated, resp)
	}
}

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

		c.JSON(http.StatusOK, dto.Success(data, "User updated successfully"))
	}
}

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

		c.JSON(http.StatusOK, dto.Success[any](nil, "User deleted successfully"))
	}
}

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

		c.JSON(http.StatusOK, dto.Success(data, "Get user detail successfully"))
	}
}

// Test API: http://localhost:8080/api/v1/users?page=1&limit=5
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

		page := paginationRequest.Page
		limit := paginationRequest.Limit
		offset := paginationRequest.GetOffset()

		// Call service
		data, total, err := h.userService.GetList(
			c.Request.Context(),
			limit,
			offset,
		)

		if err != nil {
			c.Error(err)
			return
		}

		resp := dto.Paginated(
			&data,
			dto.NewPagination(page, limit, int(total)),
			"Get users successfully",
		)

		c.JSON(http.StatusOK, resp)
	}
}

// Test API: http://localhost:8080/api/v1/users/filter?page=1&limit=5&name=dong2&sort=username:desc,created_at:asc
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

		page := request.Page
		limit := request.Limit

		// Call service
		data, total, err := h.userService.FilterAndPaginateAndSort(
			c.Request.Context(),
			request,
		)

		if err != nil {
			c.Error(err)
			return
		}

		resp := dto.Paginated(
			&data,
			dto.NewPagination(page, limit, int(total)),
			"Get users successfully",
		)

		c.JSON(http.StatusOK, resp)
	}
}
