package configs

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

