package cmd

import (
	"elearning-api/consts"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *ApiServer) route() {
	// Swagger UI
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		// authProtected.Use(server.AuthHandler())
		// {
		authProtected.POST("/change-password", server.authHandler.ChangePassword())
		// }
	}

	// Admin routes (protected)
	{
		adminGroup := server.router.Group("/api/v1/admin")

		adminGroup.Use(
			server.AuthHandler(),
			server.LoadRolePermissionHandler(),
			server.RequireRoleHandler(string(consts.RoleAdmin)),
		)

		adminUsers := adminGroup.Group("/users")
		{
			adminUsers.GET("", server.userHandler.GetList())
			adminUsers.GET("/:id", server.userHandler.GetByID())
			adminUsers.POST("", server.userHandler.Create())
			adminUsers.PUT("/:id", server.userHandler.Update())
			adminUsers.DELETE("/:id", server.userHandler.Delete())
		}
	}

	// User routes (protected)
	{
		userGroup := server.router.Group("/api/v1/users")
		// userGroup.Use(server.AuthHandler())

		userGroup.GET("", server.userHandler.GetList())
		userGroup.GET("/filter", server.userHandler.FilterAndPaginateAndSort())

		// Dynamic route must be placed after static route to avoid conflict
		userGroup.GET("/:id", server.userHandler.GetByID())

		userGroup.POST("", server.userHandler.Create())
		userGroup.PUT("/:id", server.userHandler.Update())
		userGroup.DELETE("/:id", server.userHandler.Delete())
	}

	// Blog routes (protected)
	{

		blogGroup := server.router.Group("/api/v1/blogs")
		blogGroup.GET("", server.blogHandler.GetList())
		// Dynamic route must be placed after static route to avoid conflict
		blogGroup.GET("/:slug", server.blogHandler.GetByID())

		blogGroup.Use(
			server.AuthHandler(),
			server.LoadRolePermissionHandler(),
		)

		blogGroup.POST("",
			server.RequirePermissionHandler("blog_create"),
			server.blogHandler.Create(),
		)

		blogGroup.PUT("/:id",
			server.RequirePermissionHandler("blog_update"),
			server.blogHandler.Update(),
		)

		blogGroup.DELETE("/:id",
			server.RequirePermissionHandler("blog_delete"),
			server.blogHandler.Delete(),
		)
	}

	// Subscription routes (public)
	{
		subGroup := server.router.Group("/api/v1/subscriptions")
		// Public endpoint returning mocked subscription plans
		subGroup.GET("", server.subscriptionHandler.GetSupcriptions())
	}

	// Course routes (public)
	{
		courseGroup := server.router.Group("/api/v1/courses")
		// public course listing and detail (mock)
		courseGroup.GET("", server.courseHandler.GetCourses())
		courseGroup.GET("/:slug", server.courseHandler.GetCourseBySlug())
	}

	// User course routes (mock)
	{
		meGroup := server.router.Group("/api/v1/users/me")
		// update progress and list enrolled courses
		meGroup.POST("/courses/progress", server.userCourseHandler.UpdateCourseProgress())
		meGroup.GET("/courses", server.userCourseHandler.GetMyCourses())
	}
}
