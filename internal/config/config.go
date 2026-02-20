package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nurgal1ev/yotabo-go/internal/infrastructure/postgres"
	"log"
	"sync"
)

type Config struct {
	App      AppConfig
	Postgres postgres.Config
}

type AppConfig struct {
	AuthToken string `env:"APP_AUTH_TOKEN" env-required:"true"`
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadConfig("deployment/.env", instance); err != nil {
			if err := cleanenv.ReadEnv(instance); err != nil {
				log.Fatal(err)
			}
		}
	})

	return instance
}
