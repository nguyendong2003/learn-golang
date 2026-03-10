package di

import (
	"elearning-api/repository"
)

type Repositories struct {
	User repository.UserRepository
	Role repository.RoleRepository
}

func InitRepositories(db repository.DbRepository) *Repositories {
	return &Repositories{
		User: repository.NewUserRepository(db),
		Role: repository.NewRoleRepository(db),
	}
}
