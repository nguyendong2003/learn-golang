package config

type Config struct {
	Database Database    `validate:"required" yaml:"database"`
	JWT      JWTConfig   `validate:"required" yaml:"jwt"`
	Email    EmailConfig `validate:"required" yaml:"email"`
	Redis    RedisConfig `validate:"required" yaml:"redis"`
	Frontend string      `validate:"required" yaml:"frontend"`
}
