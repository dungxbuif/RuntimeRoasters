package database

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type PostgresConfig struct {
	URL          string
	MaxOpenConns int    // default: 20
	MaxIdleConns int    // default: 5
	LogLevel     string // silent, error, warn, info
}

type DB struct {
	*gorm.DB
}

type TxManager interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

// NewPostgres initializes and returns a wrapped gorm.DB instance connected to Postgres
func NewPostgres(cfg PostgresConfig) (*DB, error) {
	gormCfg := &gorm.Config{
		Logger: newGormLogger(cfg.LogLevel),
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

// WithTx helps manage transactions in a GORM-idiomatic way.
// It wraps the context with the transaction instance so repositories can extract it.
func (db *DB) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// GetTx extracts the gorm.DB transaction from the context if it exists,
// otherwise returns the default database instance with the original context.
func (db *DB) GetTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return db.DB.WithContext(ctx)
}

// Ping checks database availability
func (db *DB) Ping(ctx context.Context) error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func parseLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Info
	}
}

func newGormLogger(level string) gormlogger.Interface {
	return gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  parseLogLevel(level),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
