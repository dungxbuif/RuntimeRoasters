package main

import (
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

	application.Run()
	if err := application.Shutdown(); err != nil {
		logger.GetLogger().Warn("failed to close warehouse consumer cleanly", zap.Error(err))
	}
}
