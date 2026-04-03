package handler

import (
	"net/http"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FollowHandler interface {
	FollowUser() gin.HandlerFunc
	UnfollowUser() gin.HandlerFunc
	GetFollowers() gin.HandlerFunc
	GetFollowings() gin.HandlerFunc
}

type followHandler struct {
	followService service.FollowService
}

func NewFollowHandler(followService service.FollowService) FollowHandler {
	return &followHandler{followService: followService}
}

// FollowUser godoc
// @Summary Follow a user
// @Description Follow another user by ID
// @Tags follows
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 201 {object} dto.ApiResponse
// @Router /api/v1/users/{id}/follow [post]
func (h *followHandler) FollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetID, err := getTargetUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		requestUserID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := h.followService.FollowUser(c.Request.Context(), requestUserID, targetID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusCreated, res)
	}
}

// UnfollowUser godoc
// @Summary Unfollow a user
// @Description Unfollow another user by ID
// @Tags follows
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/{id}/follow [delete]
func (h *followHandler) UnfollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetID, err := getTargetUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		requestUserID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := h.followService.UnfollowUser(c.Request.Context(), requestUserID, targetID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetFollowers godoc
// @Summary Get followers by user ID
// @Description Get paginated list of followers of a user
// @Tags follows
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param sortBy query string false "sortBy"
// @Param sortOrder query string false "sortOrder"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/{id}/followers [get]
func (h *followHandler) GetFollowers() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetID, err := getTargetUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.PagingRequest
		if err := util.BindAndValidateQuery(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		request.Process()

		data, total, err := h.followService.GetFollowers(c.Request.Context(), targetID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// GetFollowings godoc
// @Summary Get followings by user ID
// @Description Get paginated list of users followed by a user
// @Tags follows
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param sortBy query string false "sortBy"
// @Param sortOrder query string false "sortOrder"
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/users/{id}/followings [get]
func (h *followHandler) GetFollowings() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetID, err := getTargetUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.PagingRequest
		if err := util.BindAndValidateQuery(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		request.Process()

		data, total, err := h.followService.GetFollowings(c.Request.Context(), targetID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}

func getTargetUserID(c *gin.Context) (uuid.UUID, error) {
	id := c.Param("id")
	targetID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperror.NewBadRequestError("Invalid user ID")
	}

	return targetID, nil
}
