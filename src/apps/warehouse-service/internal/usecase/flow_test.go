package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
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
	db.AutoMigrate(&domain.Intake{}, &domain.ProductionBatch{}, &domain.RoastRun{}, &domain.Inventory{}, &domain.InboxEvent{})

	ctx := context.Background()
	logger.InitLogger("development", "debug")

	// 1. Simulate Harvest Event (Intake)
	intakeUC := NewIntakeUseCase(db)
	payload := `{"harvest_id": "HV-001", "coffee_type": "ARABICA", "origin_code": "SL", "quantity": 100.0}`
	err := intakeUC.ProcessHarvestEvent(ctx, "msg-1", []byte(payload))
	assert.NoError(t, err)

	// Verify Intake Created
	var intake domain.Intake
	db.First(&intake)
	assert.Equal(t, domain.IntakeStatusUnassigned, intake.Status)
	assert.Equal(t, 100.0, intake.Quantity)

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
	mockProducer.On("Publish", mock.Anything, "warehouse.stock.updated", mock.Anything, mock.Anything).Return(nil)

	invUC := NewInventoryUseCase(db, mockProducer, "warehouse.stock.updated")
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
