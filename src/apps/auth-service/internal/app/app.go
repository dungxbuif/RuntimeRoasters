package app

import (
	"context"
	"time"

	svcconfig "RuntimeRoasters/apps/auth-service/config"
	"RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/auth/provider"
	"RuntimeRoasters/pkg/logger"
	authv1 "RuntimeRoasters/runtime/auth/v1"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"go.uber.org/zap"
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
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	// 0. Bootstrapping Sync (Kratos -> Casbin)
	go func() {
		log.Info("starting bootstrap sync with Kratos")
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			// Use a fresh context per attempt to avoid overall deadline issues
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if err := a.UseCase.SyncCasbinWithKratos(ctx); err == nil {
				log.Info("bootstrap sync completed")
				cancel()
				return
			} else {
				log.Warn("bootstrap sync attempt failed", zap.Error(err))
			}
			cancel()

			select {
			case <-ticker.C:
			}
		}
	}()

	// Register gRPC
	a.Base.RegisterGRPC(&authv1.AuthService_ServiceDesc, a.Handler)
	a.Base.RegisterGRPC(&systemv1.SystemService_ServiceDesc, a.Handler)

	// Register Gateway (REST -> gRPC Bridge)
	a.Base.RegisterGateway(authv1.RegisterAuthServiceHandlerFromEndpoint, a.Cfg.GRPCPort)
	a.Base.RegisterGateway(systemv1.RegisterSystemServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

	a.Base.FinalizeRoutes()

	log.Info("auth service is running")
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
