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
	// 1. Config from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "54321")
	dbUser := getEnv("DB_USER", "user")
	dbPass := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "warehouse_db")
	
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
		dbHost, dbUser, dbPass, dbName, dbPort)
	
	kafkaBrokers := getEnv("KAFKA_BROKERS", "localhost:9094")
	brokers := []string{kafkaBrokers}
	groupID := getEnv("KAFKA_GROUP_ID", "warehouse-service-group")

	// 2. Database Connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// 3. Auto Migration
	fmt.Println("[WAREHOUSE] Running database migrations...")
	db.AutoMigrate(&domain.ProductionBatch{}, &domain.RoastRun{}, &domain.Inventory{}, &domain.InboxEvent{})

	// 4. Manual Dependency Injection
	// producer := kafka.NewProducer(brokers)
	intakeUseCase := usecase.NewIntakeUseCase(db)
	// processingUseCase := usecase.NewProcessingUseCase(db)
	// inventoryUseCase := usecase.NewInventoryUseCase(db, producer)
	
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
