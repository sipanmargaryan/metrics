package config

import "gorm.io/gorm"

type App struct {
	DB     *gorm.DB
	Config *Config
}
