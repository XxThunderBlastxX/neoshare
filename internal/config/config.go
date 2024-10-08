package config

import "github.com/XxThunderBlastxX/goconfigenv"

type AppConfig struct {
	AppEnv     string `env:"APP_ENV,default=dev"` // This is either dev or prod (default: dev)
	Port       int    `env:"PORT,default=8080"`
	FaviconURI string `env:"FAVICON_URI,default=https://koustav.dev/favicon.ico"`
	Auth       AuthConfig
	S3Config   S3Config
	DBConfig   DBConfig
}

type AuthConfig struct {
	ClientID          string `env:"AUTH_CLIENT_ID"`
	ClientSecret      string `env:"AUTH_CLIENT_SECRET"`
	CallbackURL       string `env:"AUTH_CALLBACK_URL"`
	Domain            string `env:"AUTH_DOMAIN"`
	Audience          string `env:"AUTH_AUDIENCE"`
	UserInfoURL       string `env:"AUTH_USER_INFO_URL"`
	LogoutCallbackURL string `env:"AUTH_LOGOUT_CALLBACK_URL"`
}

type S3Config struct {
	Endpoint  string `env:"S3_ENDPOINT"`
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Bucket    string `env:"S3_BUCKET"`
}

type DBConfig struct {
	DBUser     string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_DATABASE"`
	DBPort     string `env:"DB_PORT"`
	DBHost     string `env:"DB_HOST"`
}

func New() *AppConfig {
	config, err := goconfigenv.ParseEnv[AppConfig]()
	if err != nil {
		panic(err)
	}
	return &config
}
