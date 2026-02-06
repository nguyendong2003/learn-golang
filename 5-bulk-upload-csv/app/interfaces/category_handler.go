package interfaces

import "github.com/gin-gonic/gin"

type CategoryHandlerInterface interface {
	GetList() gin.HandlerFunc
	GetDetail() gin.HandlerFunc
	Create() gin.HandlerFunc
}
