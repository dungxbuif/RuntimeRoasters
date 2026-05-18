package app

import (
	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/config"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
)

func connectDB(cfg *config.Config) (*database.DB, error) {
	return database.NewPostgres(database.PostgresConfig{
		URL:      cfg.DatabaseURL,
		LogLevel: cfg.DBLogLevel,
	})
}

func connectValkey(cfg *config.Config) *redis.Client {
	return valkey.NewClient(valkey.Config{
		Addr: cfg.ValkeyAddr,
	})
}

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()
	if cfg.DBLogLevel == "" {
		cfg.DBLogLevel = "warn"
	}

	db, err := connectDB(&cfg)
	if err != nil {
		return nil, nil, err
	}

	vdb := connectValkey(&cfg)

	baseApp := base.NewApp(base.Options{
		Name:   "logistics-service",
		Config: cfg.BaseConfig,
	})

	app := NewApp(baseApp, &cfg, db, vdb)

	cleanup := func() {
		_ = vdb.Close()
		if sqlDB, sqlErr := db.DB.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
	}

	return app, cleanup, nil
}
