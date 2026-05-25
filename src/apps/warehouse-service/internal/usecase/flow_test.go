package usecase

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockProducer for Kafka
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

func TestFullWarehouseFlow(t *testing.T) {
	// Setup In-Memory DB
	db, _ := gorm.Open(sqlite.Open("file:warehouse-flow?mode=memory&cache=shared"), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	db.AutoMigrate(&domain.PickupRequest{}, &domain.Intake{}, &domain.ProductionBatch{}, &domain.RoastRun{}, &domain.Inventory{}, &domain.InboxEvent{})

	ctx := context.Background()
	logger.InitLogger("development", "debug")

	// 1. Simulate Harvest Event (Pickup Request)
	mockPickupProducer := new(MockProducer)
	mockPickupProducer.On("Publish", mock.Anything, events.TopicWarehousePickupRequested, mock.Anything, mock.Anything).Return(nil)
	mockPickupProducer.On("Publish", mock.Anything, events.TopicNotificationCreated, mock.Anything, mock.Anything).Return(nil)
	mockPickupProducer.On("Publish", mock.Anything, events.TopicWarehouseIntakeCreated, mock.Anything, mock.Anything).Return(nil)
	pickupUC := NewPickupUseCase(db, mockPickupProducer, events.TopicWarehousePickupRequested, events.TopicNotificationCreated, events.TopicWarehouseIntakeCreated, "WAREHOUSE-HN-001")
	payload := mustCloudEvent(t, events.TopicFarmHarvestCreated, "harvests/HV-001", HarvestCreatedEvent{HarvestID: "HV-001", CoffeeType: "ARABICA", OriginCode: "SL", Quantity: 100.0}, events.Metadata{EventID: uuid.NewString(), CorrelationID: "HV-001", HarvestID: "HV-001", OccurredAt: time.Now()})
	err := pickupUC.ProcessHarvestEvent(ctx, "msg-1", payload)
	assert.NoError(t, err)

	// Verify Pickup Created, not Intake
	var pickup domain.PickupRequest
	db.First(&pickup)
	assert.Equal(t, domain.PickupRequestStatusRequested, pickup.Status)
	assert.Equal(t, 100.0, pickup.Quantity)

	var intakeCount int64
	db.Model(&domain.Intake{}).Count(&intakeCount)
	assert.Equal(t, int64(0), intakeCount)

	arrivedPayload := mustCloudEvent(t, events.TopicLogisticsPickupArrivedAtWarehouse, "shipments/SHIP-HV-001", events.LogisticsPickupStatusChanged{
		EventID:     uuid.NewString(),
		ShipmentID:  "SHIP-HV-001",
		HarvestID:   "HV-001",
		FarmID:      pickup.FarmID,
		WarehouseID: pickup.WarehouseID,
		DriverID:    "DRIVER-DEMO-001",
		VehicleID:   "VEHICLE-DEMO-001",
		Status:      "ARRIVED_AT_WAREHOUSE",
		OccurredAt:  time.Now(),
	}, events.Metadata{EventID: uuid.NewString(), CorrelationID: "HV-001", HarvestID: "HV-001", FarmID: pickup.FarmID, WarehouseID: pickup.WarehouseID, OccurredAt: time.Now()})
	err = pickupUC.HandlePickupArrived(ctx, kafkago.Message{Topic: events.TopicLogisticsPickupArrivedAtWarehouse, Key: []byte("SHIP-HV-001"), Value: arrivedPayload})
	assert.NoError(t, err)

	// Receive creates Intake after returned pickup
	intakeCreated, err := pickupUC.ReceivePickup(ctx, pickup.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, intakeCreated.ID)
	secondIntake, err := pickupUC.ReceivePickup(ctx, pickup.ID)
	assert.NoError(t, err)
	assert.Equal(t, intakeCreated.ID, secondIntake.ID)

	var intake domain.Intake
	db.First(&intake)
	assert.Equal(t, domain.IntakeStatusUnassigned, intake.Status)
	assert.Equal(t, 100.0, intake.Quantity)
	assert.Equal(t, pickup.ID, intake.PickupID)
	db.Model(&domain.Intake{}).Count(&intakeCount)
	assert.Equal(t, int64(1), intakeCount)

	// 2. Aggregate Intakes into Batch
	aggUC := NewAggregationUseCase(db)
	batchID, err := aggUC.CreateProductionBatch(ctx, []string{intake.ID})
	assert.NoError(t, err)
	assert.NotEmpty(t, batchID)

	// Verify Batch Created
	var batch domain.ProductionBatch
	db.Preload("Intakes").First(&batch)
	assert.Equal(t, domain.BatchStatusDraft, batch.Status)
	assert.Len(t, batch.Intakes, 1)

	// 3. Start Processing (Async Simulation)
	procUC := NewProcessingUseCase(db)
	err = procUC.StartSimulation(ctx, batch.ID)
	assert.NoError(t, err)

	// Polling for simulation finish (Wait max 15s)
	t.Log("Waiting for async simulation to finish...")
	for i := 0; i < 20; i++ {
		db.First(&batch)
		if batch.Status == domain.BatchStatusReady {
			break
		}
		time.Sleep(1 * time.Second)
	}
	assert.Equal(t, domain.BatchStatusReady, batch.Status)
	assert.Greater(t, batch.TotalOutputWeight, 0.0)
	assert.Less(t, batch.TotalOutputWeight, 100.0) // Loss applied

	// 4. Finalize to Stock
	mockProducer := new(MockProducer)
	mockProducer.On("Publish", mock.Anything, events.TopicWarehouseInventoryUpdated, mock.Anything, mock.Anything).Return(nil)

	invUC := NewInventoryUseCase(db, mockProducer, events.TopicWarehouseInventoryUpdated)
	err = invUC.FinalizeBatch(ctx, batch.ID)
	assert.NoError(t, err)

	// Verify Final State
	db.First(&batch)
	assert.Equal(t, domain.BatchStatusStocked, batch.Status)

	var inventory domain.Inventory
	db.First(&inventory)
	assert.Equal(t, "SL-ARABICA-ROASTED", inventory.SKU)
	assert.Equal(t, batch.TotalOutputWeight, inventory.AvailableQuantity)

	t.Logf("Full flow verified! Final Stock: %.2f kg of %s", inventory.AvailableQuantity, inventory.SKU)
}

func TestHarvestCreatesPickupRequestIdempotently(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file:warehouse-pickup?mode=memory&cache=shared"), &gorm.Config{})
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	db.AutoMigrate(&domain.PickupRequest{}, &domain.Intake{}, &domain.InboxEvent{})
	ctx := context.Background()
	producer := new(MockProducer)
	producer.On("Publish", mock.Anything, events.TopicWarehousePickupRequested, mock.Anything, mock.Anything).Return(nil)
	producer.On("Publish", mock.Anything, events.TopicNotificationCreated, mock.Anything, mock.Anything).Return(nil)
	uc := NewPickupUseCase(db, producer, events.TopicWarehousePickupRequested, events.TopicNotificationCreated, events.TopicWarehouseIntakeCreated, "WAREHOUSE-HN-001")
	payload := mustCloudEvent(t, events.TopicFarmHarvestCreated, "harvests/HV-002", HarvestCreatedEvent{HarvestID: "HV-002", CoffeeType: "ROBUSTA", OriginCode: "BM", Quantity: 50.0}, events.Metadata{EventID: uuid.NewString(), CorrelationID: "HV-002", HarvestID: "HV-002", OccurredAt: time.Now()})

	assert.NoError(t, uc.ProcessHarvestEvent(ctx, "msg-1", payload))
	assert.NoError(t, uc.ProcessHarvestEvent(ctx, "msg-1", payload))

	var pickups int64
	db.Model(&domain.PickupRequest{}).Where("harvest_id = ?", "HV-002").Count(&pickups)
	assert.Equal(t, int64(1), pickups)

	var intakes int64
	db.Model(&domain.Intake{}).Count(&intakes)
	assert.Equal(t, int64(0), intakes)
}

func mustCloudEvent(t *testing.T, topic string, subject string, payload interface{}, metadata events.Metadata) []byte {
	t.Helper()
	cloudEvent, err := events.NewCloudEvent(context.Background(), topic, "/tests/warehouse-flow", subject, payload, metadata)
	assert.NoError(t, err)
	body, err := json.Marshal(cloudEvent)
	assert.NoError(t, err)
	return body
}
