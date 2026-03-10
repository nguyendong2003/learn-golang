package main

import (
	"elearning-api/cmd"
	docs "elearning-api/docs"
)

func main() {
	docs.SwaggerInfo.Title = "e-Smart Learn API"

	server := cmd.ApiServer{}
	server.Run()
}
