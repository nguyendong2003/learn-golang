package cmd

import (
	"os"

	"github.com/joho/godotenv"
)

func (server *ApiServer) loadEnv() error {
	if err := godotenv.Load(); err!= nil {
		return err
	}

	dsn := os.Getenv("DB_CONNECTION")
	server.config.Database.Dsn = dsn

	return nil
}
