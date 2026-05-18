package app

import (
	"bufio"
	"os"
	"strings"
	"time"

	realcasbin "github.com/casbin/casbin/v3"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	farmgrpc "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/event"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/repository"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	authgrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	casbingrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/watcher"
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

	if err := seedCasbinPolicies(casbinEnforcer, policyPath); err != nil {
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

	// 4. App-specific	// 5. Repositories
	farmRepo := repository.NewFarmRepository(db, casbinEnforcer)
	harvestRepo := repository.NewHarvestRepository(db, casbinEnforcer)
	outboxRepo := repository.NewOutboxRepository(db)

	// 6. Use Cases
	publisher := event.NewMockPublisher()
	farmUC := usecase.NewFarmUsecase(farmRepo)
	harvestUC := usecase.NewHarvestUsecase(db, harvestRepo, farmRepo, outboxRepo)

	// 7. Workers
	outboxRelay := event.NewOutboxRelay(outboxRepo, publisher, 5*time.Second)

	// 8. Handlers & Server
	farmHandler := farmgrpc.NewFarmHandler(farmUC, harvestUC)

	// 5. Build App
	app := NewApp(baseApp, &cfg, db, vdb, keyProvider, casbinEnforcer, farmHandler, outboxRelay)

	cleanup := func() {
		// No manual cleanup required for db/rdb here as shutdown handles graceful termination
	}

	return app, cleanup, nil
}

func seedCasbinPolicies(enforcer *realcasbin.SyncedEnforcer, policyPath string) error {
	policies, err := enforcer.GetPolicy()
	if err != nil {
		return err
	}
	groupingPolicies, err := enforcer.GetGroupingPolicy()
	if err != nil {
		return err
	}
	if len(policies) > 0 || len(groupingPolicies) > 0 {
		return nil
	}

	file, err := os.Open(policyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if len(parts) < 3 {
			continue
		}

		switch parts[0] {
		case "p":
			if len(parts) < 4 {
				continue
			}
			if _, err := enforcer.AddPolicy(parts[1], parts[2], parts[3]); err != nil {
				return err
			}
		case "g":
			if _, err := enforcer.AddGroupingPolicy(parts[1], parts[2]); err != nil {
				return err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return enforcer.SavePolicy()
}
