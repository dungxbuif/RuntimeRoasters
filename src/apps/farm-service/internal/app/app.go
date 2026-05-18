package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/config"
	farmgrpc "github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/event"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	rrcasbin "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	farmv1 "github.com/dungxbuif/RuntimeRoasters/runtime/farm/v1"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	DB          *database.DB
	VDB         *redis.Client
	KeyProvider provider.KeyProvider
	Enforcer    rrcasbin.Engine
	FarmHandler *farmgrpc.FarmHandler
	OutboxRelay *event.OutboxRelay
}

func NewApp(
	baseApp *base.App,
	cfg *svcconfig.Config,
	db *database.DB,
	vdb *redis.Client,
	keyProvider provider.KeyProvider,
	enforcer rrcasbin.Engine,
	farmHandler *farmgrpc.FarmHandler,
	outboxRelay *event.OutboxRelay,
) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		DB:          db,
		VDB:         vdb,
		KeyProvider: keyProvider,
		Enforcer:    enforcer,
		FarmHandler: farmHandler,
		OutboxRelay: outboxRelay,
	}
}

func (a *App) Run() {
	// Start Outbox Relay
	if a.OutboxRelay != nil {
		ctx := context.Background() // Base app handles context cancellation for server, but we can pass one here too
		go a.OutboxRelay.Start(ctx)
	}

	// Register gRPC
	a.Base.RegisterGRPC(&farmv1.FarmService_ServiceDesc, a.FarmHandler)

	// Register Gateway (REST -> gRPC Bridge)
	a.Base.RegisterGateway(farmv1.RegisterFarmServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

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
