package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	farmgrpc "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	authhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/http"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	farmv1 "github.com/dungxbuif/RuntimeRoasters/runtime/farm/v1"
	"github.com/gin-gonic/gin"
	redisclient "github.com/redis/go-redis/v9"
)

type App struct {
	Base         *base.App
	Cfg          *svcconfig.Config
	DB           *database.DB
	RDB          *redisclient.Client
	KeyProvider  provider.KeyProvider
	CasbinEngine casbin.Engine
	FarmHandler  *farmgrpc.FarmHandler
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, rdb *redisclient.Client, keyProvider provider.KeyProvider, casbinEngine casbin.Engine, farmHandler *farmgrpc.FarmHandler) *App {
	return &App{
		Base:         baseApp,
		Cfg:          cfg,
		DB:           db,
		RDB:          rdb,
		KeyProvider:  keyProvider,
		CasbinEngine: casbinEngine,
		FarmHandler:  farmHandler,
	}
}

func (a *App) Run() {
	// 1. Start Casbin Background Sync (gRPC Snapshot + Kafka Live + Polling)
	if reader, ok := a.CasbinEngine.(*casbin.ResilientReader); ok {
		// Use a background context that is cancelled on app shutdown
		reader.StartBackgroundSync(context.Background())
	}

	// Register gRPC
	a.Base.RegisterGRPC(&farmv1.FarmService_ServiceDesc, a.FarmHandler)

	// Register Gateway
	a.Base.RegisterGateway(farmv1.RegisterFarmServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

	// Register Auth Middleware for native HTTP routes
	a.Base.RegisterHTTP(func(e *gin.Engine) {
		// Example: Protect all routes under /api/v1 (native Gin endpoints)
		// Note: Gateway routes /v1/* are protected by gRPC interceptor
		api := e.Group("/api")
		api.Use(authhttp.GinMiddleware(a.KeyProvider, a.Cfg.ExpectedIssuer))
	})

	// Register Swagger
	// Assuming running from src directory or properly handled path
	a.Base.ServeSwagger("/swagger", "./runtime/farm/v1")

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
