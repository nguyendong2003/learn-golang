package cmd

import (
	"elearning-api/config"
	"elearning-api/handler"
	"elearning-api/pkg"
	"elearning-api/repository"
	"elearning-api/service"
	"log"
)

func (server *ApiServer) dependenciesInjection() {
	config, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	server.config = config
	server.initDatabase(config)
	server.initServices(config)
	server.initHandler()
}

func (s *ApiServer) initDatabase(config *config.Config) {
	db := repository.NewDbRepository(config.Database.Dsn)
	if err := db.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	s.userRepository = repository.NewUserRepository(db)
	s.roleRepository = repository.NewRoleRepository(db)
	s.blogRepository = repository.NewBlogRepository(db)
	s.refreshTokenRepository = repository.NewRefreshTokenRepository(db)
	s.categoryRepository = repository.NewCategoryRepository(db)
}

func (s *ApiServer) initServices(config *config.Config) {
	cache := pkg.NewCacheProvider(&config.Redis)
	mailer := pkg.NewEmailProvider(&config.Email, config.Frontend)

	s.userService = service.NewUserService(s.userRepository, s.roleRepository)
	s.authService = service.NewAuthService(s.userRepository, s.roleRepository, s.refreshTokenRepository, &config.JWT, mailer, cache)
	s.blogService = service.NewBlogService(s.blogRepository)
	s.categoryService = service.NewCategoryService(s.categoryRepository)
}
func (s *ApiServer) initHandler() {
	s.mainHandler = handler.NewMainHandler()
	s.userHandler = handler.NewUserHandler(s.userService)
	s.authHandler = handler.NewAuthHandler(s.authService)
	s.blogHandler = handler.NewBlogHandler(s.blogService)
	s.subscriptionHandler = handler.NewSubscriptionHandler()
	s.courseHandler = handler.NewCourseHandler()
	s.userCourseHandler = handler.NewUserCourseHandler()
	s.categoryHandler = handler.NewCategoryHandler(s.categoryService)
	s.feedbackHandler = handler.NewFeedbackHandler()
}
