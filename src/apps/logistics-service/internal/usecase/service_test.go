package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"RuntimeRoasters/apps/logistics-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) Publish(ctx context.Context, topic string, key string, value interface{}) error {
	args := m.Called(ctx, topic, key, value)
	return args.Error(0)
}

func (m *MockProducer) Close() error {
	return nil
}

func setupServiceTest(t *testing.T) (*gorm.DB, *Service, *MockProducer) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.SetMaxOpenConns(1)
	assert.NoError(t, db.AutoMigrate(&domain.Driver{}, &domain.Vehicle{}, &domain.Shipment{}, &domain.Location{}, &domain.ProcessedKafkaMessage{}, &domain.InboxEvent{}))
	producer := new(MockProducer)
	producer.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	service := NewService(db, nil, producer, events.TopicLogisticsDeliveryAssigned, events.TopicLogisticsDeliveryCompleted, events.TopicLogisticsGPSUpdated, "store-1")
	return db, service, producer
}

func TestFarmPickupStateMachineRequiresReturnBeforeAvailability(t *testing.T) {
	db, service, _ := setupServiceTest(t)
	ctx := identity.InjectContext(context.Background(), identity.Claims{Subject: "driver@runtimeroasters.com", Email: "driver@runtimeroasters.com", Role: "DRIVER"})
	now := time.Now()
	driver := domain.Driver{ID: "driver-1", UserID: "driver@runtimeroasters.com", Name: "Driver", Phone: "1", VehicleID: "vehicle-1", Status: domain.DriverStatusBusy, CurrentShipmentID: "shipment-1", IsAvailable: false}
	vehicle := domain.Vehicle{ID: "vehicle-1", PlateNumber: "HN-001", Type: "TRUCK", Status: domain.DriverStatusBusy}
	shipment := domain.Shipment{ID: "shipment-1", Type: domain.ShipmentTypeFarmPickup, Status: domain.ShipmentStatusAssigned, CurrentLeg: domain.ShipmentLegOutbound, FarmID: "farm-1", HarvestID: "harvest-1", WarehouseID: "warehouse-1", DriverID: driver.ID, VehicleID: vehicle.ID, AssignedAt: &now}
	assert.NoError(t, db.Create(&driver).Error)
	assert.NoError(t, db.Create(&vehicle).Error)
	assert.NoError(t, db.Create(&shipment).Error)
	assert.NoError(t, db.Model(&domain.Driver{}).Where("id = ?", driver.ID).Updates(map[string]interface{}{"is_available": false, "status": domain.DriverStatusBusy}).Error)

	_, err := service.DepartShipment(ctx, shipment.ID)
	assert.NoError(t, err)
	_, err = service.ArriveShipment(ctx, shipment.ID)
	assert.NoError(t, err)
	_, err = service.ConfirmLoad(ctx, shipment.ID)
	assert.NoError(t, err)

	var mid domain.Driver
	assert.NoError(t, db.Where("id = ?", driver.ID).Take(&mid).Error)
	assert.False(t, mid.IsAvailable)

	_, err = service.ReturnShipment(ctx, shipment.ID)
	assert.NoError(t, err)
	_, err = service.ArriveShipment(ctx, shipment.ID)
	assert.NoError(t, err)

	var updated domain.Shipment
	assert.NoError(t, db.Where("id = ?", shipment.ID).Take(&updated).Error)
	assert.Equal(t, domain.ShipmentStatusArrivedWarehouse, updated.Status)
	var returnedDriver domain.Driver
	assert.NoError(t, db.Where("id = ?", driver.ID).Take(&returnedDriver).Error)
	assert.True(t, returnedDriver.IsAvailable)
	assert.Equal(t, domain.DriverStatusIdle, returnedDriver.Status)
}

func TestInvalidTransitionDenied(t *testing.T) {
	db, service, _ := setupServiceTest(t)
	ctx := context.Background()
	shipment := domain.Shipment{ID: "shipment-invalid", Type: domain.ShipmentTypeRetailDelivery, Status: domain.ShipmentStatusAssigned, WarehouseID: "warehouse-1", DestinationStoreID: "store-1"}
	assert.NoError(t, db.Create(&shipment).Error)

	_, err := service.ReturnShipment(ctx, shipment.ID)
	assert.Error(t, err)
}

