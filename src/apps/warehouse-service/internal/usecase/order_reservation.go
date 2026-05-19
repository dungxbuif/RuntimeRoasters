package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type OrderReservationUseCase struct {
	db            *gorm.DB
	producer      kafka.Producer
	reservedTopic string
	failedTopic   string
}

func NewOrderReservationUseCase(db *gorm.DB, producer kafka.Producer, reservedTopic string, failedTopic string) *OrderReservationUseCase {
	return &OrderReservationUseCase{db: db, producer: producer, reservedTopic: reservedTopic, failedTopic: failedTopic}
}

func (uc *OrderReservationUseCase) ProcessOrderCreated(ctx context.Context, msg kafkago.Message) error {
	messageID := fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset)
	event, err := decodeOrderReservationEvent(msg.Topic, msg.Value)
	if err != nil {
		return err
	}

	return uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", messageID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := tx.Where("event_type = ? AND message_id = ?", "order_reserved", event.OrderID).Take(&existing).Error; err == nil {
			return markInbox(tx, messageID, msg.Topic)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		reserved, failed := uc.reserve(tx, event)
		if reserved != nil {
			if err := uc.producer.Publish(ctx, uc.reservedTopic, event.OrderID, reserved); err != nil {
				return err
			}
		}
		if failed != nil {
			if err := uc.producer.Publish(ctx, uc.failedTopic, event.OrderID, failed); err != nil {
				return err
			}
		}

		if err := markInbox(tx, messageID, msg.Topic); err != nil {
			return err
		}
		return markInbox(tx, event.OrderID, "order_reserved")
	})
}

func decodeOrderReservationEvent(topic string, payload []byte) (events.RetailOrderCreated, error) {
	switch topic {
	case events.TopicPaymentCompleted:
		var event events.PaymentCompleted
		if err := json.Unmarshal(payload, &event); err != nil {
			return events.RetailOrderCreated{}, err
		}
		return events.RetailOrderCreated{
			EventID:       event.EventID,
			OrderID:       event.OrderID,
			StoreID:       event.StoreID,
			Items:         event.Items,
			TotalAmount:   event.Amount,
			PaymentMethod: event.Provider,
			OccurredAt:    event.OccurredAt,
		}, nil
	default:
		var event events.RetailOrderCreated
		if err := json.Unmarshal(payload, &event); err != nil {
			return events.RetailOrderCreated{}, err
		}
		return event, nil
	}
}

func markInbox(tx *gorm.DB, messageID string, eventType string) error {
	return tx.Create(&domain.InboxEvent{
		ID:          uuid.NewString(),
		MessageID:   messageID,
		EventType:   eventType,
		ProcessedAt: time.Now(),
	}).Error
}

func (uc *OrderReservationUseCase) reserve(tx *gorm.DB, event events.RetailOrderCreated) (*events.WarehouseStockReserved, *events.WarehouseStockReservationFailed) {
	if len(event.Items) == 0 {
		return nil, &events.WarehouseStockReservationFailed{
			EventID:    uuid.NewString(),
			OrderID:    event.OrderID,
			StoreID:    event.StoreID,
			Reason:     "order has no items",
			OccurredAt: time.Now(),
		}
	}

	item := event.Items[0]
	var inventory domain.Inventory
	if err := tx.Where("sku = ?", item.SKU).Take(&inventory).Error; err != nil || inventory.AvailableQuantity < item.Quantity {
		return nil, &events.WarehouseStockReservationFailed{
			EventID:    uuid.NewString(),
			OrderID:    event.OrderID,
			StoreID:    event.StoreID,
			Reason:     "insufficient stock",
			OccurredAt: time.Now(),
		}
	}

	if err := tx.Model(&inventory).Update("available_quantity", gorm.Expr("available_quantity - ?", item.Quantity)).Error; err != nil {
		return nil, &events.WarehouseStockReservationFailed{
			EventID:    uuid.NewString(),
			OrderID:    event.OrderID,
			StoreID:    event.StoreID,
			Reason:     err.Error(),
			OccurredAt: time.Now(),
		}
	}

	return &events.WarehouseStockReserved{
		EventID:    uuid.NewString(),
		OrderID:    event.OrderID,
		StoreID:    event.StoreID,
		SKU:        item.SKU,
		Quantity:   item.Quantity,
		OccurredAt: time.Now(),
	}, nil
}
