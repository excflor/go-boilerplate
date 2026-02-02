package database

import (
	"fmt"
	"go-boilerplate/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	connString := "host=%s port=%s user=%s dbname=%s password=%s"
	dsn := fmt.Sprintf(connString,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Name,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
