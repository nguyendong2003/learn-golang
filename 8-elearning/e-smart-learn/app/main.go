package main

import (
	"elearning-api/cmd"
	docs "elearning-api/docs"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	docs.SwaggerInfo.Title = "e-Smart Learn API"
	docs.SwaggerInfo.Description = "API for Elearning platform"
	docs.SwaggerInfo.Version = "1.0"

	server := cmd.ApiServer{}
	server.Run()
}
