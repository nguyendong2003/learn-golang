package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BlogHandler interface {
	GetByID() gin.HandlerFunc
	GetList() gin.HandlerFunc
}

type blogHandler struct{}

func NewBlogHandler() BlogHandler {
	return &blogHandler{}
}

// GetList godoc
// @Summary Get list of blogs
// @Description Retrieve paginated list of blogs with optional filters
// @Tags blogs
// @Accept json
// @Produce json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Param sortBy query string false "sortBy" default(created_at)
// @Param sortOrder query string false "sortOrder" Enums(asc,desc) default(asc)
// @Param categoryId query string false "categoryId"
// @Param type query string false "type"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/blogs [get]
func (h *blogHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := gin.H{
			"process_id": "c1b2a8f4-7d91-4f0e-9b2c-6f0f1c92e1aa",
			"path":       "/api/v1/blogs",
			"status": gin.H{
				"code": 200,
				"type": "OK",
			},
			"request": gin.H{
				"limit":      10,
				"offset":     0,
				"sortBy":     "created_at",
				"sortOrder":  "asc",
				"categoryId": "",
				"type":       "",
			},
			"errors": []any{},
			"data": []gin.H{
				{
					"id":        "b1d2c3a4",
					"title":     "Introduction to Golang for Backend Development",
					"content":   "Golang is a powerful language designed for building scalable backend systems.",
					"image_url": "https://cdn.example.com/blogs/golang-intro.jpg",
					"author": gin.H{
						"id":         "u1001",
						"name":       "John Nguyen",
						"avatar_url": "https://cdn.example.com/avatars/john.jpg",
					},
					"view_count": 1204,
					"created_at": "2025-12-01T10:00:00Z",
					"updated_at": "2025-12-02T08:30:00Z",
				},
				{
					"id":        "b1d2c3a5",
					"title":     "Understanding RESTful API Design",
					"content":   "Designing RESTful APIs correctly improves maintainability.",
					"image_url": "https://cdn.example.com/blogs/rest-api.jpg",
					"author": gin.H{
						"id":         "u1002",
						"name":       "Alice Tran",
						"avatar_url": "https://cdn.example.com/avatars/alice.jpg",
					},
					"view_count": 876,
					"created_at": "2025-12-03T09:20:00Z",
					"updated_at": "2025-12-03T09:20:00Z",
				},
			},
			"metadata": gin.H{
				"limit":     10,
				"offset":    0,
				"sortBy":    "created_at",
				"sortOrder": "asc",
				"total":     125,
			},
		}

		c.JSON(http.StatusOK, response)
	}
}

// GetByID implements [BlogHandler].
// GetByID godoc
// @Summary      Retrieve a blog post by slug
// @Description  Returns a blog post identified by its slug, including author info, view count and timestamps.
// @Tags         blogs
// @Accept       json
// @Produce      json
// @Param        slug path string true "Blog Slug"
// @Success      200 {object} map[string]interface{} "OK"
// @Failure      400 {object} map[string]interface{} "Bad Request"
// @Failure      404 {object} map[string]interface{} "Not Found"
// @Failure      500 {object} map[string]interface{} "Internal Server Error"
// @Router       /api/v1/blogs/{slug} [get]
func (h *blogHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		response := gin.H{
			"process_id": "9d1c9a52-4a0b-4e4e-9c59-1fbbdb62c213",
			"path":       "/api/v1/blogs/" + slug,
			"status": gin.H{
				"code": 200,
				"type": "OK",
			},
			"request": gin.H{},
			"errors":  []any{},
			"data": gin.H{
				"id":        slug,
				"title":     "Introduction to Golang for Backend Development",
				"content":   "Golang is a statically typed programming language designed at Google. It is widely used for backend development because of its performance, simplicity, and built-in concurrency support with goroutines and channels.",
				"image_url": "https://cdn.example.com/blogs/golang-intro.jpg",
				"author": gin.H{
					"id":         "u1001",
					"name":       "John Nguyen",
					"avatar_url": "https://cdn.example.com/avatars/john.jpg",
				},
				"view_count": 1204,
				"created_at": "2025-12-01T10:00:00Z",
				"updated_at": "2025-12-02T08:30:00Z",
			},
			"metadata": nil,
		}

		c.JSON(http.StatusOK, response)
	}
}
