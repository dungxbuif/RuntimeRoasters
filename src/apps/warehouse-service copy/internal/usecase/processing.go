package usecase

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"gorm.io/gorm"
)

type ProcessingUseCase struct {
	db *gorm.DB
}

func NewProcessingUseCase(db *gorm.DB) *ProcessingUseCase {
	return &ProcessingUseCase{db: db}
}

// SimulateProcessing thực hiện đổi trạng thái và giả lập kết quả rang xay
func (uc *ProcessingUseCase) SimulateProcessing(ctx context.Context, batchID string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lấy thông tin Batch
		var batch domain.ProductionBatch
		if err := tx.Where("batch_id = ?", batchID).First(&batch).Error; err != nil {
			return err
		}

		// Chỉ cho phép xử lý nếu đang ở trạng thái RECEIVED
		if batch.Status != domain.BatchStatusReceived {
			return fmt.Errorf("batch status must be RECEIVED to start processing, current: %s", batch.Status)
		}

		// 2. Giả lập kết quả (Microservice Technique: Simulation)
		// Coffee roasting typically loses 12-20% weight
		lossPercent := 12.0 + rand.Float64()*8.0
		outputWeight := batch.IntakeWeight * (1 - lossPercent/100)

		qualityFlag := "NORMAL"
		if lossPercent > 18.5 {
			qualityFlag = "DARK_ROAST_WARNING"
		}

		// 3. Update trạng thái và kết quả đồng thời
		return tx.Model(&batch).Updates(map[string]interface{}{
			"status":              domain.BatchStatusProcessing,
			"output_weight":       outputWeight,
			"weight_loss_percent": lossPercent,
			"quality_flag":        qualityFlag,
		}).Error
	})
}
