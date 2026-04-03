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

type BlogHandler interface {
	GetBySlug() gin.HandlerFunc
	GetPublishedBlogs() gin.HandlerFunc
	GetBlogs() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	GetStatistics() gin.HandlerFunc
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

// GetPublishedBlogs godoc
// @Summary Get list of published blogs
// @Description Retrieve paginated list of published blogs with optional filters
// @Tags blogs
// @Accept json
// @Produce json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Param sortBy query string false "sortBy" default(created_at) Enums(created_at,title,view_total)
// @Param sortOrder query string false "sortOrder" Enums(asc,desc) default(desc)
// @Param category_id query string false "category ID"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/blogs [get]
func (h *blogHandler) GetPublishedBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter dto.SearchBlogRequest
		if err := util.BindAndValidateQuery(c, &filter); err != nil {
			_ = c.Error(err)
			return
		}

		blogs, total, err := h.blogService.GetPublishedBlogs(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}
		resp := dto.NewApiResponse(c)
		resp.Data = blogs
		resp.Request = dto.GetRequestClient(c)
		resp.Metadata = dto.NewPagination(filter.Limit, filter.Offset, int(total), filter.SortBy, filter.SortOrder)

		c.JSON(http.StatusOK, resp)
	}
}

// GetBySlug implements [BlogHandler].
// GetBySlug godoc
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
func (h *blogHandler) GetBySlug() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Get chapters > lessons of course if user is enrolled in course
		slug := c.Param("slug")
		if slug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "slug is required"})
			return
		}
		userId, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		blog, err := h.blogService.GetBlogBySlug(c.Request.Context(), slug, userId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp := dto.NewApiResponse(c)
		resp.Data = blog
		resp.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, resp)
	}
}

// GetBlogs godoc
// @Summary Get list of blogs with filters (admin)
// @Description Retrieve paginated list of blogs with optional filters (admin only)
// @Tags blogs
// @Accept json
// @Produce json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Param sortBy query string false "sortBy" default(created_at) Enums(created_at,title,view_total)
// @Param sortOrder query string false "sortOrder" Enums(asc,desc) default(desc)
// @Param category_id query string false "category ID"
// @Param status query string false "blog status filter" Enums(draft,published,scheduled)
// @Param keyword query string false "search keyword (title)"
// @Security BearerAuth
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Failure 500 {object} any
// @Router /api/v1/admin/blogs [get]
func (h *blogHandler) GetBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter dto.SearchBlogRequest
		if err := util.BindAndValidateQuery(c, &filter); err != nil {
			_ = c.Error(err)
			return
		}

		blogs, total, err := h.blogService.GetBlogs(c.Request.Context(), filter)
		if err != nil {
			_ = c.Error(err)
			return
		}
		resp := dto.NewApiResponse(c)
		resp.Data = blogs
		resp.Request = dto.GetRequestClient(c)
		resp.Metadata = dto.NewPagination(filter.Limit, filter.Offset, int(total), filter.SortBy, filter.SortOrder)

		c.JSON(http.StatusOK, resp)
	}
}

// GetByID godoc
// @Summary Get a blog post by ID (admin)
// @Description Retrieve a single blog post by its ID (admin only)
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Security BearerAuth
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Failure 404 {object} any
// @Failure 500 {object} any
// @Router /api/v1/admin/blogs/{id} [get]
func (h *blogHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		blogId, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid blog ID"))
			return
		}

		blog, err := h.blogService.GetByID(c.Request.Context(), blogId)
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp := dto.NewApiResponse(c)
		resp.Data = blog
		resp.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, resp)
	}
}

// GetStatistics godoc
// @Summary Get blog dashboard statistics (admin)
// @Description Get blog KPI statistics for current user (total articles, total views, average engagement, counts by status)
// @Tags blogs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} any
// @Failure 401 {object} any
// @Failure 500 {object} any
// @Router /api/v1/admin/blogs/statistics [get]
func (h *blogHandler) GetStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.blogService.GetBlogStatistics(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		resp := dto.NewApiResponse(c)
		resp.Request = dto.GetRequestClient(c)
		resp.Data = stats

		c.JSON(http.StatusOK, resp)
	}
}

// CreateBlog godoc
// @Summary Create a blog post
// @Description Create a new blog post and return created resource
// @Tags blogs
// @Accept json
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
			_ = c.Error(err)
			return
		}

		var request dto.CreateBlogRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error((err))
			return
		}

		blog, err := h.blogService.Create(c.Request.Context(), creatorID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		response.Status = dto.NewResponseStatus(http.StatusCreated)
		response.Request = dto.GetRequestClient(c)
		response.Data = blog

		c.JSON(http.StatusCreated, response)
	}
}

// UpdateBlog godoc
// @Summary Update a blog post
// @Description Update a blog post identified by id
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Param payload body dto.UpdateBlogRequest true "Update blog payload"
// @Security BearerAuth
// @Success 200 {object} any
// @Router /api/v1/blogs/{id} [put]
func (h *blogHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		blogId, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid blog ID"))
			return
		}
		authorId, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.UpdateBlogRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error((err))
			return
		}
		updatedBlog, err := h.blogService.UpdateBlogs(c.Request.Context(), authorId, blogId, request)
		if err != nil {
			_ = c.Error(err)
			return
		}
		resp := dto.NewApiResponse(c)
		resp.Request = dto.GetRequestClient(c)
		resp.Data = updatedBlog

		c.JSON(http.StatusOK, resp)
	}
}

// DeleteBlog godoc
// @Summary Delete a blog post
// @Description Delete a blog post identified by id
// @Tags blogs
// @Accept json
// @Produce json
// @Param id path string true "Blog ID"
// @Security BearerAuth
// @Success 200 {object} any
// @Router /api/v1/blogs/{id} [delete]
func (h *blogHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		blogId, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid blog ID"))
			return
		}
		authorId, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		if err := h.blogService.DeleteBlog(c.Request.Context(), authorId, blogId); err != nil {
			_ = c.Error(err)
			return
		}
		resp := dto.NewApiResponse(c)
		resp.Request = dto.GetRequestClient(c)
		resp.Data = gin.H{"deleted": true, "id": id}

		c.JSON(http.StatusOK, resp)
	}
}
