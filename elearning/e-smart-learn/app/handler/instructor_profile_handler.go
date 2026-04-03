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

type InstructorProfileHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
	GetList() gin.HandlerFunc
	GetGrowthStatistics() gin.HandlerFunc
}

type instructorProfileHandler struct {
	instructorProfileService service.InstructorProfileService
}

func NewInstructorProfileHandler(
	instructorProfileService service.InstructorProfileService,
) InstructorProfileHandler {
	return &instructorProfileHandler{
		instructorProfileService: instructorProfileService,
	}
}

// Create godoc
// @Summary Create a new instructorProfile
// @Description Create a new instructorProfile with the provided name and description. InstructorProfile name must be unique.
// @Tags instructor-profiles
// @Accept json
// @Produce json
// @Param payload body dto.CreateInstructorProfileRequest true "InstructorProfile create request"
// @Success 201 {object} dto.ApiResponse{data=dto.InstructorProfileResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid request payload or validation failed"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/instructor-profiles [post]
// @Security BearerAuth
func (h *instructorProfileHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateInstructorProfileRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		data, err := h.instructorProfileService.Create(c.Request.Context(), userID, request)
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

// Update godoc
// @Summary Update an existing instructorProfile
// @Description Update a instructorProfile by ID with the provided name and description. InstructorProfile name must be unique.
// @Tags instructor-profiles
// @Accept json
// @Produce json
// @Param id path string true "InstructorProfile ID (UUID format)"
// @Param payload body dto.UpdateInstructorProfileRequest true "InstructorProfile update request"
// @Success 200 {object} dto.ApiResponse{data=dto.InstructorProfileResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format or validation failed"
// @Failure 404 {object} dto.ApiResponse "InstructorProfile not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/instructor-profiles/{id} [put]
// @Security BearerAuth
func (h *instructorProfileHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Bind & validate JSON body
		var request dto.UpdateInstructorProfileRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		// Call service
		data, err := h.instructorProfileService.Update(c.Request.Context(), id, request)
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

// Delete godoc
// @Summary Delete a instructorProfile
// @Description Delete a instructorProfile by ID. Cannot delete a instructorProfile that has associated courses.
// @Tags instructor-profiles
// @Accept json
// @Produce json
// @Param id path string true "InstructorProfile ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "InstructorProfile not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/instructor-profiles/{id} [delete]
// @Security BearerAuth
func (h *instructorProfileHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		if err := h.instructorProfileService.Delete(c.Request.Context(), id); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get instructor profile by ID
// @Description Retrieve a single instructorProfile by its ID
// @Tags instructor-profiles
// @Accept json
// @Produce json
// @Param id path string true "InstructorProfile ID (UUID format)"
// @Success 200 {object} dto.ApiResponse{data=dto.InstructorProfileResponse}
// @Failure 400 {object} dto.ApiResponse "Invalid UUID format"
// @Failure 404 {object} dto.ApiResponse "InstructorProfile not found"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/instructor-profiles/{id} [get]
func (h *instructorProfileHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		data, err := h.instructorProfileService.GetByID(c.Request.Context(), id)
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

// GetList godoc
// @Summary Get paginated list of instructor profiles
// @Description Retrieve a paginated list of instructor profiles with filtering, sorting, and pagination support
// @Tags instructor-profiles
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page (default: 10, max: 100)" default(10)
// @Param offset query int false "Number of items to skip (default: 0)" default(0)
// @Param sortBy query string false "Field to sort by (default: created_at)" default(created_at)
// @Param sortOrder query string false "Sort order: asc or desc (default: desc)" default(desc) Enums(asc,desc)
// @Param name query string false "Filter instructor profiles by name (partial match)"
// @Success 200 {object} dto.ApiResponse{data=[]dto.InstructorProfileResponse,metadata=dto.Pagination}
// @Failure 400 {object} dto.ApiResponse "Invalid query parameters"
// @Failure 500 {object} dto.ApiResponse "Internal server error"
// @Router /api/v1/instructor-profiles [get]
func (h *instructorProfileHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListInstructorProfileRequest

		if err := c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		// Process default pagination
		request.Process()

		limit := request.Limit
		offset := request.Offset
		sortBy := request.SortBy
		sortOrder := request.SortOrder

		// Call service
		data, total, err := h.instructorProfileService.GetList(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(limit, offset, int(total), sortBy, sortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// GetGrowthStatistics godoc
// @Summary Get teacher growth statistics
// @Description Get current quarter teacher growth statistics, including verified teachers and top category share.
// @Tags statistics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.TeacherGrowthStatisticsResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/instructor-profiles/statistics/teachers/growth [get]
func (h *instructorProfileHandler) GetGrowthStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.instructorProfileService.GetGrowthStatistics(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = stats

		c.JSON(http.StatusOK, res)
	}
}
