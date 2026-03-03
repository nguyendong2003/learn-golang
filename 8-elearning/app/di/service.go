package di

import "elearning-api/service"

type Services struct {
	User service.UserService
}

func InitServices(repos *Repositories) *Services {
	return &Services{
		User: service.NewUserService(repos.User),
	}
}
