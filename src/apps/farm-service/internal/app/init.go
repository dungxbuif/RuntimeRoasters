package app

import (
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	farmgrpc "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/repository"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	authgrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	casbingrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/valkey"
	"google.golang.org/grpc"
)

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()

	// 1. Core Infra
	db, err := database.NewPostgres(database.PostgresConfig{
		URL: cfg.DatabaseURL,
	})
	if err != nil {
		return nil, nil, err
	}

	vdb := valkey.NewClient(valkey.Config{
		Addr: cfg.ValkeyAddr,
	})

	// 2. Auth & Casbin
	ttl, _ := time.ParseDuration(cfg.JWKSCacheTTL)
	keyProvider, err := provider.NewJWKSCache(cfg.JWKSURL, cfg.InternalSecret, ttl)
	if err != nil {
		return nil, nil, err
	}

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
	casbinEnforcer, err := casbin.NewGormAdapterEnforcer(db.DB, modelText)
	if err != nil {
		return nil, nil, err
	}

	// 3. Server Options & Base App
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			authgrpc.GRPCUnaryInterceptor(keyProvider, cfg.ExpectedIssuer),
			casbingrpc.GRPCUnaryInterceptor(casbinEnforcer),
		),
	}

	baseApp := base.NewApp(base.Options{
		Name:              "farm-service",
		Config:            cfg.BaseConfig,
		GRPCServerOptions: grpcOpts,
	})

	// 4. App-specific Components
	farmRepo := repository.NewFarmRepository(db, casbinEnforcer)

	farmUsecase := usecase.NewFarmUsecase(farmRepo)
	farmHandler := farmgrpc.NewFarmHandler(farmUsecase)

	// 5. Build App
	app := NewApp(baseApp, &cfg, db, vdb, keyProvider, casbinEnforcer, farmHandler)

	cleanup := func() {
		// No manual cleanup required for db/rdb here as shutdown handles graceful termination
	}

	return app, cleanup, nil
}
