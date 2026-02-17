package config

import "os"

type Config struct {
	AuthToken string
}

func Load() *Config {
	return &Config{
		AuthToken: os.Getenv("AUTH_TOKEN"),
	}
}
