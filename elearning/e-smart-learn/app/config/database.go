package config

type Database struct {
	Dsn string `validate:"required"`
}
