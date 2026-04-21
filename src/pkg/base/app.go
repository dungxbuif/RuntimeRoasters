package base

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/config"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Options struct {
	Name   string
	Config config.BaseConfig
}

type App struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	logger     *zap.Logger
	ginEngine  *gin.Engine
}

func NewApp(opts Options) *App {
	logger.InitLogger(opts.Config.AppEnv, opts.Config.LogLevel)
	log := logger.GetLogger()
	log.Info("bootstrapping service", zap.String("name", opts.Name))

	if opts.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(errs.GinErrorHandler())

	engine.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return &App{
		ginEngine:  engine,
		grpcServer: grpc.NewServer(),
		logger:     log,
	}
}

func (a *App) RegisterHTTP(routes func(*gin.Engine)) {
	routes(a.ginEngine)
}

func (a *App) RegisterGRPC(desc *grpc.ServiceDesc, impl interface{}) {
	a.grpcServer.RegisterService(desc, impl)
}

func (a *App) RegisterReadiness(check func() error) {
	a.ginEngine.GET("/health/ready", func(c *gin.Context) {
		if err := check(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded", "reason": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func (a *App) Run(httpPort, grpcPort int) {
	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: a.ginEngine,
	}

	go func() {
		a.logger.Info("Starting HTTP server", zap.Int("port", httpPort))
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
		if err != nil {
			a.logger.Fatal("gRPC listen error", zap.Error(err))
		}
		a.logger.Info("Starting gRPC server", zap.Int("port", grpcPort))
		if err := a.grpcServer.Serve(lis); err != nil {
			a.logger.Fatal("gRPC serve error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.logger.Info("Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Fatal("HTTP Server shutdown error", zap.Error(err))
	}

	a.grpcServer.GracefulStop()

	a.logger.Info("Shutdown complete")
}
