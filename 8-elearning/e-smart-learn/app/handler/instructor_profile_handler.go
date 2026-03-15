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

type InstructorProfileHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc

	GetByID() gin.HandlerFunc
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
			c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			c.Error(err)
			return
		}

		data, err := h.instructorProfileService.Create(c.Request.Context(), userID, request)
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
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Bind & validate JSON body
		var request dto.UpdateInstructorProfileRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			c.Error(err)
			return
		}

		// Call service
		data, err := h.instructorProfileService.Update(c.Request.Context(), id, request)

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
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		if err := h.instructorProfileService.Delete(c.Request.Context(), id); err != nil {
			c.Error(err)
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
			c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		id, err := uuid.Parse(idRequest.ID)
		if err != nil {
			c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		// Call service
		data, err := h.instructorProfileService.GetByID(c.Request.Context(), id)
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
