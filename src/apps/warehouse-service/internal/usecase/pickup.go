package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type PickupUseCase struct {
	db                 *gorm.DB
	producer           kafka.Producer
	pickupTopic        string
	notificationTopic  string
	intakeCreatedTopic string
	defaultWarehouseID string
}

func NewPickupUseCase(db *gorm.DB, producer kafka.Producer, pickupTopic string, notificationTopic string, intakeCreatedTopic string, defaultWarehouseID string) *PickupUseCase {
	return &PickupUseCase{db: db, producer: producer, pickupTopic: pickupTopic, notificationTopic: notificationTopic, intakeCreatedTopic: intakeCreatedTopic, defaultWarehouseID: defaultWarehouseID}
}

func (uc *PickupUseCase) ProcessHarvestEvent(ctx context.Context, msgID string, payload []byte) error {
	event, metadata, err := DecodeHarvestCreated(payload)
	if err != nil {
		return err
	}
	now := time.Now()
	return uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", msgID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		pickup := domain.PickupRequest{}
		if err := tx.Where("harvest_id = ?", event.HarvestID).Take(&pickup).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			pickup = domain.PickupRequest{
				ID:               uuid.NewString(),
				HarvestID:        event.HarvestID,
				FarmID:           firstNonEmpty(metadata.FarmID, "farm-"+event.HarvestID),
				WarehouseID:      firstNonEmpty(metadata.WarehouseID, uc.defaultWarehouseID),
				OriginLocationID: firstNonEmpty(metadata.FarmID, "farm-"+event.HarvestID),
				Quantity:         event.Quantity,
				CoffeeType:       event.CoffeeType,
				OriginCode:       event.OriginCode,
				Status:           domain.PickupRequestStatusRequested,
				NotificationID:   uuid.NewString(),
				CreatedAt:        now,
				UpdatedAt:        now,
			}
			if err := tx.Create(&pickup).Error; err != nil {
				return err
			}
		}

		if err := markInbox(tx, msgID, events.TopicFarmHarvestCreated); err != nil {
			return err
		}
		if err := uc.publishPickupRequested(ctx, pickup, metadata.CorrelationID, metadata.EventID, metadata.TraceID); err != nil {
			return err
		}
		return uc.publishNotification(ctx, pickup, metadata.CorrelationID, metadata.EventID, metadata.TraceID)
	})
}

func (uc *PickupUseCase) ListPickupRequests(ctx context.Context) ([]domain.PickupRequest, error) {
	var pickups []domain.PickupRequest
	query := uc.db.WithContext(ctx).Order("created_at DESC")
	_, warehouseIDs, allWarehouses := identity.WarehouseScopeFromContext(ctx)
	if !allWarehouses {
		if len(warehouseIDs) == 0 {
			return []domain.PickupRequest{}, nil
		}
		query = query.Where("warehouse_id IN ?", warehouseIDs)
	}
	err := query.Find(&pickups).Error
	return pickups, err
}

func (uc *PickupUseCase) DispatchPickup(ctx context.Context, id string, driverID, vehicleID string) (*domain.PickupRequest, error) {
	var pickup domain.PickupRequest
	now := time.Now()
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Take(&pickup).Error; err != nil {
			return err
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(ctx); !allWarehouses && !claims.CanAccessWarehouse(pickup.WarehouseID) {
			return fmt.Errorf("warehouse access denied")
		}
		if pickup.Status == domain.PickupRequestStatusReceived || pickup.Status == domain.PickupRequestStatusCancelled {
			return fmt.Errorf("pickup request %s cannot be dispatched from status %s", id, pickup.Status)
		}
		if pickup.Status != domain.PickupRequestStatusDispatched {
			if err := tx.Model(&pickup).Updates(map[string]interface{}{
				"status":        domain.PickupRequestStatusDispatched,
				"dispatched_at": &now,
			}).Error; err != nil {
				return err
			}
			pickup.Status = domain.PickupRequestStatusDispatched
			pickup.DispatchedAt = &now
		}
		return uc.publishPickupRequested(ctx, pickup, pickup.HarvestID, driverID, vehicleID)
	})
	if err != nil {
		return nil, err
	}
	return &pickup, nil
}

func (uc *PickupUseCase) HandlePickupArrived(ctx context.Context, msg kafkago.Message) error {
	cloudEvent, err := events.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	event, err := events.DataAs[events.LogisticsPickupStatusChanged](cloudEvent)
	if err != nil {
		return err
	}
	return uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.InboxEvent
		messageID := kafka.MessageID(msg)
		if err := tx.Where("message_id = ?", messageID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		updates := map[string]interface{}{
			"status":      domain.PickupRequestStatusArrivedWarehouse,
			"shipment_id": event.ShipmentID,
		}
		query := tx.Model(&domain.PickupRequest{})
		if event.HarvestID != "" {
			query = query.Where("harvest_id = ?", event.HarvestID)
		} else {
			query = query.Where("shipment_id = ?", event.ShipmentID)
		}
		result := query.Updates(updates)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("pickup request not found for shipment %s harvest %s", event.ShipmentID, event.HarvestID)
		}
		if err := markInbox(tx, messageID, msg.Topic); err != nil {
			return err
		}
		return nil
	})
}

