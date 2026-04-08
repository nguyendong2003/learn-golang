package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"elearning-api/config"
	"elearning-api/handler"
	"elearning-api/job"
	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/service"
	"elearning-api/util"
	"elearning-api/worker"

	"github.com/hibiken/asynq"
)

func (server *ApiServer) dependenciesInjection() {
	logger := util.WithLayer(util.LayerBootstrap)
	config, err := NewConfig("config.yaml")
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	server.config = config
	server.initDatabase(config)
	server.initServices(config)
	server.initHandler()
}

func (s *ApiServer) initDatabase(config *config.Config) {
	logger := util.WithLayer(util.LayerRepository)
	db := repository.NewDbRepository(config.Database.Dsn)
	if err := db.InitializeDB(); err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	s.dbRepository = db
	s.userRepository = repository.NewUserRepository(db)
	s.roleRepository = repository.NewRoleRepository(db)
	s.blogRepository = repository.NewBlogRepository(db)
	s.planRepository = repository.NewPlanRepository(db)
	s.refreshTokenRepository = repository.NewRefreshTokenRepository(db)
	s.categoryRepository = repository.NewCategoryRepository(db)
	s.enrollmentRepository = repository.NewEnrollmentRepository(db)
	s.courseRepository = repository.NewCourseRepository(db)
	s.instructorProfileRepository = repository.NewInstructorProfileRepository(db)
	s.feedbackRepository = repository.NewFeedbackRepository(db)
	s.courseEventRepository = repository.NewCourseEventRepository(db)
	s.chapterRepository = repository.NewChapterRepository(db)
	s.lessonRepository = repository.NewLessonRepository(db)
	s.followRepository = repository.NewFollowRepository(db)
	s.subscriptionRepository = repository.NewSubscriptionRepository(db)
	s.paymentRepository = repository.NewPaymentRepository(db)
	s.subscriptionRevenueShareRepository = repository.NewSubscriptionRevenueShareRepository(db)
	s.stripeEventRepository = repository.NewStripeEventRepository(db)
	s.coursePurchaseRepository = repository.NewCoursePurchaseRepository(db)
	s.coursePurchaseDetailRepository = repository.NewCoursePurchaseDetailRepository(db)
	s.coursePurchaseRevenueShareRepository = repository.NewCoursePurchaseRevenueShareRepository(db)
	s.revenueRepository = repository.NewRevenueRepository(db)
	s.couponRepository = repository.NewCouponRepository(db)
	s.cartRepository = repository.NewCartRepository(db)
	s.presignUploadTrackingRepository = repository.NewPresignedUploadTrackingRepository(db)
}

func (s *ApiServer) initServices(config *config.Config) {
	logger := util.WithLayer(util.LayerInfra)
	cache := pkg.NewCacheProvider(&config.Redis)
	if cache == nil {
		logger.Error("failed to connect to redis")
		os.Exit(1)
	}
	mailer := pkg.NewEmailProvider(&config.Email, config.Frontend)
	if mailer == nil {
		logger.Error("failed to initialize email provider")
		os.Exit(1)
	}
	storage := pkg.NewStorageProvider(&config.Minio)
	if storage == nil {
		logger.Error("failed to initialize storage provider")
		os.Exit(1)
	}
	s.storageProvider = storage

	googleOAuthConfig := &config.OAuth.Google

	host := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	redisOpt := asynq.RedisClientOpt{Addr: host}
	s.asynqClient = asynq.NewClient(redisOpt)

	s.mailService = mailer
	s.userService = service.NewUserService(
		s.dbRepository,
		s.userRepository,
		s.roleRepository,
		s.blogRepository,
		s.instructorProfileRepository,
		s.enrollmentRepository,
		storage)
	s.authService = service.NewAuthService(
		s.userRepository,
		s.roleRepository,
		s.refreshTokenRepository,
		&config.JWT,
		googleOAuthConfig,
		mailer,
		cache)
	s.uploadService = service.NewUploadService(storage, s.presignUploadTrackingRepository)
	s.blogService = service.NewBlogService(
		s.dbRepository,
		s.blogRepository,
		s.categoryRepository,
		s.userRepository,
		s.asynqClient,
		s.followRepository,
		s.uploadService,
	)
	s.categoryService = service.NewCategoryService(s.categoryRepository)
	s.planService = service.NewPlanService(s.planRepository, s.subscriptionRepository, &config.Stripe)
	s.enrollmentService = service.NewEnrollmentService(
		s.enrollmentRepository,
		s.subscriptionRepository,
		s.paymentRepository,
		s.subscriptionRevenueShareRepository,
	)
	s.instructorProfileService = service.NewInstructorProfileService(s.instructorProfileRepository)
	s.courseService = service.NewCourseService(
		s.courseRepository,
		s.coursePurchaseRepository,
		s.categoryService,
		s.instructorProfileService,
		s.courseEventRepository,
		s.enrollmentRepository,
		s.userRepository,
		s.asynqClient,
		s.enrollmentService,
	)
	s.feedbackService = service.NewFeedbackService(s.feedbackRepository, s.userRepository, s.enrollmentRepository)
	s.revenueService = service.NewRevenueService(s.revenueRepository)
	s.lessonService = service.NewLessonService(s.lessonRepository, s.dbRepository, s.courseRepository, s.chapterRepository, s.uploadService)
	s.followService = service.NewFollowService(s.followRepository, s.userRepository)
	s.meetingService = service.NewMeetingService(s.courseEventRepository, s.courseRepository, s.userRepository)
	s.meetingHub = service.NewMeetingHub()
	s.couponService = service.NewCouponService(s.couponRepository, &config.Stripe)
	s.subscriptionService = service.NewSubscriptionService(
		s.userRepository,
		s.planRepository,
		s.courseRepository,
		s.enrollmentRepository,
		s.subscriptionRepository,
		s.paymentRepository,
		s.subscriptionRevenueShareRepository,
		s.coursePurchaseRevenueShareRepository,
		s.coursePurchaseRepository,
		s.coursePurchaseDetailRepository,
		s.couponRepository,
		s.stripeEventRepository,
		&config.Stripe,
	)
	s.cartService = service.NewCartService(
		s.cartRepository,
		s.userRepository,
		s.courseRepository,
		s.couponRepository,
		s.coursePurchaseRepository,
		s.coursePurchaseDetailRepository,
		&config.Stripe,
	)
}

