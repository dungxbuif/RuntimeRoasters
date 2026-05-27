package app

import (
	"context"
	"net/http"
	"strconv"

	svcconfig "RuntimeRoasters/apps/trace-service/config"
	"RuntimeRoasters/apps/trace-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
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
	return &App{
		Base:      baseApp,
		Cfg:       cfg,
		DB:        db,
		Service:   service,
		Consumers: consumers,
		Guards:    guards,
	}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("trace-service starting", zap.Strings("topics", a.Cfg.TraceTopics))

	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleEvent); err != nil {
				log.Error("consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
			}
		}()
	}

	a.Base.RegisterHTTP(a.routes)
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
	public := r.Group("/v1/traces/public")
	public.GET("/topology/config", func(c *gin.Context) {
		cfg, err := a.Service.GetTopologyConfig(c.Request.Context(), true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cfg)
	})
	public.GET("/topology/history", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))
		history, err := a.Service.GetTopologyHistory(c.Request.Context(), usecase.TopologyHistoryFilter{
			FlowID:     c.Query("flow_id"),
			Cursor:     c.Query("cursor"),
			Limit:      limit,
			PublicOnly: true,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"events": history})
	})

	legacy := r.Group("/v1/trace")
	if a.Guards != nil {
		legacy.Use(a.Guards.Authn, a.Guards.Authz)
	}
	legacy.GET("/:id", func(c *gin.Context) {
		doc, found, err := a.Service.GetTraceDocument(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"message": "trace not found"})
			return
		}
		c.JSON(http.StatusOK, doc)
	})

	legacy.GET("/:id/document", func(c *gin.Context) {
		doc, found, err := a.Service.GetTraceDocument(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"message": "trace not found"})
			return
		}
		c.JSON(http.StatusOK, doc)
	})

	legacy.GET("/:id/events", func(c *gin.Context) {
		events, err := a.Service.GetTrace(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"events": events})
	})

	v1 := r.Group("/v1/traces")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	v1.GET("/topology/config", func(c *gin.Context) {
		cfg, err := a.Service.GetTopologyConfig(c.Request.Context(), false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cfg)
	})

	v1.GET("/topology/history", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.Query("limit"))
		history, err := a.Service.GetTopologyHistory(c.Request.Context(), usecase.TopologyHistoryFilter{
			FlowID:   c.Query("flow_id"),
			EntityID: c.Query("entity_id"),
			Cursor:   c.Query("cursor"),
			Limit:    limit,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"events": history})
	})
}
