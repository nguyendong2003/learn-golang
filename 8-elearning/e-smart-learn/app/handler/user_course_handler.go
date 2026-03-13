package handler

import (
	"net/http"

	"elearning-api/dto"
	"elearning-api/model"

	"github.com/gin-gonic/gin"
)

// UserCourseHandler handles user-specific course endpoints (mocked).
type UserCourseHandler interface {
	UpdateCourseProgress() gin.HandlerFunc
	GetMyCourses() gin.HandlerFunc
}

type userCourseHandler struct{}

// NewUserCourseHandler creates a new UserCourseHandler.
func NewUserCourseHandler() UserCourseHandler { return &userCourseHandler{} }

// UpdateCourseProgress godoc
// @Summary Update course progress
// @Description Update progress for a lesson the user completed (mocked)
// @Tags user-courses
// @Accept json
// @Produce json
// @Param payload body object true "course progress payload"
// @Success 200 {object} any
// @Router /api/v1/users/me/courses/progress [post]
func (h *userCourseHandler) UpdateCourseProgress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload struct {
			CourseId string `json:"courseId"`
			LessonId string `json:"lessonId"`
		}
		_ = c.BindJSON(&payload)

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/users/me/courses/progress"
		resp.Request = gin.H{"courseId": payload.CourseId, "lessonId": payload.LessonId}
		resp.Data = gin.H{"lesson_id": payload.LessonId, "completed": true}
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}

// GetMyCourses godoc
// @Summary Track user course progress
// @Description Return list of enrolled courses with progress (mocked)
// @Tags user-courses
// @Accept json
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param sortBy query string false "sortBy"
// @Param sortOrder query string false "sortOrder"
// @Success 200 {object} any
// @Router /api/v1/users/me/courses [get]
func (h *userCourseHandler) GetMyCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 10
		offset := 0

		// mock course model
		m := &model.Course{
			Title:       "Golang Backend Development",
			Description: "Learn how to build scalable backend services using Go.",
			Image:       "https://cdn.example.com/courses/go-backend.jpg",
			Slug:        "course_001",
		}
		m.Category = &model.Category{Name: "Backend Development"}
		m.InstructorProfile = &model.InstructorProfile{}
		m.InstructorProfile.User = &model.User{Name: "Nguyen Van A"}

		courseDTO := dto.NewCourseDetailResponse(m)

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/users/me/courses"
		resp.Request = gin.H{"limit": limit, "offset": offset, "sortBy": c.DefaultQuery("sortBy", "created_at"), "sortOrder": c.DefaultQuery("sortOrder", "desc")}
		resp.Data = []gin.H{
			{
				"course":      courseDTO,
				"progress":    gin.H{"percentage": 65, "completed_lessons": 13, "total_lessons": 20, "last_accessed_at": "2024-03-01T10:00:00Z"},
				"enrolled_at": "2024-02-10T10:00:00Z",
				"updated_at":  "2024-03-01T10:00:00Z",
			},
		}
		resp.Metadata = gin.H{"limit": limit, "offset": offset, "total": 125}

		c.JSON(http.StatusOK, resp)
	}
}
