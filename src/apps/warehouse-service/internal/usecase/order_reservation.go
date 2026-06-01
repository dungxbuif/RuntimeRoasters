package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type OrderReservationUseCase struct {
	db                     *gorm.DB
	rdb                    *redis.Client
	locker                 inventoryLocker
	producer               kafka.Producer
	reservedTopic          string
	failedTopic            string
	dispatchRequestedTopic string
}

func NewOrderReservationUseCase(db *gorm.DB, rdb *redis.Client, producer kafka.Producer, reservedTopic string, failedTopic string) *OrderReservationUseCase {
	return &OrderReservationUseCase{
		db:                     db,
		rdb:                    rdb,
		locker:                 redisInventoryLocker{client: rdb},
		producer:               producer,
		reservedTopic:          reservedTopic,
		failedTopic:            failedTopic,
		dispatchRequestedTopic: events.TopicWarehouseDispatchRequested,
	}
}

type inventoryLocker interface {
	Acquire(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
	Release(ctx context.Context, key string) error
}

type redisInventoryLocker struct {
	client *redis.Client
}

func (l redisInventoryLocker) Acquire(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	if l.client == nil {
		return false, errors.New("inventory lock client is not configured")
	}
	return l.client.SetNX(ctx, key, value, ttl).Result()
}

func (l redisInventoryLocker) Release(ctx context.Context, key string) error {
	if l.client == nil {
		return nil
	}
	return l.client.Del(ctx, key).Err()
}

func (uc *OrderReservationUseCase) ProcessOrderCreated(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)

	// Handle Completion Policy
	if msg.Topic == events.TopicLogisticsDriverReturnedToBase {
		return uc.finalizeDispatch(ctx, msg.Value)
	}

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

		// 1. Perform Inventory Reservation with Lock
		reserved, failed := uc.reserveWithLock(ctx, tx, request.Order)

		if reserved != nil {
			// 2. Create DispatchRequest for WAREHOUSE_MGR
			itemsJSON, _ := json.Marshal(request.Order.Items)
			dispatchReq := &domain.DispatchRequest{
				ID:          uuid.NewString(),
				OrderID:     reserved.OrderID,
				StoreID:     reserved.StoreID,
				WarehouseID: reserved.WarehouseID,
				Items:       string(itemsJSON),
				Status:      domain.DispatchRequestStatusReserved,
			}
			if err := tx.Create(dispatchReq).Error; err != nil {
				return err
			}

			// 3. Publish Events
			if err := uc.publishReservationEvent(ctx, uc.reservedTopic, reserved.OrderID, reserved, request.CorrelationID, request.CausationID, request.TraceID); err != nil {
				return err
			}

			// Publish Dispatch Requested event
			dispatchMsg := events.WarehouseDispatchRequested{
				EventID:     uuid.NewString(),
				OrderID:     reserved.OrderID,
				StoreID:     reserved.StoreID,
				WarehouseID: reserved.WarehouseID,
				OccurredAt:  time.Now(),
			}
			if err := uc.publishReservationEvent(ctx, uc.dispatchRequestedTopic, reserved.OrderID, &dispatchMsg, request.CorrelationID, request.CausationID, request.TraceID); err != nil {
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
		return nil
	})
}

func (uc *OrderReservationUseCase) reserveWithLock(ctx context.Context, tx *gorm.DB, event events.RetailOrderCreated) (*events.WarehouseStockReserved, *events.WarehouseStockReservationFailed) {
	if len(event.Items) == 0 {
		return nil, &events.WarehouseStockReservationFailed{
			EventID:    uuid.NewString(),
			OrderID:    event.OrderID,
			StoreID:    event.StoreID,
			Reason:     "order has no items",
			OccurredAt: time.Now(),
		}
	}

	// For simplicity in this demo, we assume all items are in the same warehouse or use a default
	warehouseID := "WAREHOUSE-HN-001"

	// Acquire locks for each SKU
	for _, item := range event.Items {
		lockKey := fmt.Sprintf("lock:inventory:%s:%s", warehouseID, item.SKU)
		ok, err := uc.locker.Acquire(ctx, lockKey, event.OrderID, 10*time.Second)
		if err != nil || !ok {
			return nil, &events.WarehouseStockReservationFailed{
				EventID:    uuid.NewString(),
				OrderID:    event.OrderID,
				StoreID:    event.StoreID,
				Reason:     "concurrency lock failure",
				OccurredAt: time.Now(),
			}
		}
		defer uc.locker.Release(ctx, lockKey)
	}

	// Perform actual reservation
	for _, item := range event.Items {
		var inventory domain.Inventory
		if err := tx.Where("sku = ? AND warehouse_id = ?", item.SKU, warehouseID).Take(&inventory).Error; err != nil || inventory.AvailableQuantity < item.Quantity {
			return nil, &events.WarehouseStockReservationFailed{
				EventID:    uuid.NewString(),
				OrderID:    event.OrderID,
				StoreID:    event.StoreID,
				Reason:     fmt.Sprintf("insufficient stock for SKU %s", item.SKU),
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
	}

	return &events.WarehouseStockReserved{
		EventID:     uuid.NewString(),
		OrderID:     event.OrderID,
		StoreID:     event.StoreID,
		WarehouseID: warehouseID,
		SKU:         event.Items[0].SKU, // Simple single-item reporting
		Quantity:    event.Items[0].Quantity,
		OccurredAt:  time.Now(),
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

func (uc *OrderReservationUseCase) finalizeDispatch(ctx context.Context, payload []byte) error {
	cloudEvent, err := events.ParseCloudEvent(payload)
	if err != nil {
		return err
	}
	event, err := events.DataAs[events.LogisticsDeliveryStatusChanged](cloudEvent)
	if err != nil {
		return err
	}

	return uc.db.WithContext(ctx).Model(&domain.DispatchRequest{}).
		Where("order_id = ?", event.OrderID).
		Update("status", domain.DispatchRequestStatusCompleted).Error
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

func warehouseReservationMetadata(payload interface{}) events.Metadata {
	switch event := payload.(type) {
	case *events.WarehouseStockReserved:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID, WarehouseID: event.WarehouseID}
	case *events.WarehouseStockReservationFailed:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID, WarehouseID: event.WarehouseID}
	case *events.WarehouseDispatchRequested:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID, WarehouseID: event.WarehouseID}
	default:
		return events.Metadata{}
	}
}
