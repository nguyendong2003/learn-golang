package di

import (
	"elearning-api/repository"
	"log"
)

type Container struct {
	DB       repository.DbRepository
	Repos    *Repositories
	Services *Services
}

func NewContainer(db repository.DbRepository) *Container {
	repos := InitRepositories(db)
	services := InitServices(repos)

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: services,
	}
}

func MustInitialize(dsn string) *Container {
	db := repository.NewDbRepository(dsn)
	if err := db.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return NewContainer(db)
}
