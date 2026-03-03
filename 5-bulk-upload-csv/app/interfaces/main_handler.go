package interfaces

import "github.com/gin-gonic/gin"

type MainHandlerInterface interface {
	HealthCheck() gin.HandlerFunc
}
