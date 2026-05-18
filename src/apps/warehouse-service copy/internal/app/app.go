package app

import (
	"context"

	svcconfig "github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/worker"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

type App struct {
	Cfg           *svcconfig.Config
	DB            *database.DB
	HarvestWorker *worker.HarvestWorker
	Consumer      kafka.Consumer
}

func NewApp(cfg *svcconfig.Config, db *database.DB, consumer kafka.Consumer, harvestWorker *worker.HarvestWorker) *App {
	return &App{
		Cfg:           cfg,
		DB:            db,
		Consumer:      consumer,
		HarvestWorker: harvestWorker,
	}
}

func (a *App) Run(ctx context.Context) error {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("service started, waiting for harvest events", zap.String("topic", a.Cfg.KafkaHarvestTopic))
	return a.HarvestWorker.Start(ctx)
}

func (a *App) Shutdown() error {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("shutting down")
	return a.Consumer.Close()
}
