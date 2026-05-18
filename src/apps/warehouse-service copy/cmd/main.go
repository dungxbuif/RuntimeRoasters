package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/app"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		logger.GetLogger().Fatal("failed to initialize warehouse service", zap.Error(err))
	}
	defer cleanup()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := application.Run(ctx); err != nil {
			logger.GetLogger().Error("warehouse worker stopped with error", zap.Error(err))
		}
	}()

	<-sigChan
	cancel()
	if err := application.Shutdown(); err != nil {
		logger.GetLogger().Warn("failed to close warehouse consumer cleanly", zap.Error(err))
	}
}
