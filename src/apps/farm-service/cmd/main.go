package main

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/redis"
)

func main() {
	// 1. Load config
	cfg := svcconfig.Load()

	// 2. Bootstrap framework
	app := base.NewApp(base.Options{
		Name:   "farm-service",
		Config: cfg.BaseConfig,
	})

	// 3. Infrastructure (Connect DB + Redis)
	db, err := database.NewPostgres(database.PostgresConfig{URL: cfg.DatabaseURL})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(redis.Config{Addr: cfg.RedisAddr})

	// 4. Register readiness check
	app.RegisterReadiness(func() error {
		if err := db.Ping(context.Background()); err != nil {
			return err
		}
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			return err
		}
		return nil
	})

	// 5. Run service
	if cfg.AppPort == 0 {
		cfg.AppPort = 8080 // safe default
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 50051 // safe default
	}
	
	app.Run(cfg.AppPort, cfg.GRPCPort)
}
