package handler

import (
	"elearning-api/dto"
	"elearning-api/util"
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

		response := dto.NewSuccessResponse[any](
			nil,
			"Server healthcheck successfully",
			util.GetRequestID(c),
		)

		c.JSON(http.StatusOK, response)
	}
}
