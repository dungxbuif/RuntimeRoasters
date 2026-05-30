package base

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"RuntimeRoasters/pkg/config"
	"RuntimeRoasters/pkg/errs"
	"RuntimeRoasters/pkg/logger"
	"RuntimeRoasters/pkg/telemetry"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Options struct {
	Name              string
	Config            config.BaseConfig
	GRPCServerOptions []grpc.ServerOption
}

type App struct {
	Name               string
	internalHost       string
	httpServer         *http.Server
	grpcServer         *grpc.Server
	logger             *zap.Logger
	ginEngine          *gin.Engine
	gwMux              *runtime.ServeMux
	shutdownFn         func()
	reflectionRegistry bool
}

func NewApp(opts Options) *App {
	logger.InitLogger(opts.Config.AppEnv, opts.Config.LogLevel)
	log := logger.GetLogger()
	log.Info("bootstrapping service", zap.String("name", opts.Name))

	otelShutdown, err := telemetry.InitTracer(opts.Name, opts.Config.OTLPEndpoint)
	if err != nil {
		log.Warn("failed to initialize telemetry", zap.Error(err))
	}

	if opts.Config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(otelgin.Middleware(opts.Name)) // HTTP/Gin tracing
	engine.Use(logger.GinLoggerMiddleware())  // Zap HTTP Logging
	engine.Use(errs.GinErrorHandler())

	engine.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	gwMux := runtime.NewServeMux()

	return &App{
		Name:         opts.Name,
		internalHost: defaultInternalHost(opts.Config.InternalHost),
		ginEngine:    engine,
		grpcServer: grpc.NewServer(
			append(opts.GRPCServerOptions,
				grpc.StatsHandler(otelgrpc.NewServerHandler()),
				grpc.ChainUnaryInterceptor(logger.GRPCLoggerInterceptor()),
			)...,
		),
		gwMux:      gwMux,
		logger:     log,
		shutdownFn: otelShutdown,
	}
}

func (a *App) RegisterHTTP(routes func(*gin.Engine)) {
	routes(a.ginEngine)
}

func (a *App) FinalizeRoutes() {
	// Route everything to gateway mux for /v1 ONLY IF NOT HANDLED BY GIN
	a.ginEngine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/v1") {
			c.Status(http.StatusOK)
			a.gwMux.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	})
}

func (a *App) RegisterGRPC(desc *grpc.ServiceDesc, impl interface{}) {
	a.grpcServer.RegisterService(desc, impl)
	if !a.reflectionRegistry {
		reflection.Register(a.grpcServer)
		a.reflectionRegistry = true
	}
}

func (a *App) RegisterGateway(register func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error, grpcPort int) {
	ctx := context.Background()
	// Gateway-to-gRPC tracing
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	endpoint := fmt.Sprintf("%s:%d", a.internalHost, grpcPort)
	if err := register(ctx, a.gwMux, endpoint, opts); err != nil {
		a.logger.Fatal("failed to register gateway", zap.Error(err))
	}
}

func defaultInternalHost(host string) string {
	if strings.TrimSpace(host) == "" {
		return "127.0.0.1"
	}
	return host
}

func (a *App) ServeSwagger(path, dir string) {
	a.logger.Info("Serving swagger definitions", zap.String("path", path), zap.String("dir", dir))
	a.ginEngine.StaticFS(path, http.Dir(dir))
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

	if a.shutdownFn != nil {
		a.shutdownFn()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		a.logger.Fatal("HTTP Server shutdown error", zap.Error(err))
	}

	a.grpcServer.GracefulStop()
	a.logger.Info("Shutdown complete")
}
