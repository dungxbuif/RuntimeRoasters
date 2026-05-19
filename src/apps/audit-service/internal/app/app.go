package app

import (
	"context"
	"net/http"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/audit-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/audit-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/audit-service/internal/usecase"
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

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, service *usecase.Service, consumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{Base: baseApp, Cfg: cfg, DB: db, Service: service, Consumers: consumers, Guards: guards}
}

func (a *App) Run() {
	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterReadiness(func() error { return a.DB.Ping(context.Background()) })
	a.Base.FinalizeRoutes()
	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleEvent); err != nil {
				logger.GetLogger().Error("audit consumer stopped", zap.Error(err))
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
	v1.GET("/audit/:partition", func(c *gin.Context) {
		logs, err := a.Service.ListByPartition(c.Request.Context(), c.Param("partition"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": logs})
	})
}

func AutoMigrate(db *database.DB) error {
	return db.AutoMigrate(&domain.AuditLog{})
}
