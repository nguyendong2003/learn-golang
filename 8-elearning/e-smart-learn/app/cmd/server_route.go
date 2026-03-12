package cmd

import (
	"elearning-api/consts"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *ApiServer) route() {
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router.GET("/api/v1/health-check", server.mainHandler.HealthCheck())

	// Auth routes
	{
		authGroup := server.router.Group("/api/v1/auth")
		authGroup.POST("/login", server.authHandler.Login())
		authGroup.POST("/register", server.authHandler.Register())
		authGroup.POST("/refresh-token", server.authHandler.RefreshToken())
		authGroup.POST("/forgot-password", server.authHandler.ForgotPassword())
		authGroup.POST("/reset-password", server.authHandler.ResetPassword())
		authProtected := authGroup.Group("")
		authProtected.Use(server.AuthHandler())
		authProtected.PUT("/change-password", server.authHandler.ChangePassword())
	}

	// User route
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

		userGroup := server.router.Group("/api/v1/users")
		{
			userGroup.GET("", server.userHandler.GetList())
			userGroup.GET("/filter", server.userHandler.FilterAndPaginateAndSort())
			userGroup.GET("/:id", server.userHandler.GetByID())
			userGroup.POST("", server.userHandler.Create())
			userGroup.PUT("/:id", server.userHandler.Update())
			userGroup.DELETE("/:id", server.userHandler.Delete())
		}
	}

	// Blog routes
	{
		blogGroup := server.router.Group("/api/v1/blogs")
		blogGroup.GET("", server.blogHandler.GetList())
		blogGroup.GET("/:slug", server.blogHandler.GetByID())

		blogGroup.Use(
			server.AuthHandler(),
			// server.LoadRolePermissionHandler(),
		)

		blogGroup.POST("",
			// server.RequirePermissionHandler("blog_create"),
			server.blogHandler.Create(),
		)

		blogGroup.PUT("/:id",
			// server.RequirePermissionHandler("blog_update"),
			server.blogHandler.Update(),
		)

		blogGroup.DELETE("/:id",
			// server.RequirePermissionHandler("blog_delete"),
			server.blogHandler.Delete(),
		)
	}

	// Category route
	{
		categoryGroup := server.router.Group("/api/v1/categories")
		categoryGroup.GET("/all", server.categoryHandler.GetAll())
		categoryGroup.GET("", server.categoryHandler.GetList())
		categoryGroup.GET("/:id", server.categoryHandler.GetByID())

		categoryProtected := categoryGroup.Group("")
		categoryProtected.Use(server.AuthHandler(), server.LoadRolePermissionHandler())
		categoryProtected.POST("", server.RequirePermissionHandler("category_create"), server.categoryHandler.Create())
		categoryProtected.PUT("/:id", server.RequirePermissionHandler("category_update"), server.categoryHandler.Update())
		categoryProtected.DELETE("/:id", server.RequirePermissionHandler("category_delete"), server.categoryHandler.Delete())
	}

	// Course route
	{
		courseGroup := server.router.Group("/api/v1/courses")
		courseGroup.GET("", server.courseHandler.GetCourses())
		courseGroup.GET("/:slug", server.courseHandler.GetCourseBySlug())
	}

	// Subscription routes (mocked)
	{
		subGroup := server.router.Group("/api/v1/subscriptions")
		subGroup.GET("", server.subscriptionHandler.GetPlans())
	}

	{
		feedbackGroup := server.router.Group("/api/v1/feedbacks")
		// feedbackGroup.Use(server.AuthHandler())

		feedbackGroup.GET("", server.feedbackHandler.GetFeedbacks())
	}

	{
		meGroup := server.router.Group("/api/v1/users/me")
		// update progress and list enrolled courses
		meGroup.POST("/courses/progress", server.userCourseHandler.UpdateCourseProgress())
		meGroup.GET("/courses", server.userCourseHandler.GetMyCourses())
	}
}
