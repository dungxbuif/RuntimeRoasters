package app

import (
	"context"
	"net/http"

	svcconfig "RuntimeRoasters/apps/warehouse-service/config"
	warehousegrpc "RuntimeRoasters/apps/warehouse-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"RuntimeRoasters/apps/warehouse-service/internal/worker"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Base          *base.App
	Cfg           *svcconfig.Config
	DB            *database.DB
	PickupService *usecase.PickupUseCase
	WarehouseUC   *usecase.WarehouseUseCase
	Aggregation   *usecase.AggregationUseCase
	Processing    *usecase.ProcessingUseCase
	Inventory     *usecase.InventoryUseCase
	SystemHandler *warehousegrpc.SystemHandler
	HarvestWorker *worker.HarvestWorker
	OrderWorker   *worker.OrderWorker
	Consumers     []kafka.Consumer
	Guards        *security.HTTPGuards
}

func NewApp(baseApp *base.App, cfg *svcconfig.Config, db *database.DB, pickupSvc *usecase.PickupUseCase, warehouseUC *usecase.WarehouseUseCase, aggregation *usecase.AggregationUseCase, processing *usecase.ProcessingUseCase, inventory *usecase.InventoryUseCase, systemHandler *warehousegrpc.SystemHandler, consumers []kafka.Consumer, harvestWorker *worker.HarvestWorker, orderWorker *worker.OrderWorker, guards *security.HTTPGuards) *App {
	return &App{
		Base:          baseApp,
		Cfg:           cfg,
		DB:            db,
		PickupService: pickupSvc,
		WarehouseUC:   warehouseUC,
		Aggregation:   aggregation,
		Processing:    processing,
		Inventory:     inventory,
		SystemHandler: systemHandler,
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
	v1 := r.Group("/v1/warehouse")
	if a.Guards != nil {
		v1.Use(a.Guards.Authn, a.Guards.Authz)
	}

	sys := r.Group("/v1/system")
	sys.Use(func(c *gin.Context) {
		secret := c.GetHeader("X-Internal-Secret")
		if secret == "" || secret != a.Cfg.InternalSecret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or missing internal secret"})
			return
		}
		c.Next()
	})
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

	v1.GET("/warehouses", func(c *gin.Context) {
		warehouses, err := a.WarehouseUC.ListWarehouses(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
	})
	v1.POST("/warehouses", func(c *gin.Context) {
		var wh domain.Warehouse
		if err := c.ShouldBindJSON(&wh); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		created, err := a.WarehouseUC.CreateWarehouse(c.Request.Context(), wh)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"warehouse": created})
	})

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
		var req struct {
			DriverID  string `json:"driver_id"`
			VehicleID string `json:"vehicle_id"`
		}
		_ = c.ShouldBindJSON(&req)
		res, err := a.PickupService.DispatchPickup(c.Request.Context(), c.Param("id"), req.DriverID, req.VehicleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pickup_request": res})
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
