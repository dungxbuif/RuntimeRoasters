package app

import (
	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/worker"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func InitializeApp() (*App, func(), error) {
	cfg := configDefaults(svcconfig.Load())

	logger.InitLogger(cfg.AppEnv, cfg.LogLevel)
	log := logger.GetLogger().With(zap.String("service", cfg.AppName))

	db, err := database.NewPostgres(database.PostgresConfig{
		URL:      cfg.DatabaseURL,
		LogLevel: cfg.DBLogLevel,
	})
	if err != nil {
		return nil, nil, err
	}

	log.Info("running database migrations")
	if err := db.AutoMigrate(&domain.ProductionBatch{}, &domain.RoastRun{}, &domain.Inventory{}, &domain.InboxEvent{}); err != nil {
		return nil, nil, err
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	intakeUseCase := usecase.NewIntakeUseCase(db.DB)
	orderReservationUC := usecase.NewOrderReservationUseCase(db.DB, producer, cfg.KafkaStockReservedTopic, cfg.KafkaStockFailedTopic)
	harvestConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-harvest", cfg.KafkaHarvestTopic)
	orderConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-orders", cfg.KafkaOrderTopic)
	harvestWorker := worker.NewHarvestWorker(harvestConsumer, intakeUseCase)
	orderWorker := worker.NewOrderWorker(orderConsumer, orderReservationUC)

	app := NewApp(&cfg, db, []kafka.Consumer{harvestConsumer, orderConsumer}, harvestWorker, orderWorker)

	cleanup := func() {
		_ = producer.Close()
		_ = harvestConsumer.Close()
		_ = orderConsumer.Close()
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
		cfg.KafkaHarvestTopic = "farm.harvest.events"
	}
	if cfg.KafkaStockTopic == "" {
		cfg.KafkaStockTopic = events.TopicWarehouseStockUpdated
	}
	if cfg.KafkaOrderTopic == "" {
		cfg.KafkaOrderTopic = events.TopicPaymentCompleted
	}
	if cfg.KafkaStockReservedTopic == "" {
		cfg.KafkaStockReservedTopic = events.TopicWarehouseStockReserved
	}
	if cfg.KafkaStockFailedTopic == "" {
		cfg.KafkaStockFailedTopic = events.TopicWarehouseStockReservationFailed
	}
	return cfg
}
