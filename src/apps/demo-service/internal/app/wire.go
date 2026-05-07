//go:build wireinject
// +build wireinject

package app

import (
	"time"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/config"
	demogrpc "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/redis"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func provideConfigPtr(cfg svcconfig.Config) *svcconfig.Config {
	return &cfg
}

func provideKeyProvider(cfg svcconfig.Config) (provider.KeyProvider, error) {
	ttl, _ := time.ParseDuration(cfg.JWKSCacheTTL)
	return provider.NewJWKSCache(cfg.JWKSURL, cfg.InternalSecret, ttl)
}

func provideGRPCServerOptions(keyProvider provider.KeyProvider, cfg svcconfig.Config) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(authgrpc.GRPCUnaryInterceptor(keyProvider, cfg.ExpectedIssuer)),
	}
}

func provideBaseOptions(cfg svcconfig.Config, grpcOpts []grpc.ServerOption) base.Options {
	return base.Options{
		Name:              "demo-service",
		Config:            cfg.BaseConfig,
		GRPCServerOptions: grpcOpts,
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
		provideKeyProvider,
		provideGRPCServerOptions,
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
