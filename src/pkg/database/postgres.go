package database

import (
	"context"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresConfig struct {
	URL          string
	MaxOpenConns int // default: 20
	MaxIdleConns int // default: 5
	LogLevel     string // silent, error, warn, info
}

type DB struct {
	*gorm.DB
}

// NewPostgres initializes and returns a wrapped gorm.DB instance connected to Postgres
func NewPostgres(cfg PostgresConfig) (*DB, error) {
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(parseLogLevel(cfg.LogLevel)),
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), gormCfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 20
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &DB{db}, nil
}

// WithTx helps manage transactions in a GORM-idiomatic way
func (db *DB) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

// Ping checks database availability
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func parseLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}
