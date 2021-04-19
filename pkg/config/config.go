package config

type Config struct {
	AWSAccessKey       string `envconfig:"AWS_ACCESS_KEY" required:"true"`
	AWSSecretAccessKey string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
	AWSRegion          string `envconfig:"AWS_REGION" required:"true"`
	AWSKmsID           string `envconfig:"AWS_KMS_ID" required:"true"`
}

var Cfg = &Config{}
