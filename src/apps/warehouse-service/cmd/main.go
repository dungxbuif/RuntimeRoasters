package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/worker"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Config (Hardcoded for demo, usually from env)
	dsn := "host=localhost user=user password=password dbname=warehouse_db port=54321 sslmode=disable"
	brokers := []string{"localhost:9094"}
	groupID := "warehouse-service-group"

	// 2. Database Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 3. Auto Migration
	fmt.Println("[WAREHOUSE] Running database migrations...")
	db.AutoMigrate(&domain.ProductionBatch{}, &domain.InboxEvent{})

	// 4. Manual Dependency Injection
	intakeUseCase := usecase.NewIntakeUseCase(db)
	consumer := kafka.NewConsumer(brokers, groupID, "farm.harvest.events")
	harvestWorker := worker.NewHarvestWorker(consumer, intakeUseCase)

	// 5. Graceful Shutdown Setup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 6. Start Worker in background
	go func() {
		fmt.Println("[WAREHOUSE] Service started. Waiting for harvest events...")
		if err := harvestWorker.Start(ctx); err != nil {
			log.Printf("worker stopped with error: %v", err)
		}
	}()

	<-sigChan
	fmt.Println("[WAREHOUSE] Shutting down...")
	cancel()
	consumer.Close()
}
