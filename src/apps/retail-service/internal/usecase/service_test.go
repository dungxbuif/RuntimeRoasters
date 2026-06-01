package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSagaEventsAdvancePaidOrderFulfillment(t *testing.T) {
	db := setupRetailDB(t)
	service := NewService(db, nil, events.TopicRetailOrderCreated)
	order := domain.Order{ID: "order-1", StoreID: "store-1", Status: domain.OrderStatusPending, TotalAmount: 100}
	require.NoError(t, db.Create(&order).Error)

	steps := []struct {
		topic  string
		data   interface{}
		status domain.OrderStatus
	}{
		{
			topic:  events.TopicPaymentCompleted,
			data:   events.PaymentCompleted{EventID: uuid.NewString(), OrderID: order.ID, StoreID: order.StoreID, Provider: "stripe", Amount: 100, Currency: "USD", OccurredAt: time.Now()},
			status: domain.OrderStatusPaymentCompleted,
		},
		{
			topic:  events.TopicWarehouseStockReserved,
			data:   events.WarehouseStockReserved{EventID: uuid.NewString(), OrderID: order.ID, StoreID: order.StoreID, WarehouseID: "WAREHOUSE-HN-001", SKU: "SKU-AR-001", Quantity: 1, OccurredAt: time.Now()},
			status: domain.OrderStatusReserved,
		},
		{
			topic:  events.TopicWarehouseDispatchRequested,
			data:   events.WarehouseDispatchRequested{EventID: uuid.NewString(), DispatchID: "dispatch-1", OrderID: order.ID, StoreID: order.StoreID, WarehouseID: "WAREHOUSE-HN-001", OccurredAt: time.Now()},
			status: domain.OrderStatusDispatchRequested,
		},
		{
			topic:  events.TopicLogisticsDeliveryAssigned,
			data:   events.LogisticsShipmentAssigned{EventID: uuid.NewString(), ShipmentID: "shipment-1", OrderID: order.ID, StoreID: order.StoreID, DriverID: "driver-1", OccurredAt: time.Now()},
			status: domain.OrderStatusShipping,
		},
		{
			topic:  events.TopicLogisticsDeliveryDriverConfirmed,
			data:   events.LogisticsDeliveryStatusChanged{EventID: uuid.NewString(), ShipmentID: "shipment-1", OrderID: order.ID, StoreID: order.StoreID, DriverID: "driver-1", VehicleID: "vehicle-1", Status: "DELIVERED", OccurredAt: time.Now()},
			status: domain.OrderStatusDelivered,
		},
	}

	for _, step := range steps {
		payload := retailCloudEvent(t, step.topic, step.data, order.ID, order.StoreID)
		require.NoError(t, service.HandleSagaEvent(context.Background(), kafkago.Message{
			Topic: step.topic,
			Key:   []byte(order.ID),
			Value: payload,
		}))

		var updated domain.Order
		require.NoError(t, db.Where("id = ?", order.ID).Take(&updated).Error)
		assert.Equal(t, step.status, updated.Status)
	}

	require.NoError(t, service.ConfirmOrderReceipt(context.Background(), order.ID))
	var completed domain.Order
	require.NoError(t, db.Where("id = ?", order.ID).Take(&completed).Error)
	assert.Equal(t, domain.OrderStatusCompleted, completed.Status)
}

func TestConfirmOrderReceiptRequiresDeliveredOrder(t *testing.T) {
	db := setupRetailDB(t)
	service := NewService(db, nil, events.TopicRetailOrderCreated)
	order := domain.Order{ID: "order-1", StoreID: "store-1", Status: domain.OrderStatusShipping}
	require.NoError(t, db.Create(&order).Error)

	err := service.ConfirmOrderReceipt(context.Background(), order.ID)
	assert.Error(t, err)
}

func setupRetailDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(&domain.Store{}, &domain.Order{}, &domain.OutboxEvent{}, &domain.InboxEvent{}))
	return db
}

func retailCloudEvent(t *testing.T, topic string, data interface{}, orderID string, storeID string) []byte {
	t.Helper()
	cloudEvent, err := events.NewCloudEvent(context.Background(), topic, "/tests/retail-service", "orders/"+orderID, data, events.Metadata{
		EventID:       uuid.NewString(),
		CorrelationID: orderID,
		OrderID:       orderID,
		StoreID:       storeID,
		OccurredAt:    time.Now(),
	})
	require.NoError(t, err)
	payload, err := json.Marshal(cloudEvent)
	require.NoError(t, err)
	return payload
}
