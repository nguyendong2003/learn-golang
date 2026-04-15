package cmd

import (
	"os"

	"elearning-api/config"
	"elearning-api/handler"
	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/service"
	"elearning-api/util"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

const API_SERVER_DEFAULT_SERVER string = "8080"

type ApiServer struct {
	config *config.Config

	dbRepository repository.DbRepository

	userRepository                       repository.UserRepository
	blogRepository                       repository.BlogRepository
	roleRepository                       repository.RoleRepository
	planRepository                       repository.PlanRepository
	refreshTokenRepository               repository.RefreshTokenRepository
	categoryRepository                   repository.CategoryRepository
	enrollmentRepository                 repository.EnrollmentRepository
	courseRepository                     repository.CourseRepository
	instructorProfileRepository          repository.InstructorProfileRepository
	feedbackRepository                   repository.FeedbackRepository
	courseEventRepository                repository.CourseEventRepository
	chapterRepository                    repository.ChapterRepository
	lessonRepository                     repository.LessonRepository
	followRepository                     repository.FollowRepository
	subscriptionRepository               repository.SubscriptionRepository
	paymentRepository                    repository.PaymentRepository
	subscriptionRevenueShareRepository   repository.SubscriptionRevenueShareRepository
	stripeEventRepository                repository.StripeEventRepository
	coursePurchaseRepository             repository.CoursePurchaseRepository
	coursePurchaseDetailRepository       repository.CoursePurchaseDetailRepository
	coursePurchaseRevenueShareRepository repository.CoursePurchaseRevenueShareRepository
	revenueRepository                    repository.RevenueRepository
	couponRepository                     repository.CouponRepository
	courseCouponRepository               repository.CourseCouponRepository
	cartRepository                       repository.CartRepository
	presignUploadTrackingRepository      repository.PresignedUploadTrackingRepository

	userService              service.UserService
	authService              service.AuthService
	blogService              service.BlogService
	uploadService            service.UploadService
	categoryService          service.CategoryService
	courseService            service.CourseService
	instructorProfileService service.InstructorProfileService
	planService              service.PlanService
	enrollmentService        service.EnrollmentService
	feedbackService          service.FeedbackService
	revenueService           service.RevenueService
	mailService              pkg.EmailProvider
	storageProvider          pkg.StorageProvider
	asynqClient              *asynq.Client
	lessonService            service.LessonService
	followService            service.FollowService
	subscriptionService      service.SubscriptionService
	couponService            service.CouponService
	cartService              service.CartService
	meetingService           service.MeetingService
	meetingHub               *service.MeetingHub

	mainHandler              handler.MainHandler
	userHandler              handler.UserHandler
	authHandler              handler.AuthHandler
	blogHandler              handler.BlogHandler
	courseHandler            handler.CourseHandler
	userCourseHandler        handler.UserCourseHandler
	categoryHandler          handler.CategoryHandler
	feedbackHandler          handler.FeedbackHandler
	instructorProfileHandler handler.InstructorProfileHandler
	lessonHandler            handler.LessonHandler
	fileUploadHandler        handler.FileUploadHandler
	userBlogHandler          handler.UserBlogHandler
	followHandler            handler.FollowHandler
	meetingHandler           handler.MeetingHandler
	revenueHandler           handler.RevenueHandler
	planHandler              handler.PlanHandler
	subscriptionHandler      handler.SubscriptionHandler
	couponHandler            handler.CouponHandler
	cartHandler              handler.CartHandler

	router *gin.Engine
}

func (server *ApiServer) Run() {
	util.SetupLoggerFromEnv()
	logger := util.WithLayer(util.LayerApp)

	// Initialize validator
	if err := util.InitValidator(); err != nil {
		logger.Error("failed to initialize validator", "error", err)
		os.Exit(1)
		return
	}
	server.router = gin.Default()

	// Apply middleware
	server.router.Use(server.RequestIDHandler())
	server.router.Use(server.RequestLoggingHandler())
	server.router.Use(server.CorsHandler())
	server.router.Use(server.CaptureRequestBody())
	server.router.Use(server.ErrorHandler())

	server.dependenciesInjection()
	server.route()
	server.initBackgroundJob()

	if err := server.router.Run(":" + API_SERVER_DEFAULT_SERVER); err != nil {
		logger.Error("failed to run http server", "port", API_SERVER_DEFAULT_SERVER, "error", err)
		os.Exit(1)
	}
}
