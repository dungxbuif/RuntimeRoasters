package app

import (
	"context"
	"os"
	"time"

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
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func connnectDB(cfg *config.Config) (*database.DB, error) {
	return database.NewPostgres(database.PostgresConfig{
		URL:      cfg.DatabaseURL,
		LogLevel: cfg.DBLogLevel,
	})
}

func connectVakey(cfg *config.Config) *redis.Client {
	return valkey.NewClient(valkey.Config{
		Addr: cfg.ValkeyAddr,
	})
}

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()
	if cfg.DBLogLevel == "" {
		cfg.DBLogLevel = "warn"
	}
	if cfg.KafkaHarvestTopic == "" {
		cfg.KafkaHarvestTopic = "farm.harvest.events"
	}
	if cfg.KafkaPolicyTopic == "" {
		cfg.KafkaPolicyTopic = "auth.policy.changed"
	}

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
	modelBytes, err := os.ReadFile(modelPath)
	if err != nil {
		return nil, nil, err
	}

	authSnapshotClient, err := casbingrpc.NewAuthSnapshotClient(cfg.AuthServiceAddr, cfg.AppName)
	if err != nil {
		return nil, nil, err
	}

	policyConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.AppName+"-auth-policy", cfg.KafkaPolicyTopic)
	casbinEnforcer, err := casbin.NewResilientReader(authSnapshotClient, casbin.ReaderOptions{
		ModelText:            string(modelBytes),
		SyncInterval:         5 * time.Minute,
		RetryInterval:        2 * time.Second,
		PolicyChangeConsumer: policyConsumer,
	})
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
	casbinEnforcer.StartBackgroundSync(context.Background())

	// 4. App-specific	// 5. Repositories
	farmRepo := repository.NewFarmRepository(db, casbinEnforcer)
	harvestRepo := repository.NewHarvestRepository(db, casbinEnforcer)
	outboxRepo := repository.NewOutboxRepository(db)

	// 6. Use Cases
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	publisher := event.NewKafkaPublisher(producer, cfg.KafkaHarvestTopic)
	farmUC := usecase.NewFarmUsecase(farmRepo)
	harvestUC := usecase.NewHarvestUsecase(db, harvestRepo, farmRepo, outboxRepo)

	// 7. Workers
	outboxRelay := event.NewOutboxRelay(outboxRepo, publisher, 5*time.Second)

	// 8. Handlers & Server
	farmHandler := farmgrpc.NewFarmHandler(farmUC, harvestUC)

	// 5. Build App
	app := NewApp(baseApp, &cfg, db, vdb, keyProvider, casbinEnforcer, farmHandler, outboxRelay)

	cleanup := func() {
		_ = producer.Close()
		_ = vdb.Close()
		if sqlDB, sqlErr := db.DB.DB(); sqlErr == nil {
			_ = sqlDB.Close()
		}
	}

	return app, cleanup, nil
}
