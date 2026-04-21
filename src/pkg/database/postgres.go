package database

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type PostgresConfig struct {
	URL          string
	MaxOpenConns int // default: 20
	MaxIdleConns int // default: 5
}

type DB struct {
	*sqlx.DB
}

// NewPostgres initializes and returns a wrapped sqlx.DB instance connected to Postgres
func NewPostgres(cfg PostgresConfig) (*DB, error) {
	db, err := sqlx.Connect("pgx", cfg.URL)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 20
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Hour)

	return &DB{db}, nil
}

// WithTx helps manage transactions, particularly for the Outbox Pattern
func (db *DB) WithTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Ping checks database availability
func (db *DB) Ping(ctx context.Context) error {
	return db.PingContext(ctx)
}
