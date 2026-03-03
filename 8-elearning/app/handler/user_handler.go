package handler

import (
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetDetail() gin.HandlerFunc
	Create() gin.HandlerFunc
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

func (h *userHandler) GetDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.GetUserDetailRequest

		if err := c.ShouldBindUri(&request); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid request parameters"))
			return
		}

		data, err := h.userService.GetDetail(c.Request.Context(), request)
		if err != nil {
			c.Error(err) // để ErrorHandler xử lý
			return
		}

		resp := dto.NewSuccessResponse(
			&data,
			"",
			util.GetRequestID(c),
		)

		c.JSON(http.StatusOK, resp)
	}
}

func (h *userHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateUserRequest

		if err := c.ShouldBind(&request); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid request body"))
			return
		}

		data, err := h.userService.Create(c.Request.Context(), request)
		if err != nil {
			c.Error(err)
			return
		}

		resp := dto.NewSuccessResponse(
			&data,
			"User created successfully",
			util.GetRequestID(c),
		)

		c.JSON(http.StatusCreated, resp)
	}
}