func (s *ApiServer) initHandler() {
	s.mainHandler = handler.NewMainHandler()
	s.userHandler = handler.NewUserHandler(s.userService, s.subscriptionService)
	s.authHandler = handler.NewAuthHandler(s.authService, s.config.Frontend)
	s.blogHandler = handler.NewBlogHandler(s.blogService)
	s.userCourseHandler = handler.NewUserCourseHandler(s.enrollmentService)
	s.instructorProfileHandler = handler.NewInstructorProfileHandler(s.instructorProfileService)
	s.categoryHandler = handler.NewCategoryHandler(s.categoryService)
	s.feedbackHandler = handler.NewFeedbackHandler(s.feedbackService)
	s.lessonHandler = handler.NewLessonHandler(s.lessonService)
	s.fileUploadHandler = handler.NewFileUploadHandler(s.uploadService)
	s.userBlogHandler = handler.NewUserBlogHandler(s.userService)
	s.followHandler = handler.NewFollowHandler(s.followService)
	s.meetingHandler = handler.NewMeetingHandler(s.meetingService, s.meetingHub, &s.config.JWT)
	s.revenueHandler = handler.NewRevenueHandler(s.revenueService)
	s.courseHandler = handler.NewCourseHandler(
		s.courseService,
		s.subscriptionService,
		s.instructorProfileService,
		s.uploadService,
		s.enrollmentService,
	)
	s.planHandler = handler.NewPlanHandler(s.planService)
	s.subscriptionHandler = handler.NewSubscriptionHandler(s.subscriptionService)
	s.couponHandler = handler.NewCouponHandler(s.couponService)
	s.cartHandler = handler.NewCartHandler(s.cartService)
}

func (s *ApiServer) initBackgroundJob() {
	logger := util.WithLayer(util.LayerWorker)
	host := fmt.Sprintf("%s:%d", s.config.Redis.Host, s.config.Redis.Port)
	redisOpt := asynq.RedisClientOpt{Addr: host}
	srv := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
	})

	emailWorkerHandler := worker.NewEmailWorkerHandler(
		s.courseEventRepository,
		s.userRepository,
		s.mailService,
		s.asynqClient,
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(job.TypeEventNotification, emailWorkerHandler.HandleEventNotification)
	mux.HandleFunc(job.TypeSendSingleEmail, emailWorkerHandler.HandleSendEmail)

	blogWorkerHandler := worker.NewBlogScheduledWorkerHandler(s.blogRepository)
	mux.HandleFunc(job.TypeBlogPublish, blogWorkerHandler.HandlePublish)

	cleanupWorker := worker.NewPresignedUploadCleanupWorker(s.presignUploadTrackingRepository, s.storageProvider)
	go cleanupWorker.Start(context.Background())

	go func() {
		if err := srv.Run(mux); err != nil {
			logger.Error("could not run asynq server", "error", err)
			os.Exit(1)
		}
	}()

	go s.startStripeSubscriptionSyncCron()
	go s.startStripeCoursePurchaseSyncCron()
}

func (s *ApiServer) startStripeSubscriptionSyncCron() {
	logger := util.WithLayer(util.LayerWorker)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logger.Info("started stripe subscription sync cron", "interval", (30 * time.Second).String())

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		if err := s.subscriptionService.SyncPendingStripeSubscriptions(ctx); err != nil {
			logger.Error("stripe subscription sync failed", "error", err)
		}
		cancel()

		<-ticker.C
	}
}

func (s *ApiServer) startStripeCoursePurchaseSyncCron() {
	logger := util.WithLayer(util.LayerWorker)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	logger.Info("started stripe course purchase sync cron", "interval", (30 * time.Second).String())

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		if err := s.subscriptionService.SyncPendingStripeCoursePurchases(ctx); err != nil {
			logger.Error("stripe course purchase sync failed", "error", err)
		}
		cancel()

		<-ticker.C
	}
}
