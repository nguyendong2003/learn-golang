package di

import (
	"elearning-api/repository"
)

type Repositories struct {
	User repository.UserRepository
}

func InitRepositories(db repository.DbRepository) *Repositories {
	return &Repositories{
		User: repository.NewUserRepository(db),
	}
}
