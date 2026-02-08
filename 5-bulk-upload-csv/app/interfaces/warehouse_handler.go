package interfaces

import "github.com/gin-gonic/gin"

type WarehouseHandlerInterface interface {
	GetList() gin.HandlerFunc
	GetDetail() gin.HandlerFunc
	Create() gin.HandlerFunc
}
