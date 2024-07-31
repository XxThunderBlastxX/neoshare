package config

import "github.com/XxThunderBlastxX/goconfigenv"

type AppConfig struct {
	Port int `env:"PORT,default=8080"`
	Auth AuthConfig
}

type AuthConfig struct {
	ClientId     string `env:"AUTH_CLIENT_ID"`
	ClientSecret string `env:"AUTH_CLIENT_SECRET"`
	CallbackURL  string `env:"AUTH_CALLBACK_URL"`
	Domain       string `env:"AUTH_DOMAIN"`
}

func New() *AppConfig {
	config, err := goconfigenv.ParseEnv[AppConfig]()
	if err != nil {
		panic(err)
	}
	return &config
}
