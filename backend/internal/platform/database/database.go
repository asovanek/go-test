package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open returns a configured GORM connection.
func Open(dsn string, silent bool) (*gorm.DB, error) {
	cfg := &gorm.Config{}
	if silent {
		cfg.Logger = logger.Default.LogMode(logger.Silent)
	}
	db, err := gorm.Open(postgres.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
