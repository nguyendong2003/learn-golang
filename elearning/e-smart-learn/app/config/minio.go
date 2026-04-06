package config

type MinioConfig struct {
	Endpoint         string `yaml:"endpoint" validate:"required"`
	ExternalEndpoint string `yaml:"external_endpoint" validate:"required"`
	AccessKey        string `yaml:"access_key" validate:"required"`
	SecretKey        string `yaml:"secret_key" validate:"required"`
	UseSSL           bool   `yaml:"use_ssl" validate:"required"`
	ExternalUseSSL   *bool  `yaml:"external_use_ssl"`
	BucketName       string `yaml:"bucket_name" validate:"required"`
}
