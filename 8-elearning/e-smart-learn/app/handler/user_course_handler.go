package handler

import (
	"net/http"

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

		response := gin.H{
			"process_id": "c1b2a8f4-7d91-4f0e-9b2c-6f0f1c92e1aa",
			"path":       "/api/v1/users/me/courses/progress",
			"status":     gin.H{"code": 200, "type": "OK"},
			"request":    gin.H{"courseId": payload.CourseId, "lessonId": payload.LessonId},
			"errors":     []any{},
			"data":       gin.H{"lesson_id": payload.LessonId, "completed": true},
			"metadata":   nil,
		}

		c.JSON(http.StatusOK, response)
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
		sortBy := c.DefaultQuery("sortBy", "created_at")
		sortOrder := c.DefaultQuery("sortOrder", "desc")

		response := gin.H{
			"process_id": "c1b2a8f4-7d91-4f0e-9b2c-6f0f1c92e1aa",
			"path":       "/api/v1/users/me/courses",
			"status":     gin.H{"code": 200, "type": "OK"},
			"request":    gin.H{"limit": limit, "offset": offset, "sortBy": sortBy, "sortOrder": sortOrder},
			"errors":     []any{},
			"data": []gin.H{
				{
					"course": gin.H{
						"id":                "course_001",
						"title":             "Golang Backend Development",
						"short_description": "Learn how to build scalable backend services using Go.",
						"image_url":         "https://cdn.example.com/courses/go-backend.jpg",
						"category":          gin.H{"id": "cat_backend", "title": "Backend Development"},
						"instructor":        gin.H{"id": "ins_001", "fullname": "Nguyen Van A"},
					},
					"progress":    gin.H{"percentage": 65, "completed_lessons": 13, "total_lessons": 20, "last_accessed_at": "2024-03-01T10:00:00Z"},
					"enrolled_at": "2024-02-10T10:00:00Z",
					"updated_at":  "2024-03-01T10:00:00Z",
				},
			},
			"metadata": gin.H{"limit": limit, "offset": offset, "total": 125},
		}

		c.JSON(http.StatusOK, response)
	}
}
