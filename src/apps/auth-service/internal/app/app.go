package app

import (
	"context"
	"time"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
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
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		log.Info("starting bootstrap sync with Kratos")
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			if err := a.UseCase.SyncCasbinWithKratos(ctx); err == nil {
				log.Info("bootstrap sync completed")
				return
			} else {
				log.Warn("bootstrap sync attempt failed", zap.Error(err))
			}

			select {
			case <-ctx.Done():
				log.Warn("bootstrap sync deadline exceeded", zap.Error(ctx.Err()))
				return
			case <-ticker.C:
			}
		}
	}()

	// Register gRPC
	a.Base.RegisterGRPC(&authv1.AuthService_ServiceDesc, a.Handler)

	// Register Gateway (REST -> gRPC Bridge)
	a.Base.RegisterGateway(authv1.RegisterAuthServiceHandlerFromEndpoint, a.Cfg.GRPCPort)

	a.Base.FinalizeRoutes()

	log.Info("auth service is running")
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
