package app

import (
	"context"
	"io"
	"net/http"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/payment-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/provider"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
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
	return &App{Base: baseApp, Cfg: cfg, DB: db, Service: service, OrderConsumers: orderConsumers, CompensationConsumers: compensationConsumers, Guards: guards}
}

func (a *App) Run() {
	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterReadiness(func() error { return a.DB.Ping(context.Background()) })
	a.Base.FinalizeRoutes()

	for _, consumer := range a.OrderConsumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleOrderCreated); err != nil {
				logger.GetLogger().Error("payment consumer stopped", zap.Error(err))
			}
		}()
	}
	for _, consumer := range a.CompensationConsumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleCompensationEvent); err != nil {
				logger.GetLogger().Error("payment compensation consumer stopped", zap.Error(err))
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
	v1.GET("/payments", func(c *gin.Context) {
		payments, err := a.Service.ListPayments(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payments": payments})
	})
	v1.GET("/payments/orders/:order_id", func(c *gin.Context) {
		payment, err := a.Service.GetPaymentByOrder(c.Request.Context(), c.Param("order_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"payment": payment})
	})

	r.POST("/v1/webhooks/stripe", a.stripeWebhook())
	r.POST("/v1/webhooks/vnpay", a.vnpayWebhook())
}

func (a *App) stripeWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		timestamp, signature := provider.StripeSignatureTimestamp(c.GetHeader("Stripe-Signature"))
		if err := a.Service.SimulateWebhook(c.Request.Context(), provider.ProviderStripe, body, timestamp, signature); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "accepted", "simulated": true})
	}
}

func (a *App) vnpayWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if err := a.Service.SimulateWebhook(c.Request.Context(), provider.ProviderVNPay, body, c.GetHeader("X-VNPAY-Timestamp"), c.GetHeader("X-VNPAY-Signature")); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "accepted", "simulated": true})
	}
}

func AutoMigrate(db *database.DB) error {
	return db.AutoMigrate(&domain.Payment{}, &domain.InboxEvent{}, &domain.WebhookEvent{})
}

var _ = events.TopicPaymentCompleted
