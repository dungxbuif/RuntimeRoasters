package app

import (
	"context"
	"time"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/trace-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/trace-service/internal/search"
	"github.com/dungxbuif/RuntimeRoasters/apps/trace-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/security"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func InitializeApp() (*App, func(), error) {
	cfg := defaults(svcconfig.Load())
	logger.InitLogger(cfg.AppEnv, cfg.LogLevel)
	db, err := database.NewPostgres(database.PostgresConfig{URL: cfg.DatabaseURL, LogLevel: cfg.DBLogLevel})
	if err != nil {
		return nil, nil, err
	}
	if err := AutoMigrate(db); err != nil {
		return nil, nil, err
	}
	searchClient := initSearchClient(cfg)
	service := usecase.NewService(db.DB, searchClient)
	consumers := make([]kafka.Consumer, 0, len(cfg.TraceTopics))
	for _, topic := range cfg.TraceTopics {
		consumers = append(consumers, kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaGroupID+"-"+topic, topic))
	}
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
	baseApp := base.NewApp(base.Options{Name: cfg.AppName, Config: cfg.BaseConfig})
	app := NewApp(baseApp, &cfg, db, service, consumers, guards)
	cleanup := func() {
		guards.Close()
		for _, c := range consumers {
			_ = c.Close()
		}
		if sqlDB, err := db.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	return app, cleanup, nil
}

func defaults(cfg svcconfig.Config) svcconfig.Config {
	if cfg.AppName == "" {
		cfg.AppName = "trace-service"
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
	if cfg.AppPort == 0 {
		cfg.AppPort = 8087
	}
	if cfg.GRPCPort == 0 {
		cfg.GRPCPort = 50057
	}
	if cfg.KafkaGroupID == "" {
		cfg.KafkaGroupID = "trace-service-group"
	}
	if cfg.ElasticsearchURL == "" {
		cfg.ElasticsearchURL = "http://localhost:9200"
	}
	if cfg.ElasticsearchIndex == "" {
		cfg.ElasticsearchIndex = "coffee_traceability"
	}
	if len(cfg.TraceTopics) == 0 {
		cfg.TraceTopics = events.TraceableTopics
	}
	return cfg
}

func initSearchClient(cfg svcconfig.Config) *search.ElasticsearchClient {
	if !cfg.ElasticsearchEnabled {
		return nil
	}
	client := search.NewElasticsearchClient(cfg.ElasticsearchURL, cfg.ElasticsearchIndex)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.EnsureIndex(ctx); err != nil {
		logger.GetLogger().Warn("elasticsearch disabled because index initialization failed", zap.Error(err))
		return nil
	}
	logger.GetLogger().Info("elasticsearch read model initialized", zap.String("index", cfg.ElasticsearchIndex), zap.String("url", cfg.ElasticsearchURL))
	return client
}
