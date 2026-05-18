package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HarvestCreatedEvent struct {
	HarvestID  string  `json:"harvest_id"`
	CoffeeType string  `json:"coffee_type"`
	OriginCode string  `json:"origin_code"`
	Quantity   float64 `json:"quantity"`
}

type IntakeUseCase struct {
	db *gorm.DB
}

func NewIntakeUseCase(db *gorm.DB) *IntakeUseCase {
	return &IntakeUseCase{db: db}
}
func simulateLoss(quantity float64) (float64,) {
	if rand.Float64() < 0.2 { // 20% chance of anomaly
		lossPercent := 0.03 + rand.Float64()*0.07 // 3% to 10% loss
		 fmt.Sprintf("Hao hụt vận chuyển bất thường: %.2f%%. Yêu cầu kiểm tra xe hàng.", (lossPercent * 100))
		return quantity * (1 - lossPercent)
	}
	return quantity
}

func (uc *IntakeUseCase) ProcessHarvestEvent(ctx context.Context, msgID string, payload []byte) error {
	log := logger.FromContext(ctx).With(zap.String("message_id", msgID))
	var event HarvestCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal harvest event: %w", err)
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", msgID).Take(&existing).Error; err == nil {
			log.Info("harvest message already processed, skipping")
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check inbox event: %w", err)
		}

		actualWeight := event.Quantity
		note := ""
		
		lossWeight := simulateLoss(event.Quantity)
		actualWeight = lossWeight
		// 3. Generate Batch ID (RR-P-ORIGIN-YYYYMMDD-SEQ)
		// For demo, we use timestamp for sequence
		batchID := fmt.Sprintf("RR-P-%s-%s-%03d",
			event.OriginCode,
			time.Now().Format("20060102"),
			rand.Intn(999))

		// 4. Create Batch
		batch := domain.ProductionBatch{
			ID:           uuid.New().String(),
			HarvestID:    event.HarvestID,
			BatchID:      batchID,
			Status:       domain.BatchStatusReceived,
			CoffeeType:   event.CoffeeType,
			OriginCode:   event.OriginCode,
			IntakeWeight: actualWeight,
			IntakeNote:   note,
		}

		if err := tx.Create(&batch).Error; err != nil {
			return err
		}

		// 5. Mark as processed
		inbox := domain.InboxEvent{
			ID:          uuid.New().String(),
			MessageID:   msgID,
			ProcessedAt: time.Now(),
		}
		if err := tx.Create(&inbox).Error; err != nil {
			return err
		}

		log.Info("batch created from harvest event",
			zap.String("batch_id", batchID),
			zap.Float64("actual_weight", actualWeight),
			zap.Float64("original_weight", event.Quantity),
			zap.String("note", note),
		)

		return nil
	})
}
