package cmd

import (
	"bulk-upload-csv/repository"
	"log"
)

func (server *ApiServer) dependenciesInjection() {
	server.dbRepository = repository.NewDbRepository(server.config.Database.Dsn)
	err := server.dbRepository.InitializeDB()
	if err != nil {
		log.Fatal(err)
	}

}