func TestDriverCannotUpdateAnotherShipmentLocation(t *testing.T) {
	db, service, _ := setupServiceTest(t)
	ctx := identity.InjectContext(context.Background(), identity.Claims{Subject: "driver@runtimeroasters.com", Email: "driver@runtimeroasters.com", Role: "DRIVER"})
	driver := domain.Driver{ID: "driver-1", UserID: "driver@runtimeroasters.com", Name: "Driver", Phone: "1", IsAvailable: false}
	other := domain.Shipment{ID: "shipment-other", Type: domain.ShipmentTypeRetailDelivery, Status: domain.ShipmentStatusAssigned, DriverID: "driver-2", DestinationStoreID: "store-1"}
	assert.NoError(t, db.Create(&driver).Error)
	assert.NoError(t, db.Create(&other).Error)

	err := service.UpdateDriverLocation(ctx, "", other.ID, 21.01, 105.85)
	assert.Error(t, err)
}

func TestWarehouseAssignedDeliveryCreatesAssignedShipmentAndMarksFleetBusy(t *testing.T) {
	db, service, _ := setupServiceTest(t)
	driver := domain.Driver{ID: "driver-1", UserID: "driver@runtimeroasters.com", Name: "Driver", Phone: "1", IsAvailable: true, Status: domain.DriverStatusIdle}
	vehicle := domain.Vehicle{ID: "vehicle-1", PlateNumber: "HN-001", Type: "TRUCK", Status: domain.DriverStatusIdle}
	require.NoError(t, db.Create(&driver).Error)
	require.NoError(t, db.Create(&vehicle).Error)

	event := events.LogisticsShipmentAssigned{
		EventID:    "evt-1",
		OrderID:    "order-1",
		StoreID:    "store-1",
		DriverID:   driver.ID,
		VehicleID:  vehicle.ID,
		OccurredAt: time.Now(),
	}
	payload := logisticsCloudEvent(t, events.TopicLogisticsDeliveryAssigned, event, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: event.OrderID,
		OrderID:       event.OrderID,
		StoreID:       event.StoreID,
		WarehouseID:   "WAREHOUSE-HN-001",
		DriverID:      event.DriverID,
		VehicleID:     event.VehicleID,
		OccurredAt:    event.OccurredAt,
	})

	require.NoError(t, service.HandleWarehouseEvent(context.Background(), kafkaMessage(events.TopicLogisticsDeliveryAssigned, event.OrderID, payload)))

	var shipment domain.Shipment
	require.NoError(t, db.Where("order_id = ?", event.OrderID).Take(&shipment).Error)
	assert.Equal(t, domain.ShipmentStatusAssigned, shipment.Status)
	assert.Equal(t, domain.ShipmentTypeRetailDelivery, shipment.Type)
	assert.Equal(t, event.DriverID, shipment.DriverID)
	assert.Equal(t, event.VehicleID, shipment.VehicleID)

	var updatedDriver domain.Driver
	require.NoError(t, db.Where("id = ?", driver.ID).Take(&updatedDriver).Error)
	assert.False(t, updatedDriver.IsAvailable)
	assert.Equal(t, domain.DriverStatusBusy, updatedDriver.Status)
	assert.Equal(t, shipment.ID, updatedDriver.CurrentShipmentID)

	var updatedVehicle domain.Vehicle
	require.NoError(t, db.Where("id = ?", vehicle.ID).Take(&updatedVehicle).Error)
	assert.Equal(t, domain.DriverStatusBusy, updatedVehicle.Status)
}

func logisticsCloudEvent(t *testing.T, topic string, data interface{}, metadata events.Metadata) []byte {
	t.Helper()
	cloudEvent, err := events.NewCloudEvent(context.Background(), topic, "/tests/logistics-service", "orders/"+metadata.OrderID, data, metadata)
	require.NoError(t, err)
	payload, err := json.Marshal(cloudEvent)
	require.NoError(t, err)
	return payload
}

func kafkaMessage(topic string, key string, payload []byte) kafkago.Message {
	return kafkago.Message{Topic: topic, Key: []byte(key), Value: payload}
}
