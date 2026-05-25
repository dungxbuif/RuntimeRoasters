package app

import (
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	authgrpc "RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	authcasbin "RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/auth/provider"
	"RuntimeRoasters/pkg/kafka"
)

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()
	if cfg.HydraAdminURL == "" {
		cfg.HydraAdminURL = "http://localhost:4445"
	}
	if cfg.KafkaPolicyTopic == "" {
		cfg.KafkaPolicyTopic = "auth.policy.changed"
	}

	// 1. Casbin Enforcer (Centralized Writer)
	enforcer, err := authcasbin.NewEnforcer(cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	// 2. Auth Provider (for verifying Admin JWT)
	ttl, _ := time.ParseDuration(cfg.JWKSCacheTTL)
	if ttl == 0 {
		ttl = 1 * time.Hour
	}
	keyProvider, err := provider.NewJWKSCache(cfg.JWKSURL, cfg.InternalSecret, ttl)
	if err != nil {
		return nil, nil, err
	}

	// 3. Kafka Producer
	producer := kafka.NewProducer(cfg.KafkaBrokers)

	// 4. Base App
	baseApp := base.NewApp(base.Options{
		Name:   "auth-service",
		Config: cfg.BaseConfig,
	})

	// 5. Handlers & Usecases
	userUsecase := usecase.NewUserUsecase(&cfg, enforcer, producer)
	handler := authgrpc.NewHandler(enforcer, userUsecase)

	// 6. Build App
	app := NewApp(baseApp, &cfg, enforcer, handler, userUsecase, keyProvider)

	cleanup := func() {
		producer.Close()
	}

	return app, cleanup, nil
}
