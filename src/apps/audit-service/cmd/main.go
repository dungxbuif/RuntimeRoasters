package main

import (
	"RuntimeRoasters/apps/audit-service/internal/app"
	"RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	application, cleanup, err := app.InitializeApp()
	if err != nil {
		logger.GetLogger().Fatal("failed to initialize audit service", zap.Error(err))
	}
	defer cleanup()
	application.Run()
}
