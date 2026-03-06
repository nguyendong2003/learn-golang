package di

import (
	"elearning-api/config"
	"elearning-api/repository"
	"log"
)

type Container struct {
	DB       repository.DbRepository
	Repos    *Repositories
	Services *Services
}

func NewContainer(config *config.Config) *Container {
	db := repository.NewDbRepository(config.Database.Dsn)
	if err := db.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	repos := InitRepositories(db)
	services := InitServices(repos, config)

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: services,
	}
}
