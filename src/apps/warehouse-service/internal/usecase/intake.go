package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
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
	cloudEvent, err := events.ParseCloudEvent(payload)
	if err != nil {
		return fmt.Errorf("failed to parse harvest CloudEvent: %w", err)
	}
	event, err := events.DataAs[HarvestCreatedEvent](cloudEvent)
	if err != nil {
		return fmt.Errorf("failed to unmarshal harvest event data: %w", err)
	}

	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Idempotency Check
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", msgID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 2. Create Intake (Unassigned raw material)
		intake := domain.Intake{
			ID:         uuid.New().String(),
			HarvestID:  event.HarvestID,
			CoffeeType: event.CoffeeType,
			OriginCode: event.OriginCode,
			Quantity:   event.Quantity,
			Status:     domain.IntakeStatusUnassigned,
		}

		if err := tx.Create(&intake).Error; err != nil {
			return err
		}

		// 3. Mark as processed
		inbox := domain.InboxEvent{
			ID:          uuid.New().String(),
			MessageID:   msgID,
			ProcessedAt: time.Now(),
		}
		return tx.Create(&inbox).Error
	})
}
