package database

import (
	"fmt"
	"time"

	"github.com/khaingminhtun/relio-backend/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(
	cfg config.DatabaseConfig,
) (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf(
			"open postgres connection: %w",
			err,
		)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf(
			"postgres ping connection : %w",
			err,
		)
	}

	// connection pool

	sqlDB.SetMaxOpenConns(10)

	sqlDB.SetMaxIdleConns(5)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
