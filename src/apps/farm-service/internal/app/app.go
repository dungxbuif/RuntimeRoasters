package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	farmgrpc "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/redis/go-redis/v9"
	farmv1 "github.com/dungxbuif/RuntimeRoasters/runtime/farm/v1"
	realcasbin "github.com/casbin/casbin/v3"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	DB          *database.DB
	VDB         *redis.Client
	KeyProvider provider.KeyProvider
	Enforcer    *realcasbin.SyncedEnforcer
	FarmHandler *farmgrpc.FarmHandler
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, vdb *redis.Client, keyProvider provider.KeyProvider, enforcer *realcasbin.SyncedEnforcer, farmHandler *farmgrpc.FarmHandler) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		DB:          db,
		VDB:         vdb,
		KeyProvider: keyProvider,
		Enforcer:    enforcer,
		FarmHandler: farmHandler,
	}
}

func (a *App) Run() {
	// Register gRPC
	a.Base.RegisterGRPC(&farmv1.FarmService_ServiceDesc, a.FarmHandler)

	a.Base.FinalizeRoutes()

	// Register readiness check
	a.Base.RegisterReadiness(func() error {
		if err := a.DB.Ping(context.Background()); err != nil {
			return err
		}
		return nil
	})

	if a.Cfg.GRPCPort == 0 {
		a.Cfg.GRPCPort = 50053
	}

	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
