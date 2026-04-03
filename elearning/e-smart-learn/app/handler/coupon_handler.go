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

type CouponHandler interface {
	Create() gin.HandlerFunc
	GetByID() gin.HandlerFunc
	Deactivate() gin.HandlerFunc
	GetAvailableCoupons() gin.HandlerFunc
	GetList() gin.HandlerFunc
}

type couponHandler struct {
	couponService service.CouponService
}

func NewCouponHandler(couponService service.CouponService) CouponHandler {
	return &couponHandler{couponService: couponService}
}

// Create godoc
// @Summary Create coupon
// @Description Create a coupon and sync to Stripe promotion code
// @Tags coupons
// @Accept json
// @Produce json
// @Param payload body dto.CreateCouponRequest true "Create coupon payload"
// @Success 201 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/coupons [post]
// @Security BearerAuth
func (h *couponHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateCouponRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.couponService.Create(c.Request.Context(), request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = result

		c.JSON(http.StatusCreated, res)
	}
}

// Deactivate godoc
// @Summary Deactivate coupon
// @Description Deactivate coupon on Stripe and local database
// @Tags coupons
// @Accept json
// @Produce json
// @Param id path string true "Coupon ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/coupons/{id}/deactivate [put]
// @Security BearerAuth
func (h *couponHandler) Deactivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		couponID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid coupon UUID format"))
			return
		}

		if err := h.couponService.Deactivate(c.Request.Context(), couponID); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Coupon deactivated"}

		c.JSON(http.StatusOK, res)
	}
}

// GetByID godoc
// @Summary Get coupon by ID
// @Description Retrieve coupon detail by UUID
// @Tags coupons
// @Accept json
// @Produce json
// @Param id path string true "Coupon ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/admin/coupons/{id} [get]
func (h *couponHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		couponID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid coupon UUID format"))
			return
		}

		result, err := h.couponService.GetByID(c.Request.Context(), couponID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = result

		c.JSON(http.StatusOK, res)
	}
}

// GetList godoc
// @Summary List coupons
// @Description Retrieve paginated list of coupons
// @Tags coupons
// @Accept json
// @Produce json
// @Param limit query int false "Items per page" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Param code query string false "Filter by coupon code"
// @Param is_active query boolean false "Filter by active state"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Router /api/v1/admin/coupons [get]
func (h *couponHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListCouponRequest

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
		data, total, err := h.couponService.GetList(c.Request.Context(), request)
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

// GetAvailableCoupons godoc
// @Summary Get available coupons for user
// @Description Retrieve list of active coupons that the user can apply
// @Tags coupons
// @Accept json
// @Produce json
// @Param limit query int false "Items per page" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Param code query string false "Filter by coupon code"
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Security BearerAuth
// @Router /api/v1/coupons [get]
func (h *couponHandler) GetAvailableCoupons() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.ListCouponRequest

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
		data, total, err := h.couponService.GetAvailableCoupons(c.Request.Context(), request)
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
