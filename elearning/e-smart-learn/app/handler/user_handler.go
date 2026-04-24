package handler

import (
	"fmt"
	"net/http"

	"elearning-api/apperror"
	"elearning-api/consts"
	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler interface {
	GetMe() gin.HandlerFunc
	UpdateMe() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	GetList() gin.HandlerFunc
	GetStatistics() gin.HandlerFunc
	GetTransactionsHistory() gin.HandlerFunc
	GetActiveStudentStatistics() gin.HandlerFunc

	ApplyToInstructor() gin.HandlerFunc
	GetPendingInstructorApplications() gin.HandlerFunc
	ApproveInstructorApplication() gin.HandlerFunc
	RejectInstructorApplication() gin.HandlerFunc
}

type userHandler struct {
	userService         service.UserService
	subscriptionService service.SubscriptionService
}

func NewUserHandler(
	userService service.UserService,
	subscriptionService service.SubscriptionService,
) UserHandler {
	return &userHandler{
		userService:         userService,
		subscriptionService: subscriptionService,
	}
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Router /api/v1/admin/users/{id} [get]
func (h *userHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}
		// Call service
		data, err := h.userService.GetByID(c.Request.Context(), userID)
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

// GetMe godoc
// @Summary Get current user details
// @Description Get details of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Failure 401 {object} any
// @Router /api/v1/users/me [get]
func (h *userHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		data, err := h.userService.GetByID(c.Request.Context(), userID)
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

// UpdateMe godoc
// @Summary Update current user details
// @Description Update details of the currently authenticated user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body dto.UpdateUserRequest true "Update user payload"
// @Success 200 {object} any
// @Failure 400 {object} any
// @Failure 401 {object} any
// @Router /api/v1/users/me [put]
func (h *userHandler) UpdateMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}
		var request dto.UpdateUserRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		data, err := h.userService.UpdateUser(c.Request.Context(), userID, request)
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

// ListUsers godoc
// @Summary List users
// @Description Get list of users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} any
// @Failure 400 {object} any
// @Router /api/v1/admin/users [get]
func (h *userHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var paginationRequest dto.PagingRequest

		// Bind query params
		if err := c.ShouldBindQuery(&paginationRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid pagination parameters"))
			return
		}

		// Process default values
		paginationRequest.Process()

		limit := paginationRequest.Limit
		offset := paginationRequest.Offset
		users, total, err := h.userService.GetList(
			c.Request.Context(),
			limit,
			offset,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = users
		res.Metadata = dto.NewPagination(limit, offset, int(total), "created_at", "desc")

		c.JSON(http.StatusOK, res)
	}
}

// GetStatistics godoc
// @Summary Get user management statistics
// @Description Get summary stats for user management dashboard (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/users/statistics [get]
func (h *userHandler) GetStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.userService.GetStatistics(c.Request.Context())
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

// GetActiveStudentStatistics godoc
// @Summary Get active student statistics
// @Description Get total active students in the system (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/users/statistics/active-students [get]
func (h *userHandler) GetActiveStudentStatistics() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.userService.GetActiveStudentStatistics(c.Request.Context())
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

// ApplyToInstructor
// @Summary Apply to become an instructor (approval needed)
// @Description Authenticated user applies to become an instructor
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.ApplyInstructorRequest true "Apply to instructor payload"
// @Success 200 {object} dto.ApiResponse "Application submitted"
// @Failure 400 {object} dto.ApiResponse "Bad request"
// @Failure 401 {object} dto.ApiResponse "Unauthorized"
// @Failure 403 {object} dto.ApiResponse "Forbidden"
// @Router /api/v1/users/me/apply-instructor [post]
func (h *userHandler) ApplyToInstructor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.ApplyInstructorRequest

		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}
		fmt.Printf("Received apply to instructor request: %+v\n", *request.CategoryID)

		err = h.userService.ApplyToInstructor(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Data = map[string]string{
			"message": "Application to become an instructor has been submitted",
		}

		c.JSON(http.StatusOK, res)
	}
}

// GetPendingInstructorApplications godoc
// @Summary Get pending instructor applications
// @Description Retrieve list of pending instructor applications (admin only)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Param sortBy query string false "sortBy" default(created_at)
// @Param sortOrder query string false "sortOrder" Enums(asc,desc) default(desc)
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/instructor-applications [get]
func (h *userHandler) GetPendingInstructorApplications() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.PagingRequest
		if err := c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid pagination parameters"))
			return
		}
		request.Process()

		data, err := h.userService.GetPendingInstructorApplications(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, len(data), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// ApproveInstructorApplication godoc
// @Summary Approve instructor application
// @Description Approve a pending instructor application (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/instructor-applications/{id}/approve [post]
func (h *userHandler) ApproveInstructorApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		_, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userResponse, err := h.userService.UpdateInstructorApplicationStatus(
			c.Request.Context(),
			idRequest.ID,
			consts.InstructorProfileApproved,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Data = map[string]any{
			"message": "Instructor application approved",
			"user":    userResponse,
		}

		c.JSON(http.StatusOK, res)
	}
}

// RejectInstructorApplication godoc
// @Summary Reject instructor application
// @Description Reject a pending instructor application (admin only)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/instructor-applications/{id}/reject [post]
func (h *userHandler) RejectInstructorApplication() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		_, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		userResponse, err := h.userService.UpdateInstructorApplicationStatus(
			c.Request.Context(),
			idRequest.ID,
			consts.InstructorProfileRejected,
		)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Data = map[string]any{
			"message": "Instructor application rejected",
			"user":    userResponse,
		}

		c.JSON(http.StatusOK, res)
	}
}

// GetTransactionsHistory godoc
// @Summary Get subscription transactions history
// @Description Get list of all subscription transactions for the current user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/users/me/transactions [get]
func (u *userHandler) GetTransactionsHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := u.subscriptionService.GetTransactionsHistory(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}
