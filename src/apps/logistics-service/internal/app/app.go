package app

import (
	"context"
	"net/http"
	"strconv"

	"RuntimeRoasters/apps/logistics-service/config"
	logisticsgrpc "RuntimeRoasters/apps/logistics-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/logistics-service/internal/usecase"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type App struct {
	Base          *base.App
	Cfg           *config.Config
	DB            *database.DB
	VDB           *redis.Client
	Service       *usecase.Service
	SystemHandler *logisticsgrpc.SystemHandler
	Consumers     []kafka.Consumer
	Guards        *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *config.Config, db *database.DB, vdb *redis.Client, service *usecase.Service, systemHandler *logisticsgrpc.SystemHandler, consumers []kafka.Consumer, guards *security.HTTPGuards) *App {
	return &App{
		Base:          baseApp,
		Cfg:           cfg,
		DB:            db,
		VDB:           vdb,
		Service:       service,
		SystemHandler: systemHandler,
		Consumers:     consumers,
		Guards:        guards,
	}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("logistics-service starting")

	for _, consumer := range a.Consumers {
		c := consumer
		go func() {
			if err := c.Listen(context.Background(), a.Service.HandleWarehouseEvent); err != nil {
				log.Error("consumer stopped", zap.String("topic", c.Topic()), zap.Error(err))
			}
		}()
	}

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
	v1 := r.Group("/v1/logistics")
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

	v1.POST("/shipments/:id/depart", func(c *gin.Context) {
		shipment, err := a.Service.DepartShipment(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})

	v1.POST("/shipments/:id/arrive", func(c *gin.Context) {
		shipment, err := a.Service.ArriveShipment(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})

	v1.POST("/shipments/:id/confirm-load", func(c *gin.Context) {
		shipment, err := a.Service.ConfirmLoad(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})

	v1.POST("/shipments/:id/confirm-delivery", func(c *gin.Context) {
		err := a.Service.ConfirmDelivery(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "delivered"})
	})

	v1.POST("/shipments/:id/return", func(c *gin.Context) {
		shipment, err := a.Service.ReturnShipment(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shipment": shipment})
	})

	v1.POST("/gps", func(c *gin.Context) {
		a.handleGPSUpdate(c)
	})

	v1.POST("/drivers/location", func(c *gin.Context) {
		a.handleGPSUpdate(c)
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
		available, _ := strconv.ParseBool(c.Query("available"))
		vehicles, err := a.Service.ListVehicles(c.Request.Context(), available)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"vehicles": vehicles})
	})

	v1.GET("/locations", func(c *gin.Context) {
		locations, err := a.Service.ListLocations(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	})
}

func (a *App) handleGPSUpdate(c *gin.Context) {
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
	if err := a.Service.UpdateDriverLocation(c.Request.Context(), req.DriverID, req.ShipmentID, req.Latitude, req.Longitude); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
