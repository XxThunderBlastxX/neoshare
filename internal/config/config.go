package config

import "github.com/XxThunderBlastxX/goconfigenv"

type AppConfig struct {
	Port int `env:"PORT,default=8080"`
}

func New() *AppConfig {
	config, err := goconfigenv.ParseEnv[AppConfig]()
	if err != nil {
		panic(err)
	}
	return &config
}
