//go:build wireinject
// +build wireinject

package app

import (
	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/google/wire"
)

func provideConfigPtr(cfg svcconfig.Config) *svcconfig.Config {
	return &cfg
}

func provideCasbinEnforcer(cfg svcconfig.Config) (*casbin.Enforcer, error) {
	return casbin.NewEnforcer(cfg.DatabaseURL)
}

func provideBaseOptions(cfg svcconfig.Config) base.Options {
	return base.Options{
		Name:   "auth-service",
		Config: cfg.BaseConfig,
	}
}

func InitializeApp() (*App, func(), error) {
	wire.Build(
		svcconfig.Load,
		provideConfigPtr,
		provideCasbinEnforcer,
		provideBaseOptions,
		base.NewApp,
		grpc.NewHandler,
		NewApp,
	)
	return &App{}, nil, nil
}
