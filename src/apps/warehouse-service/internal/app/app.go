package app

import (
	"context"
	"net/http"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/worker"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Base          *base.App
	Cfg           *svcconfig.Config
	DB            *database.DB
	PickupService *usecase.PickupUseCase
	Aggregation   *usecase.AggregationUseCase
	Processing    *usecase.ProcessingUseCase
	Inventory     *usecase.InventoryUseCase
	HarvestWorker *worker.HarvestWorker
	OrderWorker   *worker.OrderWorker
	Consumers     []kafka.Consumer
	Guards        *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, pickupSvc *usecase.PickupUseCase, aggregation *usecase.AggregationUseCase, processing *usecase.ProcessingUseCase, inventory *usecase.InventoryUseCase, consumers []kafka.Consumer, harvestWorker *worker.HarvestWorker, orderWorker *worker.OrderWorker, guards *security.HTTPGuards) *App {
	return &App{
		Base:          baseApp,
		Cfg:           cfg,
		DB:            db,
		PickupService: pickupSvc,
		Aggregation:   aggregation,
		Processing:    processing,
		Inventory:     inventory,
		Consumers:     consumers,
		HarvestWorker: harvestWorker,
		OrderWorker:   orderWorker,
		Guards:        guards,
	}
}

func (a *App) Run() {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("service started, waiting for harvest events", zap.String("topic", a.Cfg.KafkaHarvestTopic))
	go func() {
		if err := a.OrderWorker.Start(context.Background()); err != nil {
			log.Error("order worker stopped with error", zap.Error(err))
		}
	}()
	go func() {
		if err := a.HarvestWorker.Start(context.Background()); err != nil {
			log.Error("harvest worker stopped with error", zap.Error(err))
		}
	}()
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
	v1 := r.Group("/v1/warehouse")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	v1.GET("/intakes", func(c *gin.Context) {
		intakes, err := a.Aggregation.GetUnassignedIntakes(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"intakes": intakes})
	})

	v1.GET("/pickup-requests", func(c *gin.Context) {
		pickups, err := a.PickupService.ListPickupRequests(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pickup_requests": pickups})
	})
	v1.POST("/pickup-requests/:id/dispatch", func(c *gin.Context) {
		pickup, err := a.PickupService.DispatchPickup(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pickup_request": pickup})
	})
	v1.POST("/pickup-requests/:id/receive", func(c *gin.Context) {
		intake, err := a.PickupService.ReceivePickup(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"intake": intake})
	})

	v1.GET("/batches", func(c *gin.Context) {
		var batches []domain.ProductionBatch
		query := a.DB.WithContext(c.Request.Context()).Preload("Intakes").Preload("RoastRuns").Order("created_at DESC")
		_, warehouseIDs, allWarehouses := identity.WarehouseScopeFromContext(c.Request.Context())
		if !allWarehouses {
			if len(warehouseIDs) == 0 {
				c.JSON(http.StatusOK, gin.H{"batches": []domain.ProductionBatch{}})
				return
			}
			query = query.Where("warehouse_id IN ?", warehouseIDs)
		}
		if err := query.Find(&batches).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"batches": batches})
	})
	v1.GET("/batches/:id", func(c *gin.Context) {
		var batch domain.ProductionBatch
		if err := a.DB.WithContext(c.Request.Context()).Preload("Intakes").Preload("RoastRuns").Where("id = ? OR batch_id = ?", c.Param("id"), c.Param("id")).Take(&batch).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(c.Request.Context()); !allWarehouses && !claims.CanAccessWarehouse(batch.WarehouseID) {
			c.JSON(http.StatusNotFound, gin.H{"message": "batch not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"batch": batch})
	})
	v1.POST("/batches", func(c *gin.Context) {
		var req struct {
			IntakeIDs []string `json:"intake_ids"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		batchID, err := a.Aggregation.CreateProductionBatch(c.Request.Context(), req.IntakeIDs)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"batch_id": batchID})
	})
	v1.PATCH("/batches/:id/intake", func(c *gin.Context) {
		var batch domain.ProductionBatch
		if err := a.DB.WithContext(c.Request.Context()).Preload("Intakes").Where("id = ? OR batch_id = ?", c.Param("id"), c.Param("id")).Take(&batch).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(c.Request.Context()); !allWarehouses && !claims.CanAccessWarehouse(batch.WarehouseID) {
			c.JSON(http.StatusNotFound, gin.H{"message": "batch not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"batch": batch})
	})
	v1.POST("/batches/:id/process", func(c *gin.Context) {
		var batch domain.ProductionBatch
		if err := a.DB.WithContext(c.Request.Context()).Where("id = ? OR batch_id = ?", c.Param("id"), c.Param("id")).Take(&batch).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(c.Request.Context()); !allWarehouses && !claims.CanAccessWarehouse(batch.WarehouseID) {
			c.JSON(http.StatusNotFound, gin.H{"message": "batch not found"})
			return
		}
		if err := a.Processing.StartSimulation(c.Request.Context(), batch.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "PROCESSING"})
	})
	v1.POST("/batches/:id/finalize", func(c *gin.Context) {
		var batch domain.ProductionBatch
		if err := a.DB.WithContext(c.Request.Context()).Where("id = ? OR batch_id = ?", c.Param("id"), c.Param("id")).Take(&batch).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		if err := a.Inventory.FinalizeBatch(c.Request.Context(), batch.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "STOCKED"})
	})
	v1.GET("/inventory", func(c *gin.Context) {
		var inventory []domain.Inventory
		query := a.DB.WithContext(c.Request.Context()).Order("updated_at DESC")
		_, warehouseIDs, allWarehouses := identity.WarehouseScopeFromContext(c.Request.Context())
		if !allWarehouses {
			if len(warehouseIDs) == 0 {
				c.JSON(http.StatusOK, gin.H{"inventory": []domain.Inventory{}})
				return
			}
			query = query.Where("warehouse_id IN ?", warehouseIDs)
		}
		if err := query.Find(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"inventory": inventory})
	})
	v1.GET("/dispatch-requests", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"dispatch_requests": []interface{}{}})
	})
	v1.POST("/dispatch-requests/:id/dispatch", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"id": c.Param("id"), "status": "ACCEPTED"})
	})
}
