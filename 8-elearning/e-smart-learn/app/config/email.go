package config

type EmailConfig struct {
	SMTPHost string `validate:"required" yaml:"SMTP_HOST"`
	SMTPPort int    `validate:"required" yaml:"SMTP_PORT"`
	SMTPUser string `validate:"required" yaml:"SMTP_USER"`
	SMTPPass string `validate:"required" yaml:"SMTP_PASS"`
}