package app

import (
	"context"

	svcconfig "RuntimeRoasters/apps/warehouse-service/config"
	warehousegrpc "RuntimeRoasters/apps/warehouse-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"RuntimeRoasters/apps/warehouse-service/internal/worker"
	"RuntimeRoasters/pkg/base"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func InitializeApp() (*App, func(), error) {
	cfg := configDefaults(svcconfig.Load())

	db, err := database.NewPostgres(database.PostgresConfig{
		URL:      cfg.DatabaseURL,
		LogLevel: cfg.DBLogLevel,
	})
	if err != nil {
		return nil, nil, err
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	pickupUseCase := usecase.NewPickupUseCase(db.DB, producer, cfg.KafkaPickupTopic, cfg.KafkaNotificationTopic, cfg.KafkaIntakeCreatedTopic, cfg.DefaultWarehouseID)
	aggregationUC := usecase.NewAggregationUseCase(db.DB)
	processingUC := usecase.NewProcessingUseCase(db.DB)
	inventoryUC := usecase.NewInventoryUseCase(db.DB, producer, cfg.KafkaStockTopic)
	orderReservationUC := usecase.NewOrderReservationUseCase(db.DB, producer, cfg.KafkaStockReservedTopic, cfg.KafkaStockFailedTopic)
	harvestConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-harvest", cfg.KafkaHarvestTopic)
	pickupArrivedConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-pickup-arrived", cfg.KafkaPickupArrivedTopic)
	orderConsumers := make([]kafka.Consumer, 0, 2)
	for _, topic := range uniqueTopics(cfg.KafkaOrderTopic, events.TopicPaymentCompleted) {
		orderConsumers = append(orderConsumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-orders-"+topic, topic))
	}
	harvestWorker := worker.NewHarvestWorker(harvestConsumer, pickupUseCase)
	orderWorker := worker.NewOrderWorker(orderConsumers, orderReservationUC)

	consumers := append([]kafka.Consumer{harvestConsumer, pickupArrivedConsumer}, orderConsumers...)
	go func() {
		if err := pickupArrivedConsumer.Listen(context.Background(), pickupUseCase.HandlePickupArrived); err != nil {
			logger.GetLogger().Error("pickup arrived consumer stopped", zap.Error(err))
		}
	}()

	guards, err := security.NewHTTPGuards(security.HTTPOptions{
		ServiceName:      cfg.AppName,
		JWKSURL:          cfg.JWKSURL,
		InternalSecret:   cfg.InternalSecret,
		JWKSCacheTTL:     cfg.JWKSCacheTTL,
		ExpectedIssuer:   cfg.ExpectedIssuer,
		AuthServiceAddr:  cfg.AuthServiceAddr,
		KafkaBrokers:     cfg.KafkaBrokers,
		KafkaPolicyTopic: cfg.KafkaPolicyTopic,
	})
	if err != nil {
		return nil, nil, err
	}

	systemHandler := usecase.NewSystemHandler(db)

	baseApp := base.NewApp(base.Options{Name: cfg.AppName, Config: cfg.BaseConfig})
	app := NewApp(baseApp, &cfg, db, pickupUseCase, aggregationUC, processingUC, inventoryUC, systemHandler, consumers, harvestWorker, orderWorker, guards)

	cleanup := func() {
		guards.Close()
		_ = producer.Close()
		_ = harvestConsumer.Close()
		_ = pickupArrivedConsumer.Close()
		for _, consumer := range orderConsumers {
			_ = consumer.Close()
		}
	}

	return app, cleanup, nil
}

func configDefaults(cfg svcconfig.Config) svcconfig.Config {
	if cfg.AppName == "" {
		cfg.AppName = "warehouse-service"
	}
	if cfg.AppEnv == "" {
		cfg.AppEnv = "development"
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}
	if cfg.DBLogLevel == "" {
		cfg.DBLogLevel = "warn"
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "warehouse-service-group"
	}
	if cfg.KafkaHarvestTopic == "" {
		cfg.KafkaHarvestTopic = events.TopicFarmHarvestCreated
	}
	if cfg.KafkaStockTopic == "" {
		cfg.KafkaStockTopic = events.TopicWarehouseInventoryUpdated
	}
	if cfg.KafkaPickupTopic == "" {
		cfg.KafkaPickupTopic = events.TopicWarehousePickupRequested
	}
	if cfg.KafkaPickupArrivedTopic == "" {
		cfg.KafkaPickupArrivedTopic = events.TopicLogisticsPickupArrivedAtWarehouse
	}
	if cfg.KafkaIntakeCreatedTopic == "" {
		cfg.KafkaIntakeCreatedTopic = events.TopicWarehouseIntakeCreated
	}
	if cfg.KafkaNotificationTopic == "" {
		cfg.KafkaNotificationTopic = events.TopicNotificationCreated
	}
	if cfg.KafkaOrderTopic == "" {
		cfg.KafkaOrderTopic = events.TopicPaymentSimulatedCompleted
	}
	if cfg.KafkaStockReservedTopic == "" {
		cfg.KafkaStockReservedTopic = events.TopicWarehouseStockReserved
	}
	if cfg.KafkaStockFailedTopic == "" {
		cfg.KafkaStockFailedTopic = events.TopicWarehouseStockReservationFailed
	}
	if cfg.DefaultWarehouseID == "" {
		cfg.DefaultWarehouseID = "WAREHOUSE-HN-001"
	}
	return cfg
}

func uniqueTopics(topics ...string) []string {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(topics))
	for _, topic := range topics {
		if topic == "" {
			continue
		}
		if _, ok := seen[topic]; ok {
			continue
		}
		seen[topic] = struct{}{}
		unique = append(unique, topic)
	}
	return unique
}
