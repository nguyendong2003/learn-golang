package handler

import (
	"net/http"
	"strconv"

	"elearning-api/dto"
	"elearning-api/model"

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
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		// create mock models
		m1 := &model.Course{
			Title:       "Golang Backend Development",
			Description: "Learn how to build scalable backend services using Go.",
			Image:       "https://cdn.example.com/courses/go-backend.jpg",
			Slug:        "course_001",
			AverageRate: 4.7,
			OldPrice:    99,
			Price:       49,
		}
		m1.Category = &model.Category{Name: "Backend Development"}
		m1.Instructor = &model.InstructorProfile{}
		m1.Instructor.User = &model.User{Name: "Nguyen Van A"}

		m2 := &model.Course{
			Title:       "Microservices with Go",
			Description: "Design and implement microservices architecture using Go.",
			Image:       "https://cdn.example.com/courses/go-microservices.jpg",
			Slug:        "course_002",
			AverageRate: 4.8,
			OldPrice:    120,
			Price:       59,
		}
		m2.Category = &model.Category{Name: "Backend Development"}
		m2.Instructor = &model.InstructorProfile{}
		m2.Instructor.User = &model.User{Name: "Tran Van B"}

		list := dto.NewListCourseResponse([]*model.Course{m1, m2})

		resp := dto.NewApiResponse(c)
		resp.Request = gin.H{"limit": limit, "offset": (page - 1) * limit, "sortBy": c.DefaultQuery("sortBy", "created_at"), "sortOrder": c.DefaultQuery("sortOrder", "asc"), "categoryId": c.DefaultQuery("categoryId", ""), "type": c.DefaultQuery("type", ""), "search": c.DefaultQuery("search", "")}
		resp.Data = list
		resp.Metadata = dto.NewPagination(limit, (page-1)*limit, 125, c.DefaultQuery("sortBy", "created_at"), c.DefaultQuery("sortOrder", "asc"))

		c.JSON(http.StatusOK, resp)
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

		m := &model.Course{
			Title:       "Golang Backend Development",
			Description: "A complete course covering Go fundamentals, REST APIs, concurrency, and production-ready backend systems.",
			Image:       "https://cdn.example.com/courses/go-backend.jpg",
			Slug:        "course_001",
			AverageRate: 4.7,
			OldPrice:    99,
			Price:       49,
		}
		m.Category = &model.Category{Name: "Backend Development"}
		m.Instructor = &model.InstructorProfile{}
		m.Instructor.User = &model.User{Name: "Nguyen Van A"}

		detail := dto.NewCourseDetailResponse(m)

		resp := dto.NewApiResponse(c)
		resp.Request = gin.H{}
		resp.Data = detail
		resp.Metadata = nil

		// ensure path includes slug
		resp.Path = "/api/v1/courses/" + slug

		c.JSON(http.StatusOK, resp)
	}
}
