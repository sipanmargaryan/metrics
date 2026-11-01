package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"metrics/internal/config"
)

func ConnecPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
