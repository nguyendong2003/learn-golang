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

type LessonHandler interface {
	Create() gin.HandlerFunc
	UpdateLessons() gin.HandlerFunc
	GetByCourseID() gin.HandlerFunc
}

type lessonHandler struct {
	lessonService service.LessonService
}

func NewLessonHandler(
	lessonService service.LessonService,
) LessonHandler {
	return &lessonHandler{
		lessonService: lessonService,
	}
}

// Create godoc
// @Summary Create a chapter with lessons
// @Description Create one chapter and multiple lessons in a single request using presigned media URLs
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.CreateLessonsRequest true "Chapters and lessons create request"
// @Success 201 {object} dto.ApiResponse{data=[]dto.ChapterResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request payload"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/lessons [post]
// @Security BearerAuth
func (h *lessonHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateLessonsRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		id := c.Param("id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid user ID format"))
			return
		}
		data, err := h.lessonService.Create(c.Request.Context(), userID, courseID, request.Chapters)
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

// UpdateLessons godoc
// @Summary Update course with all chapters and lessons
// @Description Update course structure with chapters and lessons (supports create, update, delete)
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.UpdateCourseWithChaptersRequest true "Course chapters and lessons update request"
// @Success 200 {object} dto.ApiResponse{data=[]dto.ChapterResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request"
// @Failure 403 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/lessons [put]
// @Security BearerAuth
func (h *lessonHandler) UpdateLessons() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.UpdateCourseWithChaptersRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		id := c.Param("id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid user ID format"))
			return
		}
		data, err := h.lessonService.UpdateCourseWithChapters(c.Request.Context(), userID, courseID, request)
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

// GetByCourseID godoc
// @Summary Get all chapters and lessons of a course
// @Description Retrieve the full structure of a course including chapters and their associated lessons
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=[]dto.ChapterResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid course ID"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id}/lessons [get]
func (h *lessonHandler) GetByCourseID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		data, err := h.lessonService.GetByCourseID(c.Request.Context(), courseID)
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
