package app

import (
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	authgrpc "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	authcasbin "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
)

func InitializeApp() (*App, func(), error) {
	cfg := config.Load()

	// 1. Casbin Enforcer (Centralized Writer)
	enforcer, err := authcasbin.NewEnforcer(cfg.DatabaseURL)
	if err != nil {
		return nil, nil, err
	}

	// 2. Base App
	baseApp := base.NewApp(base.Options{
		Name:   "auth-service",
		Config: cfg.BaseConfig,
	})

	// 3. Handlers
	handler := authgrpc.NewHandler(enforcer)

	// 4. Build App
	app := NewApp(baseApp, &cfg, enforcer, handler)

	cleanup := func() {}

	return app, cleanup, nil
}
