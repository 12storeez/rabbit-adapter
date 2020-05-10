package config

import (
	"github.com/subosito/gotenv"
	"os"
)

type Config struct {
	Rabbit   *Rabbit
	App *App
}

type App struct {
	Port    string
	Version string
	ApiKey  string
}

type Rabbit struct {
	Url string
}

func New() *Config {
	_ = gotenv.Load(".env")
	return &Config{
		Rabbit: &Rabbit{
			Url: os.Getenv("RABBIT_URL"),
		},
		App: &App{
			Port:    os.Getenv("PORT"),
			Version: os.Getenv("VERSION"),
		},
	}
}
