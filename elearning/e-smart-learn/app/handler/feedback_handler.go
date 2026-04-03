package handler

import (
	"net/http"
	"strconv"
	"strings"

	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FeedbackHandler interface {
	GetFeedbacks() gin.HandlerFunc
	CreateFeedback() gin.HandlerFunc
	GetFeaturedFeedbacks() gin.HandlerFunc
}

type feedbackHandler struct {
	feedbackService service.FeedbackService
}

func NewFeedbackHandler(feedbackService service.FeedbackService) FeedbackHandler {
	return &feedbackHandler{
		feedbackService: feedbackService,
	}
}

// CreateFeedback godoc
// @Summary Create a new feedback
// @Description Create feedback for a course by current user
// @Tags feedbacks
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param request body dto.CreateFeedbackRequest true "Create Feedback Request"
// @Success 201 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/feedbacks [post]
// @Security BearerAuth
func (h *feedbackHandler) CreateFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateFeedbackRequest
		if err := util.BindAndValidateJSON(c, &req); err != nil {
			_ = c.Error(err)
			return
		}
		courseId := c.Param("id")
		courseUUID, err := uuid.Parse(courseId)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID"))
			return
		}
		req.CourseID = courseUUID
		userId, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		resp, err := h.feedbackService.CreateFeedback(c.Request.Context(), userId, req)
		if err != nil {
			_ = c.Error(err)
			return
		}
		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusCreated)
		response.Request = dto.GetRequestClient(c)
		response.Data = resp
		response.Metadata = nil

		c.JSON(http.StatusCreated, response)
	}
}

// GetFeedbacks godoc
// @Summary Get list of feedbacks
// @Description Retrieve feedback list by course
// @Tags feedbacks
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/courses/{id}/feedbacks [get]
func (h *feedbackHandler) GetFeedbacks() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		courseId, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID"))
			return
		}
		var pagging dto.PagingRequest
		if err := c.ShouldBindQuery(&pagging); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}
		feedbacks, err := h.feedbackService.GetFeedbacks(c.Request.Context(), courseId, pagging.Limit, pagging.Offset)
		if err != nil {
			_ = c.Error(err)
			return
		}
		response := dto.NewApiResponse(c)
		response.Request = dto.GetRequestClient(c)
		response.Data = feedbacks

		c.JSON(http.StatusOK, response)
	}
}

// GetFeaturedFeedbacks godoc
// @Summary Get list of featured feedbacks
// @Description Retrieve list of featured feedbacks
// @Tags feedbacks
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of feedbacks (default: 5, max: 4)"
// @Success 200 {object} any
// @Router /api/v1/feedbacks/featured [get]
func (h *feedbackHandler) GetFeaturedFeedbacks() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "5")
		limit, err := strconv.Atoi(strings.TrimSpace(limitStr))
		if err != nil || limit <= 0 {
			_ = c.Error(apperror.NewBadRequestError("Invalid limit parameter"))
			return
		}
		feedbacks, err := h.feedbackService.GetFeaturedFeedbacks(c.Request.Context(), limit)
		if err != nil {
			_ = c.Error(err)
			return
		}
		response := dto.NewApiResponse(c)
		response.Request = dto.GetRequestClient(c)
		response.Data = feedbacks

		c.JSON(http.StatusOK, response)
	}
}
