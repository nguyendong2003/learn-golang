package cmd

import (
	"elearning-api/consts"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *ApiServer) route() {
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server.router.GET("/api/v1/health-check", server.mainHandler.HealthCheck())
	server.router.GET("/api/v1/events/:event_id/meeting/ws", server.meetingHandler.JoinEventMeeting())

	// Auth routes
	{
		server.router.GET("/api/v1/oauth/google/login", server.authHandler.GoogleLogin())
		server.router.GET("/api/v1/oauth/google/callback", server.authHandler.GoogleCallback())

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

	// Me route
	{
		meGroup := server.router.Group("/api/v1/users/me")
		meGroup.Use(server.AuthHandler())
		meGroup.Use(server.LoadRolePermissionHandler())
		meGroup.GET("", server.userHandler.GetMe())
		meGroup.PUT("", server.userHandler.UpdateMe())
		meGroup.POST("/courses/progress", server.userCourseHandler.UpdateCourseProgress())
		meGroup.GET("/courses", server.userCourseHandler.GetMyCourses())
		meGroup.GET("/courses/taught", server.RequireRoleHandler(string(consts.RoleInstructor)), server.courseHandler.GetMyTaughtCourses())
		meGroup.POST("/blogs/read", server.userBlogHandler.MarkBlogAsRead())
		meGroup.GET("/blogs/read", server.userBlogHandler.GetReadBlogs())
		meGroup.POST("/apply-instructor", server.userHandler.ApplyToInstructor())

		meGroup.GET("/revenue/overview", server.RequireInstructorProfileHandler(), server.revenueHandler.GetInstructorRevenueOverview())
		meGroup.GET("/revenue/statistics", server.RequireInstructorProfileHandler(), server.revenueHandler.GetInstructorStatistics())
		meGroup.GET("/revenue/statistics/day", server.RequireInstructorProfileHandler(), server.revenueHandler.GetInstructorStatisticsByDay())
		meGroup.GET("/revenue/statistics/month", server.RequireInstructorProfileHandler(), server.revenueHandler.GetInstructorStatisticsByMonth())
		meGroup.GET("/revenue/statistics/year", server.RequireInstructorProfileHandler(), server.revenueHandler.GetInstructorStatisticsByYear())
		meGroup.GET("/transactions", server.userHandler.GetTransactionsHistory())
	}

	// Follow routes
	{
		followGroup := server.router.Group("/api/v1/users")
		followGroup.GET("/:id/followers", server.followHandler.GetFollowers())
		followGroup.GET("/:id/followings", server.followHandler.GetFollowings())

		followProtected := followGroup.Group("")
		followProtected.Use(server.AuthHandler())
		followProtected.POST("/:id/follow", server.followHandler.FollowUser())
		followProtected.DELETE("/:id/follow", server.followHandler.UnfollowUser())
	}

	// Blog routes
	{
		blogGroup := server.router.Group("/api/v1/blogs")
		blogGroup.GET("", server.blogHandler.GetPublishedBlogs())
		blogGroup.GET("/:slug",
			server.OptionalAuthHandler(),
			server.LoadRolePermissionOptionalHandler(),
			server.blogHandler.GetBySlug())
		blogAuth := blogGroup.Group("")
		{
			blogAuth.Use(
				server.AuthHandler(),
				server.LoadRolePermissionHandler(),
			)
			blogAuth.POST("",
				server.RequirePermissionHandler("blog_create"),
				server.blogHandler.Create(),
			)
			blogAuth.PUT("/:id",
				server.RequirePermissionHandler("blog_update"),
				server.blogHandler.Update(),
			)
			blogAuth.DELETE("/:id",
				server.RequirePermissionHandler("blog_delete"),
				server.blogHandler.Delete(),
			)
		}
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
		courseGroup.Use(server.OptionalAuthHandler())
		courseGroup.GET("", server.courseHandler.GetList())
		courseGroup.GET("/:id", server.courseHandler.GetByID())
		courseGroup.GET("/slug/:slug",
			server.OptionalAuthHandler(),
			server.LoadRolePermissionOptionalHandler(),
			server.courseHandler.GetBySlug(),
		)
		courseGroup.GET("/:id/events", server.courseHandler.GetEvents())

		courseBuyer := courseGroup.Group("")
		courseBuyer.Use(server.AuthHandler())
		courseBuyer.POST("/:id/purchase-checkout-session", server.courseHandler.CreatePurchaseCheckoutSession())
		courseBuyer.POST("/:id/purchase-checkout-preview", server.courseHandler.PreviewPurchaseCheckout())

		courseProtected := courseGroup.Group("")
		courseProtected.Use(server.AuthHandler(), server.LoadRolePermissionHandler())
		courseProtected.POST("/:id/enroll",
			server.RequireActiveSubscriptionHandler(),
			server.courseHandler.EnrollInCourse(),
		)
		courseProtected.POST("",
			server.RequirePermissionHandler("course_create"),
			server.courseHandler.Create(),
		)
		courseProtected.PUT("/:id",
			server.RequirePermissionHandler("course_update"),
			server.courseHandler.Update(),
		)
		courseProtected.PUT("/:id/update-status",
			server.RequirePermissionHandler("course_update"),
			server.LoadCourseHandler(),
			server.RequireCourseOwnerHandler(),
			server.courseHandler.UpdateStatus(),
		)
		courseProtected.DELETE("/:id",
			server.RequirePermissionHandler("course_delete"),
			server.LoadCourseHandler(),
			server.RequireCourseOwnerHandler(),
			server.courseHandler.Delete(),
		)
		//
		courseProtected.POST("/:id/publish",
			server.RequirePermissionHandler("course_update"),
			server.LoadCourseHandler(),
			server.RequireCourseOwnerHandler(),
			server.courseHandler.PublishCourse(),
		)

		courseProtected.POST("/:id/submit",
			server.RequireInstructorProfileHandler(),
			server.LoadCourseHandler(),
			server.RequireCourseOwnerHandler(),
			server.courseHandler.SubmitForReview(),
		)
		//
		courseProtected.POST("/:id/events",
			server.RequirePermissionHandler("event_create"),
			server.courseHandler.CreateEvent())
		courseProtected.PUT("/:id/events/:event_id",
			server.RequirePermissionHandler("event_update"),
			server.courseHandler.UpdateEvent())
		courseProtected.DELETE("/:id/events/:event_id",
			server.RequirePermissionHandler("event_delete"),
			server.courseHandler.DeleteEvent())

		courseAdmin := courseProtected.Group("/admin")
		courseAdmin.Use(server.RequireRoleHandler(string(consts.RoleAdmin)))
		{
			courseAdmin.GET("/statistics", server.courseHandler.GetStatistics())
		}
	}
	// Course purchase verify by session
	server.router.GET("/api/v1/course-purchase/:session_id", server.courseHandler.GetCoursePurchaseBySessionID())

	// Admin review endpoints for courses
	server.router.POST("/api/v1/admin/courses/:id/approve",
		server.AuthHandler(),
		server.LoadRolePermissionHandler(),
		server.RequireRoleHandler(string(consts.RoleAdmin)),
		server.LoadCourseHandler(),
		server.courseHandler.ApproveCourse(),
	)

	server.router.POST("/api/v1/admin/courses/:id/reject",
		server.AuthHandler(),
		server.LoadRolePermissionHandler(),
		server.RequireRoleHandler(string(consts.RoleAdmin)),
		server.LoadCourseHandler(),
		server.courseHandler.RejectCourse(),
	)

	// Coupon route
	{
		instructorCouponGroup := server.router.Group("/api/v1/instructor/coupons")
		instructorCouponGroup.Use(
			server.AuthHandler(),
			server.LoadRolePermissionHandler(),
			server.RequireRoleHandler(string(consts.RoleInstructor)),
			server.RequireInstructorProfileHandler(),
		)
		instructorCouponGroup.GET("", server.couponHandler.GetList())
		instructorCouponGroup.GET("/assignable", server.couponHandler.GetAssignableList())
		instructorCouponGroup.GET("/assignable-for-course", server.couponHandler.GetAssignableListForCourse())
		instructorCouponGroup.GET("/:id", server.couponHandler.GetByID())
		instructorCouponGroup.POST("", server.couponHandler.Create())
		instructorCouponGroup.PUT("/:id/deactivate", server.couponHandler.Deactivate())
	}

	// Cart route
	{
		cartGroup := server.router.Group("/api/v1/carts")
		cartGroup.Use(server.AuthHandler())
		cartGroup.GET("", server.cartHandler.GetMyCart())
		cartGroup.POST("/courses/:course_id", server.cartHandler.AddCourse())
		cartGroup.DELETE("/courses/:course_id", server.cartHandler.RemoveCourse())
		cartGroup.POST("/checkout-preview", server.cartHandler.PreviewCheckout())
		cartGroup.POST("/checkout-session", server.cartHandler.Checkout())
	}

	// Instructor profile route
	{
		instructorProfileGroup := server.router.Group("/api/v1/instructor-profiles")
		instructorProfileGroup.GET("", server.instructorProfileHandler.GetList())
		instructorProfileGroup.GET("/:id", server.instructorProfileHandler.GetByID())

		instructorProtected := instructorProfileGroup.Group("")
		instructorProtected.Use(server.AuthHandler(), server.LoadRolePermissionHandler())
		instructorProtected.POST("",
			server.RequirePermissionHandler("instructor_profile_create"),
			server.instructorProfileHandler.Create(),
		)
		instructorProtected.PUT("/:id",
			server.RequirePermissionHandler("instructor_profile_update"),
			server.RequireInstructorProfileOwnerHandler(),
			server.instructorProfileHandler.Update(),
		)
		instructorProtected.DELETE("/:id",
			server.RequirePermissionHandler("instructor_profile_delete"),
			server.RequireInstructorProfileOwnerHandler(),
			server.instructorProfileHandler.Delete(),
		)
	}

	// Lesson route
	{
		lessonGroup := server.router.Group("/api/v1/courses/:id")
		lessonGroup.GET("/lessons", server.lessonHandler.GetByCourseID())
		lessonProtected := lessonGroup.Group("")
		lessonProtected.Use(server.AuthHandler(), server.LoadRolePermissionHandler())
		lessonProtected.POST("/lessons",
			server.RequirePermissionHandler("lesson_create"),
			server.LoadCourseHandler(),
			server.lessonHandler.Create(),
		)
		lessonProtected.PUT("/lessons",
			server.RequirePermissionHandler("lesson_update"),
			server.LoadCourseHandler(),
			server.lessonHandler.UpdateLessons(),
		)
	}

	// Plan route
	{
		planGroup := server.router.Group("/api/v1/subscription-plans")
		planGroup.GET("/active", server.planHandler.GetActivePlans())
		planGroup.GET("/:id", server.planHandler.GetByID())
	}

	// Subscription route
	{
		subGroup := server.router.Group("/api/v1/subscriptions")
		subGroup.POST("/webhook/stripe", server.subscriptionHandler.WebhookStripe())

		subProtected := subGroup.Group("")
		subProtected.Use(server.AuthHandler())
		subProtected.POST("/checkout-session", server.subscriptionHandler.CreateSubscriptionCheckoutSession())
		subProtected.GET("/me", server.subscriptionHandler.GetMySubscription())

		subProtected.POST("/cancel", server.subscriptionHandler.CancelAtPeriodEnd())
		subProtected.POST("/resume", server.subscriptionHandler.Resume())
		subProtected.POST("/billing-portal", server.subscriptionHandler.CreateBillingPortalSession())
	}

	// Feedback route
	{
		server.router.GET("/api/v1/feedbacks/featured", server.feedbackHandler.GetFeaturedFeedbacks())
		feedbackGroup := server.router.Group("/api/v1/courses/:id/feedbacks")
		feedbackGroup.GET("", server.feedbackHandler.GetFeedbacks())
		feedbackAuth := feedbackGroup.Group("")
		feedbackAuth.Use(server.AuthHandler())
		feedbackAuth.POST("", server.feedbackHandler.CreateFeedback())
	}

	// File upload
	{
		uploadGroup := server.router.Group("/api/v1/upload")
		uploadGroup.POST("/image", server.fileUploadHandler.UploadImage())
		uploadGroup.POST("/presign", server.fileUploadHandler.PresignUploadURL())
	}

	{
		adminGroup := server.router.Group("/api/v1/admin")
		{
			adminGroup.Use(
				server.AuthHandler(),
				server.LoadRolePermissionHandler(),
				server.RequireRoleHandler(string(consts.RoleAdmin)),
			)

			// User administration
			adminUsers := adminGroup.Group("/users")
			{
				adminUsers.GET("", server.userHandler.GetList())
				adminUsers.GET("/:id", server.userHandler.GetByID())
				adminUsers.GET("/statistics", server.userHandler.GetStatistics())
				adminUsers.GET("/statistics/active-students", server.userHandler.GetActiveStudentStatistics())
			}

			// Application administration
			adminApplications := adminGroup.Group("/instructor-applications")
			{
				adminApplications.GET("", server.userHandler.GetPendingInstructorApplications())
				adminApplications.POST("/:id/approve", server.userHandler.ApproveInstructorApplication())
				adminApplications.POST("/:id/reject", server.userHandler.RejectInstructorApplication())
			}

			// Blog administration
			adminBlogs := adminGroup.Group("/blogs")
			adminBlogs.GET("",
				server.blogHandler.GetBlogs())
			adminBlogs.GET("/statistics",
				server.blogHandler.GetStatistics())
			adminBlogs.GET("/:id",
				server.blogHandler.GetByID())

			subAdmin := adminGroup.Group("/subscriptions")
			subAdmin.GET("/subscribers", server.subscriptionHandler.GetSubscribers())

			adminPlans := adminGroup.Group("/subscription-plans")
			adminPlans.Use(server.AuthHandler(), server.LoadRolePermissionHandler())
			adminPlans.GET("",
				server.planHandler.GetList(),
			)
			adminPlans.POST("",
				server.planHandler.Create(),
			)
			adminPlans.PUT("/:id",
				server.planHandler.Update(),
			)
			adminPlans.PUT("/:id/activate",
				server.planHandler.Activate(),
			)
			adminPlans.PUT("/:id/deactivate",
				server.planHandler.Deactivate(),
			)
			adminPlans.DELETE("/:id",
				server.planHandler.Delete(),
			)

			// Course administration
			adminCourses := adminGroup.Group("/courses")
			adminCourses.GET("/statistics/new", server.courseHandler.GetNewCoursesLast30Days())

			// Subscription administration
			adminSubscription := adminGroup.Group("/subscriptions")
			adminSubscription.GET("/statistics/members/retention", server.subscriptionHandler.GetMemberRetention())

			// Instructor profile administration
			adminInstructorProfile := adminGroup.Group("/instructor-profiles")
			adminInstructorProfile.GET("/statistics/teachers/growth", server.instructorProfileHandler.GetGrowthStatistics())

			// Revenue and growth statistics
			adminRevenue := adminGroup.Group("/revenue")
			adminRevenue.GET("/overview", server.revenueHandler.GetAdminRevenueOverview())
			adminRevenue.GET("/statistics", server.revenueHandler.GetAdminStatistics())
			adminRevenue.GET("/statistics/day", server.revenueHandler.GetAdminStatisticsByDay())
			adminRevenue.GET("/statistics/month", server.revenueHandler.GetAdminStatisticsByMonth())
			adminRevenue.GET("/statistics/year", server.revenueHandler.GetAdminStatisticsByYear())

			adminRevenue.GET("/statistics/sales-segmentation", server.revenueHandler.GetAdminSalesSegmentation())

			adminRevenue.GET("/transactions", server.revenueHandler.GetAdminTransactions())

			adminRevenue.GET("/statistics/teachers/revenue", server.revenueHandler.GetAllTeachersRevenue())
		}
	}
}
