package cmd

import (
	"bulk-upload-csv/config"
	"bulk-upload-csv/interfaces"
	"log"

	"github.com/gin-gonic/gin"
)

const API_SERVER_DEFAULT_SERVER string = "8080"

type ApiServer struct {
	config       config.Config
	dbRepository interfaces.DbRepositoryInterface

	router *gin.Engine
}

func (server *ApiServer) Run() {
	server.Start()
}

func (server *ApiServer) Start() {
	if err := server.loadEnv(); err != nil {
		log.Fatal(err)
		return
	}

	server.router = gin.Default()
}
