package handler

import (
	"elearning-api/apperror"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
	GetBySlug() gin.HandlerFunc
	GetList() gin.HandlerFunc
}

type courseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(
	courseService service.CourseService,
	instructorProfileService service.InstructorProfileService,
) CourseHandler {
	return &courseHandler{
		courseService: courseService,
	}
}

// Create godoc
// @Summary Create a new course
// @Description Create a new course with the provided title, description, and category. Requires an instructor account.
// @Tags courses
// @Accept json
// @Produce json
// @Param payload body dto.CreateCourseRequest true "Course create request"
// @Success 201 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request payload or validation failed"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses [post]
// @Security BearerAuth
func (h *courseHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateCourseRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			return
		}

		createdCourse, err := h.courseService.Create(c.Request.Context(), userID, request)
		if err != nil {
			c.Error(err)
			return
		}

		data, err := h.courseService.GetByID(c.Request.Context(), uuid.MustParse(createdInstructorProfile.ID))
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusCreated, res)
	}
}

// Update godoc
// @Summary Update an existing course
// @Description Update a course by ID. Requires instructor permissions.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Param payload body dto.UpdateCourseRequest true "Course update request"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format or validation failed"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [put]
// @Security BearerAuth
func (h *courseHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Bind and validate request body
		var request dto.UpdateCourseRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			return
		}

		// Call service
		updatedCourse, err := h.courseService.Update(c.Request.Context(), userID, id, request)
		if err != nil {
			c.Error(err)
			return
		}

		data, err := h.courseService.GetByID(c.Request.Context(), uuid.MustParse(updatedCourse.ID))
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// Delete godoc
// @Summary Delete a course
// @Description Delete a course by ID.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [delete]
// @Security BearerAuth
func (h *courseHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Check if user is instructor and this course belongs to the instructor
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			return
		}

		// Call service
		if err := h.courseService.Delete(c.Request.Context(), userID, id); err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get course by ID
// @Description Retrieve a single course by its ID.
// @Tags courses
// @Accept json
// @Produce json
// @Param id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/{id} [get]
func (h *courseHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		data, err := h.courseService.GetByID(c.Request.Context(), id)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetBySlug godoc
// @Summary Get course by slug
// @Description Retrieve a single course by its slug.
// @Tags courses
// @Accept json
// @Produce json
// @Param slug path string true "Course slug"
// @Success 200 {object} dto.ApiResponse{data=dto.CourseResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid slug format"
// @Failure 404 {object} dto.ApiResponse "Course not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses/slug/{slug} [get]
func (h *courseHandler) GetBySlug() gin.HandlerFunc {
	return func(c *gin.Context) {
		var slugRequest dto.CourseSlugRequest

		if err := c.ShouldBindUri(&slugRequest); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid slug in URI"))
			return
		}

		// Call service
		data, err := h.courseService.GetBySlug(c.Request.Context(), slugRequest.Slug)
		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data

		c.JSON(http.StatusOK, res)
	}
}

// GetList godoc
// @Summary Get paginated list of courses
// @Description Retrieve a paginated list of courses with filtering, sorting, and pagination support.
// @Tags courses
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10, max: 100)" default(10)
// @Param offset query int false "Number of items to skip (default: 0)" default(0)
// @Param sortBy query string false "Field to sort by (default: created_at)" default(created_at)
// @Param sortOrder query string false "Sort order: asc or desc (default: desc)" default(desc) Enums(asc,desc)
// @Param title query string false "Filter courses by title (partial match)"
// @Param category_id query string false "Filter courses by category id (UUID format)"
// @Param status query string false "Filter courses by status" Enums(DRAFT,PUBLISHED,ARCHIVED)
// @Success 200 {object} dto.ApiResponse{data=[]dto.CourseResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse "Invalid query parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/courses [get]
func (h *courseHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListCourseRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		// Process default pagination
		request.Process()

		limit := request.Limit
		offset := request.Offset
		sortBy := request.SortBy
		sortOrder := request.SortOrder

		// Call service
		data, total, err := h.courseService.GetList(c.Request.Context(), request)

		if err != nil {
			c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}
