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

type PlanHandler interface {
	Create() gin.HandlerFunc
	Update() gin.HandlerFunc
	Activate() gin.HandlerFunc
	Deactivate() gin.HandlerFunc
	Delete() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	GetActivePlans() gin.HandlerFunc
	GetList() gin.HandlerFunc
}

type planHandler struct {
	planService service.PlanService
}

func NewPlanHandler(planService service.PlanService) PlanHandler {
	return &planHandler{
		planService: planService,
	}
}

// Create godoc
// @Summary Create subscription plan
// @Description Create a subscription plan and sync to Stripe and Postgres
// @Tags plans
// @Accept json
// @Produce json
// @Param payload body dto.CreatePlanRequest true "Create plan payload"
// @Success 201 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans [post]
// @Security BearerAuth
func (h *planHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreatePlanRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.planService.Create(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Status = dto.NewResponseStatus(http.StatusCreated)
		res.Request = dto.GetRequestClient(c)
		res.Data = result

		c.JSON(http.StatusCreated, res)
	}
}

// Update godoc
// @Summary Update subscription plan
// @Description Update subscription plan and sync changed fields to Stripe and Postgres
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID (UUID format)"
// @Param payload body dto.UpdatePlanRequest true "Update plan payload"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans/{id} [put]
// @Security BearerAuth
func (h *planHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		planID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		var request dto.UpdatePlanRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		updated, err := h.planService.Update(c.Request.Context(), planID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updated

		c.JSON(http.StatusOK, res)
	}
}

// Activate godoc
// @Summary Activate subscription plan
// @Description Activate plan for new sales
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans/{id}/activate [put]
// @Security BearerAuth
func (h *planHandler) Activate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		planID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		updated, err := h.planService.Activate(c.Request.Context(), planID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updated

		c.JSON(http.StatusOK, res)
	}
}

// Deactivate godoc
// @Summary Deactivate subscription plan
// @Description Deactivate plan for new sales and set related active subscriptions to cancel at period end
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans/{id}/deactivate [put]
// @Security BearerAuth
func (h *planHandler) Deactivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		planID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		updated, err := h.planService.Deactivate(c.Request.Context(), planID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = updated

		c.JSON(http.StatusOK, res)
	}
}

// Delete godoc
// @Summary Delete subscription plan
// @Description Delete subscription plan in Postgres and deactivate linked Stripe price/product
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans/{id} [delete]
// @Security BearerAuth
func (h *planHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		planID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		if err := h.planService.Delete(c.Request.Context(), planID); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get subscription plan by id
// @Description Get a subscription plan detail by id
// @Tags plans
// @Accept json
// @Produce json
// @Param id path string true "Plan ID"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/subscription-plans/{id} [get]
func (h *planHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		var idRequest dto.UUIDRequest

		if err := c.ShouldBindUri(&idRequest); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID in URI"))
			return
		}

		planID, err := uuid.Parse(idRequest.ID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid UUID format"))
			return
		}

		plan, err := h.planService.GetByID(c.Request.Context(), planID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = plan

		c.JSON(http.StatusOK, res)
	}
}

// GetActivePlans godoc
// @Summary List subscription plans
// @Description Return list of subscription plans
// @Tags plans
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/v1/subscription-plans/active [get]
func (h *planHandler) GetActivePlans() gin.HandlerFunc {
	return func(c *gin.Context) {
		plans, err := h.planService.GetActivePlans(c.Request.Context())
		if err != nil {
			_ = c.Error((err))
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = plans

		c.JSON(http.StatusOK, res)
	}
}

// GetList godoc
// @Summary List all subscription plans
// @Description Return all plans for admin, including inactive plans
// @Tags plans
// @Accept json
// @Produce json
// @Param limit query int false "Page size" default(20)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/subscription-plans [get]
// @Security BearerAuth
func (h *planHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListPlanRequest

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
		data, total, err := h.planService.GetList(c.Request.Context(), request)
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
