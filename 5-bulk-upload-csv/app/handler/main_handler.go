package handler

import (
	"bulk-upload-csv/dto"
	"bulk-upload-csv/interfaces"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MainHandler struct{}

func NewMainHandler() interfaces.MainHandlerInterface {
	return &MainHandler{}
}

func (h *MainHandler) HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		response:= dto.NewApiResponse(c.FullPath())
		c.JSON(http.StatusOK, response)
	}
}