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
	authgrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	casbingrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/grpc"
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

func provideCasbinClient(cfg svcconfig.Config) (casbin.AuthSnapshotClient, error) {
	return casbingrpc.NewAuthSnapshotClient(cfg.AuthServiceAddr, "demo-service")
}

func provideCasbinEngine(client casbin.AuthSnapshotClient) (casbin.Engine, error) {
	// In a real app, modelText would come from a shared config or file
	modelText := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act))
`
	return casbin.NewResilientReader(client, casbin.ReaderOptions{
		ModelText:     modelText,
		SyncInterval:  10 * time.Minute,
		RetryInterval: 5 * time.Second,
	})
}

func provideGRPCServerOptions(keyProvider provider.KeyProvider, casbinEngine casbin.Engine, cfg svcconfig.Config) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			authgrpc.GRPCUnaryInterceptor(keyProvider, cfg.ExpectedIssuer),
			casbingrpc.GRPCUnaryInterceptor(casbinEngine),
		),
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
		provideCasbinClient,
		provideCasbinEngine,
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
