package app

import (
	"fmt"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	authhttp "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/http"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
	baseauthhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/http"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
	"github.com/gin-gonic/gin"
)

type App struct {
	Base        *base.App
	Cfg         *svcconfig.Config
	Enforcer    *casbin.Enforcer
	Handler     *grpc.Handler
	UserHandler *authhttp.UserHandler
	KeyProvider provider.KeyProvider
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, enforcer *casbin.Enforcer, handler *grpc.Handler, userHandler *authhttp.UserHandler, keyProvider provider.KeyProvider) *App {
	return &App{
		Base:        baseApp,
		Cfg:         cfg,
		Enforcer:    enforcer,
		Handler:     handler,
		UserHandler: userHandler,
		KeyProvider: keyProvider,
	}
}

func (a *App) Run() {
	// Register gRPC
	a.Base.RegisterGRPC(&authv1.AuthService_ServiceDesc, a.Handler)

	// Register HTTP
	a.Base.RegisterHTTP(func(e *gin.Engine) {
		v1 := e.Group("/v1/admin")
		
		// Auth Middleware for Admin APIs
		v1.Use(baseauthhttp.GinMiddleware(a.KeyProvider, a.Cfg.ExpectedIssuer))
		
		v1.POST("/users", a.UserHandler.CreateUser)
		v1.GET("/users", a.UserHandler.ListUsers)
	})

	fmt.Println("Auth Service is running...")
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}
