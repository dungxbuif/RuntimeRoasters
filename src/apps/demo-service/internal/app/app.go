package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/config"
	demogrpc "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	demov1 "github.com/dungxbuif/RuntimeRoasters/runtime/demo/v1"
	redisclient "github.com/redis/go-redis/v9"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	DB          *database.DB
	RDB         *redisclient.Client
	DemoHandler *demogrpc.DemoHandler
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, rdb *redisclient.Client, demoHandler *demogrpc.DemoHandler) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		DB:          db,
		RDB:         rdb,
		DemoHandler: demoHandler,
	}
}

func (a *App) Run() {
	// Register gRPC
	a.Base.RegisterGRPC(&demov1.DemoService_ServiceDesc, a.DemoHandler)

	// Register Gateway
	a.Base.RegisterGateway(demov1.RegisterDemoServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

	// Register Swagger
	// Assuming running from src directory or properly handled path
	a.Base.ServeSwagger("/swagger", "./runtime/demo/v1")

	// Register readiness check
	a.Base.RegisterReadiness(func() error {
		if err := a.DB.Ping(context.Background()); err != nil {
			return err
		}
		if err := a.RDB.Ping(context.Background()).Err(); err != nil {
			return err
		}
		return nil
	})

	if a.Cfg.AppPort == 0 {
		a.Cfg.AppPort = 8080
	}
	if a.Cfg.GRPCPort == 0 {
		a.Cfg.GRPCPort = 50051
	}

	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
