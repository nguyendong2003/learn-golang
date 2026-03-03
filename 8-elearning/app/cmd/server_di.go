package cmd

import (
	"elearning-api/di"
	"elearning-api/handler"
)

func (server *ApiServer) dependenciesInjection() {
	container := di.MustInitialize(server.config.Database.Dsn)

	server.dbRepository = container.DB
	server.mainHandler = handler.NewMainHandler()
	server.userHandler = handler.NewUserHandler(container.Services.User)
}
