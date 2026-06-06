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

func TestCreateStoreDerivesCodeForExistingAdminContract(t *testing.T) {
	db := setupRetailDB(t)
	service := NewService(db, nil, events.TopicRetailOrderCreated)

	created, err := service.CreateStore(context.Background(), domain.Store{
		Name: "New Demo Store", City: "Hanoi", Address: "1 Demo Street",
		Status: domain.StoreStatusActive,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, created.Code)
	assert.Contains(t, created.Code, "STORE-")
}

func TestSeedDataIsIdempotentAndRebuildsAvailability(t *testing.T) {
	db := setupRetailDB(t)
	service := NewService(db, nil, events.TopicRetailOrderCreated)
	users := map[string]string{
		"mgr.hn.hoankiem@runtimeroasters.com": "user-hk",
		"mgr.hn.caugiay@runtimeroasters.com":  "user-cg",
		"mgr.hcm.d1@runtimeroasters.com":      "user-d1",
		"mgr.hcm.d7@runtimeroasters.com":      "user-d7",
		"mgr.dn.haichau@runtimeroasters.com":  "user-hc",
	}

	first, err := service.SeedData(context.Background(), false, users)
	require.NoError(t, err)
	assert.True(t, first.Success)
	firstStatus, err := service.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.True(t, firstStatus.Seeded)
	assert.Equal(t, int64(42), firstStatus.RecordCounts["menu_items"])
	assert.Equal(t, int64(25), firstStatus.RecordCounts["sale_items"])
	assert.Equal(t, int64(210), firstStatus.RecordCounts["store_menu_inventories"])

	_, err = service.SeedData(context.Background(), false, users)
	require.NoError(t, err)
	secondStatus, err := service.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.Equal(t, firstStatus.RecordCounts, secondStatus.RecordCounts)

	var lots []domain.InventoryLot
	require.NoError(t, db.Find(&lots).Error)
	for _, lot := range lots {
		var ledgerBalance float64
		require.NoError(t, db.Model(&domain.StockMovement{}).
			Where("inventory_lot_id = ?", lot.ID).
			Select("COALESCE(SUM(quantity_delta), 0)").
			Scan(&ledgerBalance).Error)
		assert.InDelta(t, lot.AvailableQuantity, ledgerBalance, 0.001, lot.ID)
	}

	var espresso domain.MenuItem
	require.NoError(t, db.Where("id = ?", "MI-ESPRESSO-DOUBLE-S").Take(&espresso).Error)
	var availability domain.StoreMenuInventory
	require.NoError(t, db.Where("store_id = ? AND menu_item_id = ?",
		"11111111-1111-1111-1111-111111111101", espresso.ID).
		Take(&availability).Error)
	assert.Positive(t, availability.AvailableUnits)

	userMovement := domain.StockMovement{
		ID: uuid.NewString(), StoreID: "11111111-1111-1111-1111-111111111101",
		InventoryLotID: "LOT-HK-AR-001", StockSKU: "BEAN-ARABICA-ROASTED",
		MovementType: domain.StockMovementSold, QuantityDelta: -18, Unit: "GRAM",
		ReferenceType: domain.StockReferenceSaleItem, ReferenceID: uuid.NewString(),
		OccurredAt: time.Now(),
	}
	require.NoError(t, db.Create(&userMovement).Error)
	var beforeReseed domain.InventoryLot
	require.NoError(t, db.Where("id = ?", userMovement.InventoryLotID).Take(&beforeReseed).Error)

	_, err = service.SeedData(context.Background(), false, users)
	require.NoError(t, err)
	var afterReseed domain.InventoryLot
	require.NoError(t, db.Where("id = ?", userMovement.InventoryLotID).Take(&afterReseed).Error)
	assert.InDelta(t, beforeReseed.AvailableQuantity-18, afterReseed.AvailableQuantity, 0.001)
}

func setupRetailDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	require.NoError(t, db.AutoMigrate(
		&domain.Store{}, &domain.Menu{}, &domain.MenuItem{}, &domain.Order{},
		&domain.InventoryLot{}, &domain.Sale{}, &domain.SaleItem{},
		&domain.StockMovement{}, &domain.StoreMenuInventory{},
		&domain.OutboxEvent{}, &domain.InboxEvent{},
	))
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
