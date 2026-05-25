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
	registerLogisticsRoutes(v1, a)
	registerLogisticsRoutes(v1.Group("/logistics"), a)
}

func registerLogisticsRoutes(v1 *gin.RouterGroup, a *App) {
	v1.GET("/locations", func(c *gin.Context) {
		locations, err := a.Service.ListLocations(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	})
	v1.GET("/shipments", func(c *gin.Context) {
		shipments, err := a.Service.ListShipments(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipments": shipments})
	})
	v1.GET("/shipments/:id", func(c *gin.Context) {
		shipment, err := a.Service.GetShipment(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})

	v1.POST("/shipments/:id/deliver", func(c *gin.Context) {
		if err := a.Service.DeliverShipment(c.Request.Context(), c.Param("id")); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "delivered"})
	})
	v1.POST("/shipments/:id/assign", func(c *gin.Context) {
		var req struct {
			DriverID  string `json:"driver_id"`
			VehicleID string `json:"vehicle_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		shipment, err := a.Service.AssignShipment(c.Request.Context(), c.Param("id"), req.DriverID, req.VehicleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})
	v1.POST("/shipments/:id/depart", shipmentAction(a, func(ctx context.Context, id string) (*domain.Shipment, error) {
		return a.Service.DepartShipment(ctx, id)
	}))
	v1.POST("/shipments/:id/arrive", shipmentAction(a, func(ctx context.Context, id string) (*domain.Shipment, error) {
		return a.Service.ArriveShipment(ctx, id)
	}))
	v1.POST("/shipments/:id/confirm-load", shipmentAction(a, func(ctx context.Context, id string) (*domain.Shipment, error) {
		return a.Service.ConfirmLoad(ctx, id)
	}))
	v1.POST("/shipments/:id/confirm-delivery", shipmentAction(a, func(ctx context.Context, id string) (*domain.Shipment, error) {
		if err := a.Service.ConfirmDelivery(ctx, id); err != nil {
			return nil, err
		}
		return a.Service.GetShipment(ctx, id)
	}))
	v1.POST("/shipments/:id/return", shipmentAction(a, func(ctx context.Context, id string) (*domain.Shipment, error) {
		return a.Service.ReturnShipment(ctx, id)
	}))

	v1.POST("/drivers/location", func(c *gin.Context) {
		var req struct {
			DriverID   string  `json:"driver_id"`
			ShipmentID string  `json:"shipment_id"`
			Latitude   float64 `json:"latitude"`
			Longitude  float64 `json:"longitude"`
			Lat        float64 `json:"lat"`
			Lng        float64 `json:"lng"`
			Heading    float64 `json:"heading"`
			Speed      float64 `json:"speed"`
			RouteIndex int     `json:"route_index"`
			Status     string  `json:"status"`
			OccurredAt string  `json:"occurred_at"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		lat := req.Latitude
		lng := req.Longitude
		if lat == 0 {
			lat = req.Lat
		}
		if lng == 0 {
			lng = req.Lng
		}
		if err := a.Service.UpdateDriverLocation(c.Request.Context(), req.DriverID, req.ShipmentID, lat, lng); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	})
	v1.GET("/drivers", func(c *gin.Context) {
		drivers, err := a.Service.ListDrivers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"drivers": drivers})
	})
	v1.GET("/vehicles", func(c *gin.Context) {
		vehicles, err := a.Service.ListVehicles(c.Request.Context(), false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
	})
	v1.GET("/vehicles/available", func(c *gin.Context) {
		vehicles, err := a.Service.ListVehicles(c.Request.Context(), true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
	})
}

func shipmentAction(a *App, fn func(context.Context, string) (*domain.Shipment, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		shipment, err := fn(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	}
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
