package config

import (
	"os"
	"strings"
)

type Config struct {
	Environment string
	Database    DatabaseConfig
	Kafka       KafkaConfig
}

type KafkaConfig struct {
	Brokers []string
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
		Kafka: KafkaConfig{
			Brokers: strings.FieldsFunc(os.Getenv("KAFKA_BROKERS"), func(r rune) bool {
				return r == ',' || r == ' ' || r == '\t' || r == '\n'
			}),
		},
		Database: DatabaseConfig{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			Name: os.Getenv("DB_NAME"),
		},
	}
}
