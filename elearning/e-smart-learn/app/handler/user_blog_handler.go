package handler

import (
	"net/http"

	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type UserBlogHandler interface {
	MarkBlogAsRead() gin.HandlerFunc
	GetReadBlogs() gin.HandlerFunc
}

type userBlogHandler struct {
	userService service.UserService
}

func NewUserBlogHandler(userService service.UserService) UserBlogHandler {
	return &userBlogHandler{
		userService: userService,
	}
}

// MarkBlogAsRead godoc
// @Summary Mark blog as read
// @Description Mark a blog as read for the current user
// @Tags user-blogs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body dto.UpdateReadBlogHistoryRequest true "course progress payload"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/me/blogs/read [post]
func (h *userBlogHandler) MarkBlogAsRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.UpdateReadBlogHistoryRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		updatedUser, err := h.userService.MarkBlogAsRead(
			c.Request.Context(),
			userID,
			request,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Data = updatedUser
		c.JSON(http.StatusOK, response)
	}
}

// GetReadBlogs godoc
// @Summary Get read blogs
// @Description Get a list of blogs read by the current user
// @Tags user-blogs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param params query dto.ViewReadBlogHistoryRequest false "query parameters"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/me/blogs/read [get]
func (h *userBlogHandler) GetReadBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ViewReadBlogHistoryRequest
		if err := util.BindAndValidateQuery(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		request.Process()
		readBlogsResponse, err := h.userService.GetReadBlogs(
			c.Request.Context(),
			userID,
			request,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Request = request
		response.Data = readBlogsResponse
		response.Metadata = dto.NewPagination(request.Limit, request.Offset, readBlogsResponse.Total, request.SortBy, request.SortOrder)
		c.JSON(http.StatusOK, response)
	}
}
