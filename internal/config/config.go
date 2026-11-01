package config

import (
	"os"
)

type Config struct {
	Environment string
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func Load() *Config {

	return &Config{
		Environment: os.Getenv("Environment"),
		Database: DatabaseConfig{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
	}
}

