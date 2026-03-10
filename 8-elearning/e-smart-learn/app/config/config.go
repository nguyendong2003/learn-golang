package config

type Config struct {
	Database Database  `validate:"required" yaml:"database"`
	JWT      JWTConfig `validate:"required" yaml:"jwt"`
}
