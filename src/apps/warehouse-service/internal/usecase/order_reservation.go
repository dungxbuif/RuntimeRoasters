package usecase

import (
	"context"
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
	messageID := kafka.MessageID(msg)
	request, err := decodeOrderReservationEvent(msg.Topic, msg.Value)
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
		if err := tx.Where("event_type = ? AND message_id = ?", "order_reserved", request.Order.OrderID).Take(&existing).Error; err == nil {
			return markInbox(tx, messageID, msg.Topic)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		reserved, failed := uc.reserve(tx, request.Order)
		if reserved != nil {
			if err := uc.publishReservationEvent(ctx, uc.reservedTopic, reserved.OrderID, reserved, request.CorrelationID, request.CausationID, request.TraceID); err != nil {
				return err
			}
		}
		if failed != nil {
			if err := uc.publishReservationEvent(ctx, uc.failedTopic, failed.OrderID, failed, request.CorrelationID, request.CausationID, request.TraceID); err != nil {
				return err
			}
		}

		if err := markInbox(tx, messageID, msg.Topic); err != nil {
			return err
		}
		return markInbox(tx, request.Order.OrderID, "order_reserved")
	})
}

type orderReservationRequest struct {
	Order         events.RetailOrderCreated
	CorrelationID string
	CausationID   string
	TraceID       string
}

func decodeOrderReservationEvent(topic string, payload []byte) (orderReservationRequest, error) {
	cloudEvent, err := events.ParseCloudEvent(payload)
	if err != nil {
		return orderReservationRequest{}, err
	}
	correlationID := events.ExtensionString(cloudEvent, "correlationid")
	if correlationID == "" {
		correlationID = events.ExtensionString(cloudEvent, "orderid")
	}
	traceID := events.ExtensionString(cloudEvent, "traceid")
	switch topic {
	case events.TopicPaymentCompleted, events.TopicPaymentSimulatedCompleted:
		event, err := events.DataAs[events.PaymentCompleted](cloudEvent)
		if err != nil {
			return orderReservationRequest{}, err
		}
		return orderReservationRequest{Order: events.RetailOrderCreated{
			EventID:       event.EventID,
			OrderID:       event.OrderID,
			StoreID:       event.StoreID,
			Items:         event.Items,
			TotalAmount:   event.Amount,
			PaymentMethod: event.Provider,
			OccurredAt:    event.OccurredAt,
		}, CorrelationID: correlationID, CausationID: cloudEvent.ID(), TraceID: traceID}, nil
	default:
		event, err := events.DataAs[events.RetailOrderCreated](cloudEvent)
		if err != nil {
			return orderReservationRequest{}, err
		}
		if correlationID == "" {
			correlationID = event.OrderID
		}
		return orderReservationRequest{Order: event, CorrelationID: correlationID, CausationID: cloudEvent.ID(), TraceID: traceID}, nil
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

func (uc *OrderReservationUseCase) publishReservationEvent(ctx context.Context, topic string, key string, payload interface{}, correlationID string, causationID string, traceID string) error {
	metadata := warehouseReservationMetadata(payload)
	if correlationID == "" {
		correlationID = metadata.OrderID
	}
	metadata.CorrelationID = correlationID
	metadata.CausationID = causationID
	metadata.TraceID = traceID
	cloudEvent, err := events.NewCloudEvent(ctx, topic, events.SourceWarehouseService, fmt.Sprintf("orders/%s", metadata.OrderID), payload, metadata)
	if err != nil {
		return err
	}
	return uc.producer.Publish(ctx, topic, key, cloudEvent)
}

func warehouseReservationMetadata(payload interface{}) events.Metadata {
	switch event := payload.(type) {
	case *events.WarehouseStockReserved:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID}
	case *events.WarehouseStockReservationFailed:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID}
	default:
		return events.Metadata{}
	}
}
