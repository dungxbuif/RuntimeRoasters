package app

import (
	"context"
	"net/http"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	Base      *base.App
	Cfg       *svcconfig.Config
	DB        *database.DB
	VDB       *redis.Client
	Service   *usecase.Service
	Consumers []kafka.Consumer
	Guards    *security.HTTPGuards
}

func NewApp(
	baseApp *base.App,
	cfg *svcconfig.Config,
	db *database.DB,
	vdb *redis.Client,
	service *usecase.Service,
	consumers []kafka.Consumer,
	guards *security.HTTPGuards,
) *App {
	return &App{
		Base:      baseApp,
		Cfg:       cfg,
		DB:        db,
		VDB:       vdb,
		Service:   service,
		Consumers: consumers,
		Guards:    guards,
	}
}

func (a *App) Run() {
	a.Base.RegisterHTTP(a.routes)
	a.Base.RegisterReadiness(func() error {
		if err := a.DB.Ping(context.Background()); err != nil {
			return err
		}
		return nil
	})
	a.Base.FinalizeRoutes()

	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			logger.GetLogger().Info("starting logistics consumer", zap.String("topic", c.Topic()))
			if err := c.Listen(context.Background(), a.Service.HandleWarehouseEvent); err != nil {
				logger.GetLogger().Error("logistics consumer stopped", zap.Error(err), zap.String("topic", c.Topic()))
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
	v1.GET("/shipments", func(c *gin.Context) {
		shipments, err := a.Service.ListShipments(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipments": shipments})
	})

	v1.POST("/shipments/:id/deliver", func(c *gin.Context) {
		if err := a.Service.DeliverShipment(c.Request.Context(), c.Param("id")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delivered"})
	})

	v1.POST("/drivers/location", func(c *gin.Context) {
		var req struct {
			DriverID   string  `json:"driver_id"`
			ShipmentID string  `json:"shipment_id"`
			Latitude   float64 `json:"latitude"`
			Longitude  float64 `json:"longitude"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if req.DriverID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "driver_id is required"})
			return
		}
		if err := a.Service.UpdateDriverLocation(c.Request.Context(), req.DriverID, req.ShipmentID, req.Latitude, req.Longitude); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})
}

func AutoMigrate(db *database.DB) error {
	return db.AutoMigrate(
		&domain.Driver{},
		&domain.Vehicle{},
		&domain.Shipment{},
		&domain.Location{},
		&domain.ProcessedKafkaMessage{},
		&domain.InboxEvent{},
	)
}
