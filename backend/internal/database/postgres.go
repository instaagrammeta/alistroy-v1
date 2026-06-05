package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/instaagrammeta/alistroy-v1/backend/internal/config"
)

// Connect opens a GORM connection to PostgreSQL with sane production pool settings.
func Connect(cfg config.PostgresConfig, env string) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		PrepareStmt: true,
	}
	if env == "production" {
		gormCfg.Logger = logger.Default.LogMode(logger.Warn)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("postgres connect: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping: %w", err)
	}
	return db, nil
}
