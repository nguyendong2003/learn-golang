package config

type RedisConfig struct {
	Host     string `validate:"required" yaml:"host"`
	Port     int    `validate:"required" yaml:"port"`
	Password string `validate:"required" yaml:"password"`
	DB       int    `validate:"required" yaml:"db"`
}
