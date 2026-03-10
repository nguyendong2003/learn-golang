package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubscriptionHandler defines handler methods for subscription-related endpoints.
// In this project we use it to serve subscription plan mocks.
type SubscriptionHandler interface {
	GetSupcriptions() gin.HandlerFunc
}

type subscriptionHandler struct{}

// NewSubscriptionHandler creates a new instance of SubscriptionHandler.
func NewSubscriptionHandler() SubscriptionHandler {
	return &subscriptionHandler{}
}

// GetSupcriptions godoc
// @Summary List subscription plans
// @Description Return list of subscription plans (mocked)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/v1/subscriptions [get]
//
// GetSupcriptions returns a Gin handler that responds with a
// hardcoded list of subscription plans. This is a mock implementation
// intended for testing and documentation; it does not access any database.
func (h *subscriptionHandler) GetSupcriptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := gin.H{
			"process_id": "uuid",
			"path":       "/api/v1/subscriptions",
			"status": gin.H{
				"code": 200,
				"type": "OK",
			},
			"request": gin.H{},
			"errors":  []any{},
			"data": []gin.H{
				{
					"id":              "basic",
					"name":            "Basic",
					"description":     "Basic access to core features",
					"access_features": []string{"Feature A", "Feature B"},
					"price_monthly":   0,
					"price_yearly":    0,
					"is_default":      true,
				},
				{
					"id":              "pro",
					"name":            "Pro",
					"description":     "Pro plan with additional features",
					"access_features": []string{"Feature A", "Feature B", "Feature C"},
					"price_monthly":   9,
					"price_yearly":    90,
					"is_default":      false,
				},
			},
			"metadata": nil,
		}

		c.JSON(http.StatusOK, response)
	}
}
