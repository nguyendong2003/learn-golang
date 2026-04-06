package config

type StripeConfig struct {
	SecretKey                        string `validate:"required" yaml:"secret_key"`
	WebhookSecret                    string `validate:"required" yaml:"webhook_secret"`
	SuccessURL                       string `validate:"required" yaml:"success_url"`
	CancelURL                        string `validate:"required" yaml:"cancel_url"`
	BillingPortal                    string `validate:"required" yaml:"billing_portal_return_url"`
	DefaultCurrency                  string `validate:"required" yaml:"default_currency"`
	CheckoutSessionExpirationSeconds int    `yaml:"checkout_session_expiration_seconds" default:"1800"`
}
