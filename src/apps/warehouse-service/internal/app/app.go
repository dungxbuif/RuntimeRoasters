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
	OrderWorker   *worker.OrderWorker
	Consumers     []kafka.Consumer
}

func NewApp(cfg *svcconfig.Config, db *database.DB, consumers []kafka.Consumer, harvestWorker *worker.HarvestWorker, orderWorker *worker.OrderWorker) *App {
	return &App{
		Cfg:           cfg,
		DB:            db,
		Consumers:     consumers,
		HarvestWorker: harvestWorker,
		OrderWorker:   orderWorker,
	}
}

func (a *App) Run(ctx context.Context) error {
	log := logger.GetLogger().With(zap.String("service", a.Cfg.AppName))
	log.Info("service started, waiting for harvest events", zap.String("topic", a.Cfg.KafkaHarvestTopic))
	go func() {
		if err := a.OrderWorker.Start(ctx); err != nil {
			log.Error("order worker stopped with error", zap.Error(err))
		}
	}()
	return a.HarvestWorker.Start(ctx)
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
