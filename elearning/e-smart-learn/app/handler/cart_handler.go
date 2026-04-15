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

type CartHandler interface {
	AddCourse() gin.HandlerFunc
	RemoveCourse() gin.HandlerFunc
	GetMyCart() gin.HandlerFunc
	PreviewCheckout() gin.HandlerFunc
	Checkout() gin.HandlerFunc
}

type cartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) CartHandler {
	return &cartHandler{cartService: cartService}
}

// AddCourse godoc
// @Summary Add course to cart
// @Description Add a published paid course to current user's cart
// @Tags carts
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/carts/courses/{course_id} [post]
// @Security BearerAuth
func (h *cartHandler) AddCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("course_id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if err := h.cartService.AddCourse(c.Request.Context(), userID, courseID); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Course added to cart"}
		c.JSON(http.StatusOK, res)
	}
}

// RemoveCourse godoc
// @Summary Remove course from cart
// @Description Remove a course from current user's cart
// @Tags carts
// @Accept json
// @Produce json
// @Param course_id path string true "Course ID (UUID format)"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/carts/courses/{course_id} [delete]
// @Security BearerAuth
func (h *cartHandler) RemoveCourse() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("course_id")
		courseID, err := uuid.Parse(id)
		if err != nil {
			_ = c.Error(apperror.NewBadRequestError("Invalid course UUID format"))
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		if err := h.cartService.RemoveCourse(c.Request.Context(), userID, courseID); err != nil {
			_ = c.Error(err)
			return
		}

		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = gin.H{"message": "Course removed from cart"}
		c.JSON(http.StatusOK, res)
	}
}

// GetMyCart godoc
// @Summary Get my cart
// @Description Retrieve current user's cart with course items
// @Tags carts
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/carts [get]
// @Security BearerAuth
func (h *cartHandler) GetMyCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.cartService.GetMyCart(c.Request.Context(), userID)
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

// PreviewCheckout godoc
// @Summary Preview cart checkout
// @Description Preview cart pricing with entered coupon and default coupon fallback per course
// @Tags carts
// @Accept json
// @Produce json
// @Param payload body dto.CartCheckoutRequest false "Cart checkout preview request"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/carts/checkout-preview [post]
// @Security BearerAuth
func (h *cartHandler) PreviewCheckout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CartCheckoutRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.cartService.PreviewCheckout(c.Request.Context(), userID, request)
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

// Checkout godoc
// @Summary Create checkout session from cart
// @Description Create Stripe checkout session (payment mode) for all items in current user's cart with optional coupon
// @Tags carts
// @Accept json
// @Produce json
// @Param payload body dto.CartCheckoutRequest true "Cart checkout request"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/carts/checkout-session [post]
// @Security BearerAuth
func (h *cartHandler) Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CartCheckoutRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.cartService.Checkout(c.Request.Context(), userID, request)
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
