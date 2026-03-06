package cmd

import (
	"elearning-api/config"
	"elearning-api/handler"
	"elearning-api/repository"
	"elearning-api/util"
	"log"

	"github.com/gin-gonic/gin"
)

const API_SERVER_DEFAULT_SERVER string = "8080"

type ApiServer struct {
	config *config.Config

	dbRepository repository.DbRepository

	mainHandler handler.MainHandler
	userHandler handler.UserHandler
	authHandler handler.AuthHandler

	router *gin.Engine
}

func (server *ApiServer) Run() {
	// Initialize validator
	if err := util.InitValidator(); err != nil {
		log.Fatal(err)
		return
	}

	//
	server.router = gin.Default()

	// Apply middleware
	server.router.Use(server.RequestIDHandler())
	server.router.Use(server.CorsHandler())
	server.router.Use(server.ErrorHandler())

	server.dependenciesInjection()
	server.route()

	if err := server.router.Run(":" + API_SERVER_DEFAULT_SERVER); err != nil {
		log.Fatal(err)
	}
}
