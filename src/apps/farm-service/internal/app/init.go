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
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/watcher"
	casbingrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func connnectDB(cfg *config.Config) (*database.DB, error) {
	return database.NewPostgres(database.PostgresConfig{
		URL: cfg.DatabaseURL,
	})
}

func connectVakey(cfg *config.Config) *redis.Client {
	return valkey.NewClient(valkey.Config{
		Addr: cfg.ValkeyAddr,
	})
}

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()

	db, err := connnectDB(&cfg)

	if err != nil {
		return nil, nil, err
	}

	vdb := connectVakey(&cfg)

	ttl, _ := time.ParseDuration(cfg.JWKSCacheTTL)
	keyProvider, err := provider.NewJWKSCache(cfg.JWKSURL, cfg.InternalSecret, ttl)
	if err != nil {
		return nil, nil, err
	}

	// 2. Auth & Casbin
	modelPath := "configs/rbac_model.conf"
	policyPath := "configs/rbac_policy.csv"

	casbinEnforcer, err := casbin.NewGormAdapterEnforcer(db.DB, modelPath)
	if err != nil {
		return nil, nil, err
	}

	watcher.WatchCasbinFiles(casbinEnforcer, modelPath, policyPath)


	// Optional: Seed policies from CSV if needed (only for development/initial setup)
	// You might want to wrap this in a condition or only run it once.
	// _ = casbinEnforcer.LoadPolicy() // Load from DB
	// casbinEnforcer.LoadPolicyFromCSV("rbac_policy.csv") 
	// casbinEnforcer.SavePolicy() // Sync back to DB if you want to persist CSV changes


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
