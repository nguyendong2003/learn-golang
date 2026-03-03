package cmd

func (server *ApiServer) route() {
	server.router.GET("/api/v1/health-check", server.mainHandler.HealthCheck())

	{
		userGroup := server.router.Group("/api/v1/users")
		userGroup.GET("/:id", server.userHandler.GetDetail())
		userGroup.POST("", server.userHandler.Create())
	}

}
