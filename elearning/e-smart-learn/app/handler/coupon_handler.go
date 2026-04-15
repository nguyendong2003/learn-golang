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
	GetList() gin.HandlerFunc
	GetAssignableList() gin.HandlerFunc
	GetAssignableListForCourse() gin.HandlerFunc
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
// @Router /api/v1/instructor/coupons [post]
// @Security BearerAuth
func (h *couponHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.CreateCouponRequest
		if err = util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.couponService.Create(c.Request.Context(), userID, request)
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
// @Router /api/v1/instructor/coupons/{id}/deactivate [put]
// @Security BearerAuth
func (h *couponHandler) Deactivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		id := c.Param("id")
		couponID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid coupon UUID format"))
			return
		}

		if err := h.couponService.Deactivate(c.Request.Context(), userID, couponID); err != nil {
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
// @Router /api/v1/instructor/coupons/{id} [get]
func (h *couponHandler) GetByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		id := c.Param("id")
		couponID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid coupon UUID format"))
			return
		}

		result, err := h.couponService.GetByID(c.Request.Context(), userID, couponID)
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
// @Router /api/v1/instructor/coupons [get]
func (h *couponHandler) GetList() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.ListCouponRequest

		if err = c.ShouldBindQuery(&request); err != nil {
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
		data, total, err := h.couponService.GetList(c.Request.Context(), userID, request)
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

// GetAssignableList godoc
// @Summary List assignable coupons
// @Description Retrieve paginated list of assignable coupons for course creation
// @Tags coupons
// @Accept json
// @Produce json
// @Param limit query int false "Items per page" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Param code query string false "Search by coupon code"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/instructor/coupons/assignable [get]
// @Security BearerAuth
func (h *couponHandler) GetAssignableList() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.ListCouponRequest
		if err = c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		request.Process()

		data, total, err := h.couponService.GetAssignableList(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}

// GetAssignableListForCourse godoc
// @Summary List assignable coupons for course
// @Description Retrieve paginated list of assignable coupons with pre-selected coupons prioritized
// @Tags coupons
// @Accept json
// @Produce json
// @Param course_id query string true "Course ID (UUID format)"
// @Param limit query int false "Items per page" default(10)
// @Param offset query int false "Offset" default(0)
// @Param sortBy query string false "Sort field" default(created_at)
// @Param sortOrder query string false "Sort order" Enums(asc,desc) default(desc)
// @Param code query string false "Search by coupon code"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/instructor/coupons/assignable-for-course [get]
// @Security BearerAuth
func (h *couponHandler) GetAssignableListForCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		var request dto.ListCouponRequest
		if err = c.ShouldBindQuery(&request); err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid query parameters"))
			return
		}

		if request.CourseID == nil || *request.CourseID == "" {
			_ = c.Error(apperror.NewBadRequestError("course_id is required"))
			return
		}

		courseID, err := uuid.Parse(*request.CourseID)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course ID format"))
			return
		}

		request.Process()

		data, total, err := h.couponService.GetAssignableListForCourse(c.Request.Context(), userID, courseID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = data
		res.Metadata = dto.NewPagination(request.Limit, request.Offset, int(total), request.SortBy, request.SortOrder)

		c.JSON(http.StatusOK, res)
	}
}
