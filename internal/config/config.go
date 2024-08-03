package config

import "github.com/XxThunderBlastxX/goconfigenv"

type AppConfig struct {
	Port     int `env:"PORT,default=8080"`
	Auth     AuthConfig
	S3Config S3Config
}

type AuthConfig struct {
	ClientId     string `env:"AUTH_CLIENT_ID"`
	ClientSecret string `env:"AUTH_CLIENT_SECRET"`
	CallbackURL  string `env:"AUTH_CALLBACK_URL"`
	Domain       string `env:"AUTH_DOMAIN"`
}

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT"`
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Bucket    string `env:"S3_BUCKET"`
}

func New() *AppConfig {
	config, err := goconfigenv.ParseEnv[AppConfig]()
	if err != nil {
		panic(err)
	}
	return &config
}
