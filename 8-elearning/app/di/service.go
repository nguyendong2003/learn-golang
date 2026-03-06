package di

import (
	"elearning-api/config"
	"elearning-api/service"
)

type Services struct {
	User service.UserService
	Auth service.AuthService
}

func InitServices(repos *Repositories, config *config.Config) *Services {
	return &Services{
		User: service.NewUserService(repos.User),
		Auth: service.NewAuthService(repos.User, &config.JWT),
	}
}
