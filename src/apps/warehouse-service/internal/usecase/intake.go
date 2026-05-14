package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
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

func (uc *IntakeUseCase) ProcessHarvestEvent(ctx context.Context, msgID string, payload []byte) error {
	var event HarvestCreatedEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal harvest event: %w", err)
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Idempotency Check (Inbox Pattern)
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", msgID).First(&existing).Error; err == nil {
			fmt.Printf("[INTAKE] Message %s already processed. Skipping.\n", msgID)
			return nil
		}

		// 2. Business Logic: Random Intake Simulation
		// PO Request: Some batches match, some lack weight (Random)
		actualWeight := event.Quantity
		note := ""
		
		// Simulate random loss (10% chance to lose weight between 3% and 10%)
		if rand.Float64() < 0.2 { // 20% chance of anomaly
			lossPercent := 0.03 + rand.Float64()*0.07 // 3% to 10% loss
			actualWeight = event.Quantity * (1 - lossPercent)
			
			deviation := (event.Quantity - actualWeight) / event.Quantity * 100
			if deviation > 2.0 {
				note = fmt.Sprintf("Hao hụt vận chuyển bất thường: %.2f%%. Yêu cầu kiểm tra xe hàng.", deviation)
			}
		}

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

		fmt.Printf("[INTAKE] Batch created: %s, Weight: %.2fkg (Original: %.2fkg), Note: %s\n", 
			batchID, actualWeight, event.Quantity, note)

		return nil
	})
}
