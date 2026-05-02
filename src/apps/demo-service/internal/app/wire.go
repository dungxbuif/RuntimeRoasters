//go:build wireinject
// +build wireinject

package app

import (
	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/config"
	demogrpc "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/redis"
	"github.com/google/wire"
)

func provideConfigPtr(cfg svcconfig.Config) *svcconfig.Config {
	return &cfg
}

func provideBaseOptions(cfg svcconfig.Config) base.Options {
	return base.Options{
		Name:   "demo-service",
		Config: cfg.BaseConfig,
	}
}

func providePostgresConfig(cfg svcconfig.Config) database.PostgresConfig {
	return database.PostgresConfig{
		URL: cfg.DatabaseURL,
	}
}

func provideRedisConfig(cfg svcconfig.Config) redis.Config {
	return redis.Config{
		Addr: cfg.RedisAddr,
	}
}

func InitializeApp() (*App, func(), error) {
	wire.Build(
		svcconfig.Load,
		provideConfigPtr,
		provideBaseOptions,
		providePostgresConfig,
		provideRedisConfig,
		base.NewApp,
		database.NewPostgres,
		redis.NewClient,
		usecase.NewDemoUsecase,
		demogrpc.NewDemoHandler,
		NewApp,
	)
	return &App{}, nil, nil
}
