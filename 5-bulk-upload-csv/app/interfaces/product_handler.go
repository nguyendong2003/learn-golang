package interfaces

import "github.com/gin-gonic/gin"

type ProductHandlerInterface interface {
	GetList() gin.HandlerFunc
	GetDetail() gin.HandlerFunc
	Create() gin.HandlerFunc
}
