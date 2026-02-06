package cmd

import "github.com/joho/godotenv"

func (server *ApiServer) loadEnv() error {
	err := godotenv.Load()
	return err
}
