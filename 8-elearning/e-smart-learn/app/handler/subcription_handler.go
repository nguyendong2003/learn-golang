package handler

import (
	"elearning-api/dto"
	"elearning-api/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SubscriptionHandler defines handler methods for subscription-related endpoints.
// In this project we use it to serve subscription plan mocks.
type SubscriptionHandler interface {
	GetPlans() gin.HandlerFunc
}

type subscriptionHandler struct{}

// NewSubscriptionHandler creates a new instance of SubscriptionHandler.
func NewSubscriptionHandler() SubscriptionHandler {
	return &subscriptionHandler{}
}

// GetPlans godoc
// @Summary List subscription plans
// @Description Return list of subscription plans (mocked)
// @Tags subscriptions
// @Accept json
// @Produce json
// @Success 200 {object} any
// @Router /api/v1/subscriptions [get]
//
// GetPlans returns a Gin handler that responds with a
// hardcoded list of subscription plans. This is a mock implementation
// intended for testing and documentation; it does not access any database.
func (h *subscriptionHandler) GetPlans() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create mock plans
		p1 := &model.Plan{
			Name:         "Basic",
			Description:  "Basic access to core features",
			MonthlyPrice: 0,
			YearlyPrice:  0,
			IsDefault:    true,
		}
		p2 := &model.Plan{
			Name:         "Pro",
			Description:  "Pro plan with additional features",
			MonthlyPrice: 9,
			YearlyPrice:  90,
			IsDefault:    false,
		}

		list := dto.NewListPlanResponse([]*model.Plan{p1, p2})

		resp := dto.NewApiResponse(c)
		resp.Path = "/api/v1/subscriptions"
		resp.Request = gin.H{}
		resp.Data = list
		resp.Metadata = nil

		c.JSON(http.StatusOK, resp)
	}
}
