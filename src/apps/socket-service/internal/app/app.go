package app

import (
	"context"
	"net/http"

	svcconfig "RuntimeRoasters/apps/socket-service/config"
	"RuntimeRoasters/apps/socket-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"RuntimeRoasters/pkg/topology"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	Base      *base.App
	Cfg       *svcconfig.Config
	RDB       *redis.Client
	Service   *usecase.Service
	Producer  kafka.Producer
	Consumers []kafka.Consumer
	Guards    *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, rdb *redis.Client, service *usecase.Service, producer kafka.Producer, consumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{Base: baseApp, Cfg: cfg, RDB: rdb, Service: service, Producer: producer, Consumers: consumers, Guards: guards}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.BaseConfig.AppName))
	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleKafkaMessage); err != nil {
				log.Error("socket consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
			}
		}()
	}
	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterReadiness(func() error { return a.RDB.Ping(context.Background()).Err() })
	a.Base.FinalizeRoutes()
	a.Base.Run(a.Cfg.BaseConfig.AppPort, a.Cfg.BaseConfig.GRPCPort)
}

func (a *App) Shutdown() error {
	a.Service.Close()
	var err error
	for _, consumer := range a.Consumers {
		if closeErr := consumer.Close(); closeErr != nil {
			err = closeErr
		}
	}
	_ = a.Producer.Close()
	_ = a.RDB.Close()
	if a.Guards != nil {
		a.Guards.Close()
	}
	return err
}

func (a *App) routes(r *gin.Engine) {
	r.GET("/v1/realtime/public/topology/ws", func(c *gin.Context) {
		_ = a.Service.ServeWS(c.Writer, c.Request, usecase.ClientScope{Public: true, FlowID: c.Query("flow_id")})
	})

	internal := r.Group("/internal/v1/socket")
	internal.POST("/events", func(c *gin.Context) {
		var req topology.BroadcastRequested
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if err := a.Service.PublishInternalEvent(c.Request.Context(), c.GetHeader("X-RR-API-Key"), req); err != nil {
			if err.Error() == "missing api key" || err.Error() == "invalid api key" {
				c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
	})

	private := r.Group("/v1/realtime")
	if a.Guards != nil {
		private.POST("/tickets", a.Guards.Authn, a.Guards.Authz, func(c *gin.Context) {
			claims, ok := identity.FromContext(c.Request.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
			ticket, err := a.Service.IssueTicket(c.Request.Context(), claims)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to issue ticket"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"ticket": ticket})
		})

		private.GET("/notifications", a.Guards.Authn, a.Guards.Authz, func(c *gin.Context) {
			claims, ok := identity.FromContext(c.Request.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
			notifications, err := a.Service.ListNotifications(c.Request.Context(), claims)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"notifications": notifications})
		})

		private.POST("/notifications/:id/ack", a.Guards.Authn, a.Guards.Authz, func(c *gin.Context) {
			claims, ok := identity.FromContext(c.Request.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
			notification, err := a.Service.AcknowledgeNotification(c.Request.Context(), claims, c.Param("id"))
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"notification": notification})
		})

		private.POST("/notifications/:id/resolve", a.Guards.Authn, a.Guards.Authz, func(c *gin.Context) {
			claims, ok := identity.FromContext(c.Request.Context())
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
				return
			}
			notification, err := a.Service.ResolveNotification(c.Request.Context(), claims, c.Param("id"))
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"notification": notification})
		})
	}

	private.GET("/stream", func(c *gin.Context) {
		claims, err := a.Service.ExchangeTicket(c.Request.Context(), c.Query("ticket"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}
		_ = a.Service.ServeWS(c.Writer, c.Request, usecase.ClientScope{
			Role:         claims.Role,
			Subject:      claims.Subject,
			StoreIDs:     claims.StoreIDs,
			WarehouseIDs: claims.WarehouseIDs,
			FlowID:       c.Query("flow_id"),
		})
	})
}
