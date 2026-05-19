package app

import (
	"context"
	"net/http"
	"time"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/retail-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/retail-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/retail-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Base      *base.App
	Cfg       *svcconfig.Config
	DB        *database.DB
	Service   *usecase.Service
	Consumers []kafka.Consumer
	Guards    *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, svc *usecase.Service, consumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{Base: baseApp, Cfg: cfg, DB: db, Service: svc, Consumers: consumers, Guards: guards}
}

func (a *App) Run() {
	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterReadiness(func() error { return a.DB.Ping(context.Background()) })
	a.Base.FinalizeRoutes()

	go a.startOutboxRelay(context.Background())
	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleSagaEvent); err != nil {
				logger.GetLogger().Error("retail saga consumer stopped", zap.Error(err))
			}
		}()
	}

	a.Base.Run(a.Cfg.AppPort, a.Cfg.GRPCPort)
}

func (a *App) routes(r *gin.Engine) {
	v1 := r.Group("/v1")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	v1.GET("/stores", func(c *gin.Context) {
		stores, err := a.Service.ListStores(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"stores": stores})
	})

	v1.POST("/orders", func(c *gin.Context) {
		var req usecase.CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		order, err := a.Service.CreateOrder(c.Request.Context(), req, c.GetHeader("Idempotency-Key"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order": order})
	})

	v1.GET("/orders/:id", func(c *gin.Context) {
		order, err := a.Service.GetOrder(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"order": order})
	})
}

func (a *App) startOutboxRelay(ctx context.Context) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		if err := a.Service.ProcessOutbox(ctx, 20); err != nil {
			logger.GetLogger().Warn("retail outbox relay failed", zap.Error(err))
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func AutoMigrate(db *database.DB) error {
	return db.AutoMigrate(&domain.Store{}, &domain.Order{}, &domain.OutboxEvent{}, &domain.InboxEvent{})
}
