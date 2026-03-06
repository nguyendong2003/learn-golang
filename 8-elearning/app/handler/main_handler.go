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

func (h *mainHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := dto.SimpleSuccess("Server healthcheck successfully")

		c.JSON(http.StatusOK, response)
	}
}
