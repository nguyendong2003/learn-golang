package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	Google GoogleOAuthConfig `yaml:"google" validate:"required"`
}

type GoogleOAuthConfig struct {
	ClientID     string `yaml:"client_id" validate:"required"`
	ClientSecret string `yaml:"client_secret" validate:"required"`
	RedirectURL  string `yaml:"redirect_url" validate:"required"`
}

func (c *GoogleOAuthConfig) OAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}
