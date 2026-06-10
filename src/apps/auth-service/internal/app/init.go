package app

import (
	"context"
	"fmt"
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	authgrpc "RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	authcasbin "RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/repository"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/auth/provider"
	authinterceptor "RuntimeRoasters/pkg/base/auth/transport/grpc"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/valkey"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()
	if cfg.HydraAdminURL == "" {
		cfg.HydraAdminURL = "http://localhost:4445"
	}
	if cfg.KafkaPolicyTopic == "" {
		cfg.KafkaPolicyTopic = "auth.policy.changed"
	}

	// 1. Authoritative PostgreSQL Connection
	pgDb, err := database.NewPostgres(database.PostgresConfig{
		URL:          cfg.DatabaseURL,
		MaxOpenConns: 20,
		MaxIdleConns: 5,
		LogLevel:     cfg.LogLevel,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 2. Valkey Client connection
	var rClient *redis.Client
	if cfg.ValkeyAddr != "" {
		rClient = valkey.NewClient(valkey.Config{
			Addr: cfg.ValkeyAddr,
		})
	}

	// 3. Casbin Enforcer (Centralized Writer)
	enforcer, err := authcasbin.NewEnforcer(cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	// 4. Auth Provider (for verifying Admin JWT)
	ttl, _ := time.ParseDuration(cfg.JWKSCacheTTL)
	if ttl == 0 {
		ttl = 1 * time.Hour
	}
	keyProvider, err := provider.NewJWKSCache(cfg.JWKSURL, cfg.InternalSecret, ttl)
	if err != nil {
		return nil, nil, err
	}

	// 5. Kafka Producer
	producer := kafka.NewProducer(cfg.KafkaBrokers)

	// 6. Repositories & Usecases
	securityRepo := repository.NewSecurityStateRepository(pgDb)
	securityUsecase := usecase.NewSecurityUsecase(&cfg, securityRepo, pgDb, rClient, producer)
	userUsecase := usecase.NewUserUsecase(&cfg, enforcer, producer)

	localChecker := &localTokenValidityChecker{securityUsecase: securityUsecase}

	// 7. Base App
	baseApp := base.NewApp(base.Options{
		Name:   "auth-service",
		Config: cfg.BaseConfig,
		GRPCServerOptions: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				authinterceptor.GRPCUnaryInterceptor(keyProvider, cfg.ExpectedIssuer,
					authinterceptor.WithPublicRoutes(
						"/runtime.auth.v1.AuthService/AcceptLogin",     // Internal via Gateway
						"/runtime.auth.v1.AuthService/SeedSystem",      // CLI/Admin Seed
						"/runtime.system.v1.SystemService/SeedData",    // CLI/Admin Seed
					),
					authinterceptor.WithTokenValidityChecker(localChecker),
				),
			),
		},
	})

	// 8. Handlers & Build App
	handler := authgrpc.NewHandler(&cfg, enforcer, userUsecase, securityUsecase)
	app := NewApp(baseApp, &cfg, enforcer, handler, userUsecase, keyProvider)

	cleanup := func() {
		producer.Close()
		if rClient != nil {
			rClient.Close()
		}
	}

	return app, cleanup, nil
}

type localTokenValidityChecker struct {
	securityUsecase usecase.SecurityUsecase
}

func (c *localTokenValidityChecker) IsValid(ctx context.Context, subject string, iat time.Time) (bool, error) {
	return c.securityUsecase.ValidateTokenValidity(ctx, subject, iat)
}
