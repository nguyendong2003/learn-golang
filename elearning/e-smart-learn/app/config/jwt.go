package config

import "time"

type JWTConfig struct {
	AccessTokenSecret      string        `validate:"required" yaml:"access_token"`
	RefreshTokenSecret     string        `validate:"required" yaml:"refresh_token"`
	AccessTokenExpiration  time.Duration `validate:"required" yaml:"access_token_expiration"`
	RefreshTokenExpiration time.Duration `validate:"required" yaml:"refresh_token_expiration"`
	Issuer                 string        `validate:"required" yaml:"issuer"`
}
