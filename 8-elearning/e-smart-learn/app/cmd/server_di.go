package cmd

import (
	"elearning-api/di"
	"elearning-api/handler"
	"log"
)

func (server *ApiServer) dependenciesInjection() {
	config, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	server.config = config
	container := di.NewContainer(server.config)

	server.dbRepository = container.DB
	server.mainHandler = handler.NewMainHandler()
	server.userHandler = handler.NewUserHandler(container.Services.User)
	server.authHandler = handler.NewAuthHandler(container.Services.Auth)
	server.blogHandler = handler.NewBlogHandler()
	server.subscriptionHandler = handler.NewSubscriptionHandler()
	server.courseHandler = handler.NewCourseHandler()
	server.userCourseHandler = handler.NewUserCourseHandler()
}
