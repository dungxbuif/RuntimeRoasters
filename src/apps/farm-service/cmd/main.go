package main

import (
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/app"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		logger.GetLogger().Fatal("failed to initialize farm service", zap.Error(err))
	}
	defer cleanup()
	application.Run()
}
