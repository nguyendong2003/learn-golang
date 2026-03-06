package cmd

func (server *ApiServer) route() {
	server.router.GET("/api/v1/health-check", server.mainHandler.HealthCheck())

	// Auth routes (public)
	{
		authGroup := server.router.Group("/api/v1/auth")

		authGroup.POST("/login", server.authHandler.Login())
		authGroup.POST("/register", server.authHandler.Register())
		authGroup.POST("/refresh-token", server.authHandler.RefreshToken())
		authGroup.POST("/forgot-password", server.authHandler.ForgotPassword())

		// Protected auth routes
		authProtected := authGroup.Group("")
		authProtected.Use(server.AuthHandler())
		{
			authProtected.POST("/change-password", server.authHandler.ChangePassword())
		}
	}

	// User routes (protected)
	{
		userGroup := server.router.Group("/api/v1/users")
		userGroup.Use(server.AuthHandler())

		userGroup.GET("", server.userHandler.GetList())
		userGroup.GET("/filter", server.userHandler.FilterAndPaginateAndSort())

		// Dynamic route must be placed after static route to avoid conflict
		userGroup.GET("/:id", server.userHandler.GetByID())

		userGroup.POST("", server.userHandler.Create())
		userGroup.PUT("/:id", server.userHandler.Update())
		userGroup.DELETE("/:id", server.userHandler.Delete())
	}

}
