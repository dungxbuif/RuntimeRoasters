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
	"github.com/dungxbuif/RuntimeRoasters/pkg/telemetry"
	"go.uber.org/zap"
)

func InitializeApp() (*App, func(), error) {
	cfg := configDefaults(svcconfig.Load())

	logger.InitLogger(cfg.AppEnv, cfg.LogLevel)
	log := logger.GetLogger().With(zap.String("service", cfg.AppName))
	otelShutdown, err := telemetry.InitTracer(cfg.AppName, cfg.OTLPEndpoint)
	if err != nil {
		log.Warn("failed to initialize tracer", zap.Error(err))
	}

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
	if err := ensureIntakesTable(db); err != nil {
		return nil, nil, err
	}

	producer := kafka.NewProducer(cfg.KafkaBrokers)
	intakeUseCase := usecase.NewIntakeUseCase(db.DB)
	orderReservationUC := usecase.NewOrderReservationUseCase(db.DB, producer, cfg.KafkaStockReservedTopic, cfg.KafkaStockFailedTopic)
	harvestConsumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-harvest", cfg.KafkaHarvestTopic)
	orderConsumers := make([]kafka.Consumer, 0, 2)
	for _, topic := range uniqueTopics(cfg.KafkaOrderTopic, events.TopicPaymentCompleted) {
		orderConsumers = append(orderConsumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-orders-"+topic, topic))
	}
	harvestWorker := worker.NewHarvestWorker(harvestConsumer, intakeUseCase)
	orderWorker := worker.NewOrderWorker(orderConsumers, orderReservationUC)

	consumers := append([]kafka.Consumer{harvestConsumer}, orderConsumers...)
	app := NewApp(&cfg, db, consumers, harvestWorker, orderWorker)

	cleanup := func() {
		if otelShutdown != nil {
			otelShutdown()
		}
		_ = producer.Close()
		_ = harvestConsumer.Close()
		for _, consumer := range orderConsumers {
			_ = consumer.Close()
		}
	}

	return app, cleanup, nil
}

func ensureIntakesTable(db *database.DB) error {
	return db.Exec(`
		CREATE TABLE IF NOT EXISTS intakes (
			id varchar(64) PRIMARY KEY,
			harvest_id varchar(64) NOT NULL,
			coffee_type varchar(20) NOT NULL,
			origin_code varchar(10) NOT NULL,
			quantity decimal(10,2) NOT NULL,
			status varchar(20) NOT NULL DEFAULT 'UNASSIGNED',
			batch_id varchar(80),
			created_at timestamptz,
			updated_at timestamptz
		);
		CREATE INDEX IF NOT EXISTS idx_intakes_harvest_id ON intakes (harvest_id);
		CREATE INDEX IF NOT EXISTS idx_intakes_batch_id ON intakes (batch_id);
	`).Error
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
	if cfg.KafkaOrderTopic == "" {
		cfg.KafkaOrderTopic = events.TopicPaymentSimulatedCompleted
	}
	if cfg.KafkaStockReservedTopic == "" {
		cfg.KafkaStockReservedTopic = events.TopicWarehouseStockReserved
	}
	if cfg.KafkaStockFailedTopic == "" {
		cfg.KafkaStockFailedTopic = events.TopicWarehouseStockReservationFailed
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
