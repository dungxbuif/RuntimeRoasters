package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"gorm.io/gorm"
)

type ProcessingUseCase struct {
	db *gorm.DB
}

func NewProcessingUseCase(db *gorm.DB) *ProcessingUseCase {
	return &ProcessingUseCase{db: db}
}

func (uc *ProcessingUseCase) StartProcessing(ctx context.Context, batchID string) error {
	return uc.db.Model(&domain.ProductionBatch{}).
		Where("batch_id = ?", batchID).
		Update("status", domain.BatchStatusProcessing).Error
}

func (uc *ProcessingUseCase) AddRoastRun(ctx context.Context, batchID string, input, output float64, roaster string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get Batch
		var batch domain.ProductionBatch
		if err := tx.Where("batch_id = ?", batchID).First(&batch).Error; err != nil {
			return err
		}

		if batch.Status != domain.BatchStatusProcessing {
			return fmt.Errorf("batch must be in PROCESSING status to add runs")
		}

		// 2. Count existing runs to get next number
		var count int64
		tx.Model(&domain.RoastRun{}).Where("batch_id = ?", batchID).Count(&count)

		// 3. Create Run
		loss := (1 - output/input) * 100
		run := domain.RoastRun{
			ID:           uuid.New().String(),
			BatchID:      batchID,
			RunNumber:    int(count) + 1,
			InputWeight:  input,
			OutputWeight: output,
			LossPercent:  loss,
			RoasterName:  roaster,
		}

		if err := tx.Create(&run).Error; err != nil {
			return err
		}

		// 4. Update Batch Aggregates
		newTotalInput := batch.TotalInputWeight + input
		newTotalOutput := batch.TotalOutputWeight + output
		newBatchLoss := (1 - newTotalOutput/newTotalInput) * 100

		qualityFlag := "NORMAL"
		if loss > 25.0 || loss < 10.0 {
			qualityFlag = "QUALITY_WARNING"
		}

		return tx.Model(&batch).Updates(map[string]interface{}{
			"total_input_weight":  newTotalInput,
			"total_output_weight": newTotalOutput,
			"weight_loss_percent": newBatchLoss,
			"quality_flag":        qualityFlag,
		}).Error
	})
}
