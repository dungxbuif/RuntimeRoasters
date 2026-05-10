package app

import (
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	authgrpc "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	authhttp "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/http"
	authcasbin "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
)

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()

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
	handler := authgrpc.NewHandler(enforcer)
	userUsecase := usecase.NewUserUsecase(&cfg, enforcer, producer)
	userHandler := authhttp.NewUserHandler(userUsecase)

	// 6. Build App
	app := NewApp(baseApp, &cfg, enforcer, handler, userHandler, keyProvider)

	cleanup := func() {
		producer.Close()
	}

	return app, cleanup, nil
}
