package app

import (
	"context"
	"fmt"
	"time"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	Enforcer    *casbin.Enforcer
	Handler     *grpc.Handler
	UseCase     usecase.UserUsecase
	KeyProvider provider.KeyProvider
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, enforcer *casbin.Enforcer, handler *grpc.Handler, uc usecase.UserUsecase, keyProvider provider.KeyProvider) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		Enforcer:    enforcer,
		Handler:     handler,
		UseCase:     uc,
		KeyProvider: keyProvider,
	}
}

func (a *App) Run() {
	// 0. Bootstrapping Sync (Kratos -> Casbin)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()
		
		fmt.Println("Waiting for Kratos to be ready for sync...")
		time.Sleep(10 * time.Second)

		if err := a.UseCase.SyncCasbinWithKratos(ctx); err != nil {
			fmt.Printf("Initial bootstrapping sync FAILED: %v\n", err)
		}
	}()

	// Register gRPC
	a.Base.RegisterGRPC(&authv1.AuthService_ServiceDesc, a.Handler)

	a.Base.FinalizeRoutes()

	fmt.Println("Auth Service is running (gRPC Only)...")
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
