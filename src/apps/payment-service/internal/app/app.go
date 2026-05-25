package app

import (
	"context"
	"io"
	"net/http"

	svcconfig "RuntimeRoasters/apps/payment-service/config"
	"RuntimeRoasters/apps/payment-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Base                  *base.App
	Cfg                   *svcconfig.Config
	DB                    *database.DB
	Service               *usecase.Service
	OrderConsumers        []kafka.Consumer
	CompensationConsumers []kafka.Consumer
	Guards                *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, service *usecase.Service, orderConsumers []kafka.Consumer, compensationConsumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{
		Base:                  baseApp,
		Cfg:                   cfg,
		DB:                    db,
		Service:               service,
		OrderConsumers:        orderConsumers,
		CompensationConsumers: compensationConsumers,
		Guards:                guards,
	}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("payment-service starting")

	for _, consumer := range a.OrderConsumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleOrderCreated); err != nil {
				log.Error("order consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
			}
		}()
	}

	for _, consumer := range a.CompensationConsumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleCompensationEvent); err != nil {
				log.Error("compensation consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
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
	for _, consumer := range a.OrderConsumers {
		if closeErr := consumer.Close(); closeErr != nil {
			err = closeErr
		}
	}
	for _, consumer := range a.CompensationConsumers {
		if closeErr := consumer.Close(); closeErr != nil {
			err = closeErr
		}
	}
	return err
}

func (a *App) routes(r *gin.Engine) {
	v1 := r.Group("/v1/payments")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	v1.GET("", func(c *gin.Context) {
		payments, err := a.Service.ListPayments(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payments": payments})
	})

	v1.GET("/order/:id", func(c *gin.Context) {
		payment, err := a.Service.GetPaymentByOrder(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payment": payment})
	})

	v1.POST("/webhook/:provider", func(c *gin.Context) {
		payload, _ := io.ReadAll(c.Request.Body)
		err := a.Service.SimulateWebhook(
			c.Request.Context(),
			c.Param("provider"),
			payload,
			c.GetHeader("X-Timestamp"),
			c.GetHeader("X-Signature"),
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "processed"})
	})
}
