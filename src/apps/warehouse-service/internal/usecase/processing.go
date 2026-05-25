package usecase

import (
	"context"
	"math/rand"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProcessingUseCase struct {
	db *gorm.DB
}

func NewProcessingUseCase(db *gorm.DB) *ProcessingUseCase {
	return &ProcessingUseCase{db: db}
}

func (uc *ProcessingUseCase) StartSimulation(ctx context.Context, internalID string) error {
	// 1. Set to PROCESSING
	err := uc.db.Model(&domain.ProductionBatch{}).
		Where("id = ? AND status = ?", internalID, domain.BatchStatusDraft).
		Update("status", domain.BatchStatusProcessing).Error
	if err != nil {
		return err
	}

	// 2. Launch Async Simulation
	go uc.simulateRoasting(internalID)

	return nil
}

func (uc *ProcessingUseCase) simulateRoasting(internalID string) {
	ctx := context.Background()
	log := logger.FromContext(ctx).With(zap.String("batch_internal_id", internalID))
	log.Info("Starting async roasting simulation...")

	var batch domain.ProductionBatch
	if err := uc.db.Preload("Intakes").Where("id = ?", internalID).First(&batch).Error; err != nil {
		log.Error("failed to find batch for simulation", zap.Error(err))
		return
	}

	// Simulate 3 Runs
	numRuns := 3
	runWeight := batch.TotalInputWeight / float64(numRuns)
	var totalOutput float64

	for i := 1; i <= numRuns; i++ {
		time.Sleep(2 * time.Second) // Simulated roast time

		loss := 12.0 + rand.Float64()*6.0 // 12-18% loss
		output := runWeight * (1 - loss/100)
		totalOutput += output

		run := domain.RoastRun{
			BatchID:      batch.ID,
			RunNumber:    i,
			InputWeight:  runWeight,
			OutputWeight: output,
			Status:       "COMPLETED",
		}
		uc.db.Create(&run)
		log.Info("Roast run completed", zap.Int("run", i), zap.Float64("loss", loss))
	}

	// Finalize Simulation Results
	finalLoss := (1 - totalOutput/batch.TotalInputWeight) * 100
	uc.db.Model(&batch).Updates(map[string]interface{}{
		"status":              domain.BatchStatusReady,
		"total_output_weight": totalOutput,
		"weight_loss_percent": finalLoss,
	})

	log.Info("Simulation finished. Batch ready to stock.", zap.Float64("total_output", totalOutput))
}
