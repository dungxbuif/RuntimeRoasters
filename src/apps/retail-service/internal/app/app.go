package app

import (
	"context"
	"net/http"
	"time"

	svcconfig "RuntimeRoasters/apps/retail-service/config"
	retailgrpc "RuntimeRoasters/apps/retail-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/retail-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Base          *base.App
	Cfg           *svcconfig.Config
	DB            *database.DB
	Service       *usecase.Service
	SystemHandler *retailgrpc.SystemHandler
	Consumers     []kafka.Consumer
	Guards        *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, service *usecase.Service, systemHandler *retailgrpc.SystemHandler, consumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{
		Base:          baseApp,
		Cfg:           cfg,
		DB:            db,
		Service:       service,
		SystemHandler: systemHandler,
		Consumers:     consumers,
		Guards:        guards,
	}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("retail-service starting")

	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleSagaEvent); err != nil {
				log.Error("consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
			}
		}()
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for range ticker.C {
			if err := a.Service.ProcessOutbox(context.Background(), 10); err != nil {
				log.Error("failed to process outbox", zap.Error(err))
			}
		}
	}()

	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterGRPC(&systemv1.SystemService_ServiceDesc, a.SystemHandler)
	a.Base.RegisterReadiness(func() error { return a.DB.Ping(context.Background()) })
	a.Base.FinalizeRoutes()
	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}

func (a *App) Shutdown() error {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("shutting down")
	var err error
	for _, consumer := range a.Consumers {
		if closeErr := consumer.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

func (a *App) routes(r *gin.Engine) {
	v1 := r.Group("/v1/retail")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	// Direct system routes in Gin to bypass gRPC Gateway for utility routes
	sys := r.Group("/v1/system")
	if a.Guards != nil {
		sys.Use(a.Guards.Authn, a.Guards.Authz)
	}
	sys.GET("/status", func(c *gin.Context) {
		res, err := a.SystemHandler.GetStatus(c.Request.Context(), &systemv1.GetStatusRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})
	sys.POST("/seed", func(c *gin.Context) {
		var req systemv1.SeedDataRequest
		_ = c.ShouldBindJSON(&req)
		res, err := a.SystemHandler.SeedData(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	registerBusinessRoutes := func(group *gin.RouterGroup) {
		group.GET("/stores", func(c *gin.Context) {
			stores, err := a.Service.ListStores(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"stores": stores})
		})

		group.POST("/orders", func(c *gin.Context) {
			var req usecase.CreateOrderRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			order, err := a.Service.CreateOrder(c.Request.Context(), req, c.GetHeader("X-Idempotency-Key"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"order": order})
		})

		group.GET("/orders/:id", func(c *gin.Context) {
			order, err := a.Service.GetOrder(c.Request.Context(), c.Param("id"))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"order": order})
		})
	}

	registerBusinessRoutes(v1)

	legacy := r.Group("/v1")
	if a.Guards != nil {
		legacy.Use(a.Guards.Authn, a.Guards.Authz)
	}
	registerBusinessRoutes(legacy)
}
