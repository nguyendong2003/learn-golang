package handler

import (
	"elearning-api/dto"
	"elearning-api/model"
	"elearning-api/service"
	"elearning-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogHandler interface {
	GetByID() gin.HandlerFunc
	GetList() gin.HandlerFunc
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
}

type blogHandler struct {
	blogService service.BlogService
}

func NewBlogHandler(blogService service.BlogService) BlogHandler {
	return &blogHandler{
		blogService: blogService,
	}
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
		// create mock blog models
		b1 := &model.Blog{
			Title:     "Introduction to Golang for Backend Development",
			Content:   "Golang is a powerful language designed for building scalable backend systems.",
			Slug:      "b1d2c3a4",
			ViewTotal: 1204,
		}
		b1.Author = &model.InstructorProfile{
			User: &model.User{
				Name:   "John Nguyen",
				Avatar: "https://cdn.example.com/avatars/john.jpg",
			},
		}

		b2 := &model.Blog{
			Title:     "Understanding RESTful API Design",
			Content:   "Designing RESTful APIs correctly improves maintainability.",
			Slug:      "b1d2c3a5",
			ViewTotal: 876,
		}
		b2.Author = &model.InstructorProfile{
			User: &model.User{
				Name:   "Alice Tran",
				Avatar: "https://cdn.example.com/avatars/alice.jpg",
			},
		}

		list := dto.NewListBlogResponse([]*model.Blog{b1, b2})

		resp := dto.NewApiResponse(c)
		resp.Request = gin.H{"limit": 10, "offset": 0, "sortBy": "created_at", "sortOrder": "asc", "categoryId": "", "type": ""}
		resp.Data = list
		resp.Metadata = gin.H{"limit": 10, "offset": 0, "sortBy": "created_at", "sortOrder": "asc", "total": 125}

		c.JSON(http.StatusOK, resp)
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

		m := &model.Blog{
			Title:     "Introduction to Golang for Backend Development",
			Content:   "Golang is a statically typed programming language designed at Google. It is widely used for backend development because of its performance, simplicity, and built-in concurrency support with goroutines and channels.",
			Slug:      slug,
			ViewTotal: 1204,
		}
		m.Author = &model.InstructorProfile{
			User: &model.User{
				Name:   "John Nguyen",
				Avatar: "https://cdn.example.com/avatars/john.jpg",
			},
		}

		blog := dto.NewBlogDetailResponse(m)

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/blogs/" + slug
		resp.Request = gin.H{}
		resp.Data = blog
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}

// CreateBlog godoc
// @Summary Create a blog post
// @Description Create a new blog post and return created resource
// @Tags blogs
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param payload body dto.CreateBlogRequest true "Create blog payload"
// @Success 200 {object} dto.BlogResponse
// @Router /api/v1/blogs [post]
func (h *blogHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.NewApiResponse(c)
		var creatorID uuid.UUID
		creatorID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			return
		}

		var request dto.CreateBlogRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error((err))
			return
		}

		blog, err := h.blogService.Create(c.Request.Context(), creatorID, request)
		if err != nil {
			c.Error(err)
			return
		}
		response.Status = dto.NewResponseStatus(http.StatusCreated)
		response.Request = dto.GetRequestClient(c)
		response.Data = blog

		c.JSON(http.StatusCreated, response)
	}
}

// UpdateBlog godoc
// @Summary Update a blog post (mock)
// @Description Update an existing blog post identified by slug (mocked)
// @Tags blogs
// @Accept json
// @Produce json
// @Param slug path string true "Blog Slug"
// @Param payload body object true "Update blog payload"
// @Success 200 {object} any
// @Router /api/v1/blogs/{slug} [put]
func (h *blogHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		var payload struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		_ = c.BindJSON(&payload)

		m := &model.Blog{
			Title:     payload.Title,
			Content:   payload.Content,
			Slug:      slug,
			ViewTotal: 100,
		}
		m.Author = &model.InstructorProfile{
			User: &model.User{
				Name:   "John Nguyen",
				Avatar: "https://cdn.example.com/avatars/john.jpg",
			},
		}

		blog := dto.NewBlogDetailResponse(m)

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/blogs/" + slug
		resp.Request = dto.GetRequestClient(c)
		resp.Data = blog
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}

// DeleteBlog godoc
// @Summary Delete a blog post (mock)
// @Description Delete a blog post identified by slug (mocked)
// @Tags blogs
// @Accept json
// @Produce json
// @Param slug path string true "Blog Slug"
// @Success 200 {object} any
// @Router /api/v1/blogs/{slug} [delete]
func (h *blogHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("slug")

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/blogs/" + slug
		resp.Request = gin.H{"slug": slug}
		resp.Data = gin.H{"deleted": true, "id": slug}
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}
