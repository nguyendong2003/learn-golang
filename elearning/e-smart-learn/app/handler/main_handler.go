package handler

import (
	"elearning-api/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainHandler interface {
	HealthCheck() gin.HandlerFunc
}

type mainHandler struct{}

func NewMainHandler() MainHandler {
	return &mainHandler{}
}

// HealthCheck godoc
// @Summary Ping
// @Description Check service health
// @Tags main
// @Accept json
// @Produce json
// @Success 200 {object} dto.ApiResponse
// @Router /api/v1/health-check [get]
func (h *mainHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := dto.NewApiResponse(c)
		res.Request = dto.GetRequestClient(c)
		res.Data = "Health check successful"
		c.JSON(http.StatusOK, res)
	}
}