func (uc *PickupUseCase) ReceivePickup(ctx context.Context, id string) (*domain.Intake, error) {
	var intake domain.Intake
	var pickup domain.PickupRequest
	now := time.Now()
	shouldPublish := false
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Take(&pickup).Error; err != nil {
			return err
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(ctx); !allWarehouses && !claims.CanAccessWarehouse(pickup.WarehouseID) {
			return fmt.Errorf("warehouse access denied")
		}
		if pickup.Status != domain.PickupRequestStatusArrivedWarehouse && pickup.Status != domain.PickupRequestStatusReceived {
			return fmt.Errorf("pickup request %s must arrive at warehouse before receipt, current: %s", id, pickup.Status)
		}
		if err := createIntakeFromPickup(tx, pickup, &intake); err != nil {
			return err
		}
		if pickup.Status != domain.PickupRequestStatusReceived {
			if err := tx.Model(&pickup).Updates(map[string]interface{}{
				"status":      domain.PickupRequestStatusReceived,
				"received_at": &now,
			}).Error; err != nil {
				return err
			}
			pickup.Status = domain.PickupRequestStatusReceived
			pickup.ReceivedAt = &now
			shouldPublish = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if shouldPublish {
		if err := uc.publishIntakeCreated(ctx, pickup, intake); err != nil {
			return nil, err
		}
	}
	return &intake, nil
}

func (uc *PickupUseCase) publishPickupRequested(ctx context.Context, pickup domain.PickupRequest, correlationID string, causationID string, traceID string) error {
	payload := events.WarehousePickupRequested{
		EventID:     uuid.NewString(),
		PickupID:    pickup.ID,
		HarvestID:   pickup.HarvestID,
		FarmID:      pickup.FarmID,
		WarehouseID: pickup.WarehouseID,
		Quantity:    pickup.Quantity,
		Status:      string(pickup.Status),
		OccurredAt:  time.Now(),
	}
	cloudEvent, err := events.NewCloudEvent(ctx, uc.pickupTopic, events.SourceWarehouseService, fmt.Sprintf("pickups/%s", pickup.ID), payload, events.Metadata{
		EventID:       payload.EventID,
		CorrelationID: firstNonEmpty(correlationID, pickup.HarvestID),
		CausationID:   causationID,
		TraceID:       traceID,
		HarvestID:     pickup.HarvestID,
		FarmID:        pickup.FarmID,
		WarehouseID:   pickup.WarehouseID,
		OccurredAt:    payload.OccurredAt,
	})
	if err != nil {
		return err
	}
	return uc.producer.Publish(ctx, uc.pickupTopic, pickup.ID, cloudEvent)
}

func (uc *PickupUseCase) publishNotification(ctx context.Context, pickup domain.PickupRequest, correlationID string, causationID string, traceID string) error {
	if uc.notificationTopic == "" {
		return nil
	}
	payload := events.NotificationCreated{
		EventID:     pickup.NotificationID,
		Role:        "WAREHOUSE_MGR",
		WarehouseID: pickup.WarehouseID,
		Title:       "Pickup requested",
		Message:     fmt.Sprintf("Harvest %s is ready for warehouse pickup", pickup.HarvestID),
		Severity:    "INFO",
		OccurredAt:  time.Now(),
	}
	cloudEvent, err := events.NewCloudEvent(ctx, uc.notificationTopic, events.SourceWarehouseService, fmt.Sprintf("notifications/%s", pickup.NotificationID), payload, events.Metadata{
		EventID:       payload.EventID,
		CorrelationID: firstNonEmpty(correlationID, pickup.HarvestID),
		CausationID:   causationID,
		TraceID:       traceID,
		HarvestID:     pickup.HarvestID,
		FarmID:        pickup.FarmID,
		WarehouseID:   pickup.WarehouseID,
		OccurredAt:    payload.OccurredAt,
	})
	if err != nil {
		return err
	}
	return uc.producer.Publish(ctx, uc.notificationTopic, pickup.ID, cloudEvent)
}

func (uc *PickupUseCase) publishIntakeCreated(ctx context.Context, pickup domain.PickupRequest, intake domain.Intake) error {
	payload := events.WarehouseIntakeCreated{
		EventID:     uuid.NewString(),
		IntakeID:    intake.ID,
		HarvestID:   pickup.HarvestID,
		FarmID:      pickup.FarmID,
		WarehouseID: pickup.WarehouseID,
		Quantity:    intake.Quantity,
		OccurredAt:  time.Now(),
	}
	cloudEvent, err := events.NewCloudEvent(ctx, uc.intakeCreatedTopic, events.SourceWarehouseService, fmt.Sprintf("intakes/%s", intake.ID), payload, events.Metadata{
		EventID:       payload.EventID,
		CorrelationID: pickup.HarvestID,
		HarvestID:     pickup.HarvestID,
		FarmID:        pickup.FarmID,
		WarehouseID:   pickup.WarehouseID,
		OccurredAt:    payload.OccurredAt,
	})
	if err != nil {
		return err
	}
	return uc.producer.Publish(ctx, uc.intakeCreatedTopic, pickup.HarvestID, cloudEvent)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
