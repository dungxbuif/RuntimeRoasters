package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/config"
	demogrpc "github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	authhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/http"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	demov1 "github.com/dungxbuif/RuntimeRoasters/runtime/demo/v1"
	"github.com/gin-gonic/gin"
	redisclient "github.com/redis/go-redis/v9"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	DB          *database.DB
	RDB         *redisclient.Client
	KeyProvider provider.KeyProvider
	DemoHandler *demogrpc.DemoHandler
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, rdb *redisclient.Client, keyProvider provider.KeyProvider, demoHandler *demogrpc.DemoHandler) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		DB:          db,
		RDB:         rdb,
		KeyProvider: keyProvider,
		DemoHandler: demoHandler,
	}
}

func (a *App) Run() {
	// Register gRPC
	a.Base.RegisterGRPC(&demov1.DemoService_ServiceDesc, a.DemoHandler)

	// Register Gateway
	a.Base.RegisterGateway(demov1.RegisterDemoServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

	// Register Auth Middleware for native HTTP routes
	a.Base.RegisterHTTP(func(e *gin.Engine) {
		// Example: Protect all routes under /api/v1 (native Gin endpoints)
		// Note: Gateway routes /v1/* are protected by gRPC interceptor
		api := e.Group("/api")
		api.Use(authhttp.GinMiddleware(a.KeyProvider, a.Cfg.ExpectedIssuer))
	})

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
