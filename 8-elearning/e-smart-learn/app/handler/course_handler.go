package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CourseHandler handles course listing and detail endpoints (mocked).
type CourseHandler interface {
	GetCourses() gin.HandlerFunc
	GetCourseBySlug() gin.HandlerFunc
}

type courseHandler struct{}

// NewCourseHandler creates a new CourseHandler.
func NewCourseHandler() CourseHandler { return &courseHandler{} }

// GetCourses godoc
// @Summary View course list
// @Description Return a paginated list of courses (mocked)
// @Tags courses
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param sortBy query string false "sortBy"
// @Param sortOrder query string false "sortOrder"
// @Param categoryId query string false "categoryId"
// @Param type query string false "type"
// @Param search query string false "search"
// @Success 200 {object} any
// @Router /api/v1/courses [get]
func (h *courseHandler) GetCourses() gin.HandlerFunc {
	return func(c *gin.Context) {
		// read query params with defaults
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		sortBy := c.DefaultQuery("sortBy", "created_at")
		sortOrder := c.DefaultQuery("sortOrder", "asc")
		categoryId := c.DefaultQuery("categoryId", "")
		typ := c.DefaultQuery("type", "")
		search := c.DefaultQuery("search", "")

		response := gin.H{
			"process_id": "c1b2a8f4-7d91-4f0e-9b2c-6f0f1c92e1aa",
			"path":       "/api/v1/courses",
			"status":     gin.H{"code": 200, "type": "OK"},
			"request": gin.H{
				"limit":      limit,
				"offset":     (page - 1) * limit,
				"sortBy":     sortBy,
				"sortOrder":  sortOrder,
				"categoryId": categoryId,
				"type":       typ,
				"search":     search,
			},
			"errors": []any{},
			"data": []gin.H{
				{
					"id":                "course_001",
					"title":             "Golang Backend Development",
					"short_description": "Learn how to build scalable backend services using Go.",
					"image_url":         "https://cdn.example.com/courses/go-backend.jpg",
					"average_rate":      4.7,
					"old_price":         99,
					"price":             49,
					"category":          gin.H{"id": "cat_backend", "title": "Backend Development"},
					"instructor":        gin.H{"id": "ins_001", "fullname": "Nguyen Van A"},
					"created_at":        "2024-01-10T10:00:00Z",
					"updated_at":        "2024-01-15T10:00:00Z",
				},
				{
					"id":                "course_002",
					"title":             "Microservices with Go",
					"short_description": "Design and implement microservices architecture using Go.",
					"image_url":         "https://cdn.example.com/courses/go-microservices.jpg",
					"average_rate":      4.8,
					"old_price":         120,
					"price":             59,
					"category":          gin.H{"id": "cat_backend", "title": "Backend Development"},
					"instructor":        gin.H{"id": "ins_002", "fullname": "Tran Van B"},
					"created_at":        "2024-02-01T10:00:00Z",
					"updated_at":        "2024-02-05T10:00:00Z",
				},
			},
			"metadata": gin.H{"limit": limit, "offset": (page - 1) * limit, "sortBy": sortBy, "sortOrder": sortOrder, "total": 125},
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetCourseBySlug godoc
// @Summary View course detail
// @Description Return detailed information of a course by slug (mocked)
// @Tags courses
// @Accept json
// @Produce json
// @Param slug path string true "course slug"
// @Success 200 {object} any
// @Router /api/v1/courses/{slug} [get]
func (h *courseHandler) GetCourseBySlug() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		response := gin.H{
			"process_id": "c1b2a8f4-7d91-4f0e-9b2c-6f0f1c92e1aa",
			"path":       "/api/v1/courses/" + slug,
			"status":     gin.H{"code": 200, "type": "OK"},
			"request":    gin.H{},
			"errors":     []any{},
			"data": gin.H{
				"id":           "course_001",
				"title":        "Golang Backend Development",
				"description":  "A complete course covering Go fundamentals, REST APIs, concurrency, and production-ready backend systems.",
				"image_url":    "https://cdn.example.com/courses/go-backend.jpg",
				"instructor":   gin.H{"id": "ins_001", "fullname": "Nguyen Van A"},
				"average_rate": 4.7,
				"old_price":    99,
				"price":        49,
				"category":     gin.H{"id": "cat_backend", "title": "Backend Development"},
				"created_at":   "2024-01-10T10:00:00Z",
				"updated_at":   "2024-01-15T10:00:00Z",
			},
			"metadata": nil,
		}

		c.JSON(http.StatusOK, response)
	}
}
