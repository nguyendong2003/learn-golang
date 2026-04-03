package handler

import (
	"io"
	"net/http"

	"elearning-api/dto"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler interface {
	CreateSubscriptionCheckoutSession() gin.HandlerFunc
	WebhookStripe() gin.HandlerFunc
	GetMySubscription() gin.HandlerFunc
	CancelAtPeriodEnd() gin.HandlerFunc
	Resume() gin.HandlerFunc
	CreateBillingPortalSession() gin.HandlerFunc
	GetSubscribers() gin.HandlerFunc

	// Get subscription retention statistics for admin dashboard
	GetMemberRetention() gin.HandlerFunc
}

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// CreateSubscriptionCheckoutSession godoc
// @Summary Create subscription checkout session
// @Description Create Stripe checkout session for subscription plan purchase
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param payload body dto.CreateSubscriptionCheckoutSessionRequest true "Checkout session payload"
// @Success 200 {object} dto.ApiResponse
// @Failure 400 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/checkout-session [post]
// @Security BearerAuth
func (h *subscriptionHandler) CreateSubscriptionCheckoutSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request dto.CreateSubscriptionCheckoutSessionRequest
		if err := util.BindAndValidateJSON(c, &request); err != nil {
			_ = c.Error(err)
			return
		}

		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.CreateSubscriptionCheckoutSession(c.Request.Context(), userID, request)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}

// WebhookStripe godoc
// @Summary Stripe webhook handler
// @Description Receive and process Stripe webhook events
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param Stripe-Signature header string true "Stripe webhook signature"
// @Success 200 {object} map[string]bool
// @Failure 400 {object} map[string]string
// @Failure 401 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/webhook/stripe [post]
func (h *subscriptionHandler) WebhookStripe() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			_ = c.Error(err)
			return
		}

		signature := c.GetHeader("Stripe-Signature")
		if signature == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "missing stripe signature"})
			return
		}

		if err := h.subscriptionService.HandleStripeWebhook(c.Request.Context(), payload, signature); err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"received": true})
	}
}

// GetMySubscription godoc
// @Summary Get my subscription
// @Description Get current user's latest subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/me [get]
// @Security BearerAuth
func (h *subscriptionHandler) GetMySubscription() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.GetMySubscription(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}

// CancelAtPeriodEnd godoc
// @Summary Cancel subscription at period end
// @Description Mark current user's active subscription to cancel at period end
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/cancel [post]
// @Security BearerAuth
func (h *subscriptionHandler) CancelAtPeriodEnd() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.CancelAtPeriodEnd(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}

// Resume godoc
// @Summary Resume subscription
// @Description Resume current user's subscription by clearing cancel_at_period_end
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/resume [post]
// @Security BearerAuth
func (h *subscriptionHandler) Resume() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.Resume(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}

// CreateBillingPortalSession godoc
// @Summary Create billing portal session
// @Description Create Stripe billing portal session for current user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 404 {object} dto.ApiResponse
// @Router /api/v1/subscriptions/billing-portal [post]
// @Security BearerAuth
func (h *subscriptionHandler) CreateBillingPortalSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := util.GetRequestUserID(c)
		if err != nil {
			_ = c.Error(err)
			return
		}

		result, err := h.subscriptionService.CreateBillingPortalSession(c.Request.Context(), userID)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = result

		c.JSON(http.StatusOK, response)
	}
}

// GetSubscribers godoc
// @Summary Get all subscribers
// @Description Get paginated list of all active subscribers
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param limit query int false "Limit per page" default(10)
// @Param offset query int false "Offset from start" default(0)
// @Success 200 {object} dto.ApiResponse
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Router /api/v1/admin/subscriptions/subscribers [get]
// @Security BearerAuth
func (h *subscriptionHandler) GetSubscribers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pagingRequest dto.PagingRequest
		if err := util.BindAndValidateQuery(c, &pagingRequest); err != nil {
			_ = c.Error(err)
			return
		}
		pagingRequest.Process()

		subscribers, total, err := h.subscriptionService.GetSubscribers(c.Request.Context(), pagingRequest.Limit, pagingRequest.Offset)
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Request = dto.GetRequestClient(c)
		response.Data = subscribers
		response.Metadata = dto.NewPagination(pagingRequest.Limit, pagingRequest.Offset, int(total), "created_at", "desc")

		c.JSON(http.StatusOK, response)
	}
}

// GetMemberRetention godoc
// @Summary Get member retention statistics
// @Description Get active memberships and retention percentage compared to the previous month for admin dashboard.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ApiResponse{data=dto.MemberRetentionResponse}
// @Failure 401 {object} dto.ApiResponse
// @Failure 403 {object} dto.ApiResponse
// @Failure 500 {object} dto.ApiResponse
// @Router /api/v1/admin/subscriptions/statistics/members/retention [get]
func (h *subscriptionHandler) GetMemberRetention() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.subscriptionService.GetMemberRetention(c.Request.Context())
		if err != nil {
			_ = c.Error(err)
			return
		}

		response := dto.NewApiResponse(c)
		response.Status = dto.NewResponseStatus(http.StatusOK)
		response.Request = dto.GetRequestClient(c)
		response.Data = stats

		c.JSON(http.StatusOK, response)
	}
}
