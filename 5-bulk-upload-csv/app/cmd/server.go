package cmd

import (
	"bulk-upload-csv/config"
	"bulk-upload-csv/interfaces"
	"log"

	"github.com/gin-gonic/gin"
)

const API_SERVER_DEFAULT_SERVER string = "8080"

type ApiServer struct {
	config config.Config

	mainHandler     interfaces.MainHandlerInterface
	generateHandler interfaces.GenerateHandlerInterface
	dbRepository    interfaces.DbRepositoryInterface

	categoryRepository interfaces.CategoryRepositoryInterface
	categoryService    interfaces.CategoryServiceInterface
	categoryHandler    interfaces.CategoryHandlerInterface

	warehouseRepository interfaces.WarehouseRepositoryInterface
	warehouseService    interfaces.WarehouseServiceInterface
	warehouseHandler    interfaces.WarehouseHandlerInterface

	productRepository interfaces.ProductRepositoryInterface
	productService    interfaces.ProductServiceInterface
	productHandler    interfaces.ProductHandlerInterface

	router *gin.Engine
}

func (server *ApiServer) Run() {
	if err := server.loadEnv(); err != nil {
		log.Fatal(err)
		return
	}

	server.router = gin.Default()
	server.dependenciesInjection()
	server.route()

	if err := server.router.Run(":" + API_SERVER_DEFAULT_SERVER); err != nil {
		log.Fatal(err)
	}
}
