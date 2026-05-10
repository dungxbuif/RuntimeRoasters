package app

import (
	"fmt"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
)

type App struct {
	Base     *base.App
	Cfg      *svcconfig.Config
	Enforcer *casbin.Enforcer
	Handler  *grpc.Handler
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, enforcer *casbin.Enforcer, handler *grpc.Handler) *App {
	return &App{
		Base:     baseApp,
		Cfg:      cfg,
		Enforcer: enforcer,
		Handler:  handler,
	}
}

func (a *App) Run() {
	// Register gRPC
	a.Base.RegisterGRPC(&authv1.AuthService_ServiceDesc, a.Handler)

	fmt.Println("Auth Service is running...")
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
