package interfaces

import "github.com/gin-gonic/gin"

type GenerateHandlerInterface interface {
	Generate() gin.HandlerFunc
}
