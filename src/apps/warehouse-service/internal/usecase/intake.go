package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/events"
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

func (uc *IntakeUseCase) CreateIntakeFromHarvest(ctx context.Context, pickup domain.PickupRequest) (*domain.Intake, error) {
	var created domain.Intake
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return createIntakeFromPickup(tx, pickup, &created)
	})
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func createIntakeFromPickup(tx *gorm.DB, pickup domain.PickupRequest, created *domain.Intake) error {
	var existing domain.Intake
	if err := tx.Where("pickup_id = ?", pickup.ID).Take(&existing).Error; err == nil {
		*created = existing
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	intake := domain.Intake{
		ID:          uuid.New().String(),
		HarvestID:   pickup.HarvestID,
		PickupID:    pickup.ID,
		WarehouseID: pickup.WarehouseID,
		CoffeeType:  pickup.CoffeeType,
		OriginCode:  pickup.OriginCode,
		Quantity:    pickup.Quantity,
		Status:      domain.IntakeStatusUnassigned,
	}
	if err := tx.Create(&intake).Error; err != nil {
		return err
	}
	*created = intake
	return nil
}

func (uc *IntakeUseCase) ProcessHarvestEvent(ctx context.Context, msgID string, payload []byte) error {
	event, _, err := DecodeHarvestCreated(payload)
	if err != nil {
		return err
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

func DecodeHarvestCreated(payload []byte) (HarvestCreatedEvent, events.Metadata, error) {
	cloudEvent, err := events.ParseCloudEvent(payload)
	if err != nil {
		return HarvestCreatedEvent{}, events.Metadata{}, fmt.Errorf("failed to parse harvest CloudEvent: %w", err)
	}
	event, err := events.DataAs[HarvestCreatedEvent](cloudEvent)
	if err != nil {
		return HarvestCreatedEvent{}, events.Metadata{}, fmt.Errorf("failed to unmarshal harvest event data: %w", err)
	}
	metadata := events.Metadata{
		EventID:       cloudEvent.ID(),
		CorrelationID: events.ExtensionString(cloudEvent, "correlationid"),
		CausationID:   events.ExtensionString(cloudEvent, "causationid"),
		TraceID:       events.ExtensionString(cloudEvent, "traceid"),
		HarvestID:     events.ExtensionString(cloudEvent, "harvestid"),
		FarmID:        events.ExtensionString(cloudEvent, "farmid"),
		WarehouseID:   events.ExtensionString(cloudEvent, "warehouseid"),
	}
	if metadata.HarvestID == "" {
		metadata.HarvestID = event.HarvestID
	}
	if metadata.CorrelationID == "" {
		metadata.CorrelationID = event.HarvestID
	}
	return event, metadata, nil
}
