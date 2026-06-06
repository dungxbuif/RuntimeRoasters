package usecase

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRetailSeedPostgresIntegration(t *testing.T) {
	dsn := os.Getenv("RETAIL_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("RETAIL_TEST_DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	service := NewService(db, nil, events.TopicRetailOrderCreated)
	users := map[string]string{
		"mgr.hn.hoankiem@runtimeroasters.com": "user-hk",
		"mgr.hn.caugiay@runtimeroasters.com":  "user-cg",
		"mgr.hcm.d1@runtimeroasters.com":      "user-d1",
		"mgr.hcm.d7@runtimeroasters.com":      "user-d7",
		"mgr.dn.haichau@runtimeroasters.com":  "user-hc",
	}

	_, err = service.SeedData(context.Background(), false, users)
	require.NoError(t, err)
	_, err = service.SeedData(context.Background(), false, users)
	require.NoError(t, err)

	status, err := service.GetSystemStatus(context.Background())
	require.NoError(t, err)
	assert.True(t, status.Seeded)
	assert.Equal(t, int64(5), status.RecordCounts["stores"])
	assert.Equal(t, int64(42), status.RecordCounts["menu_items"])
	assert.Equal(t, int64(10), status.RecordCounts["inventory_lots"])
	assert.Equal(t, int64(25), status.RecordCounts["sales"])
	assert.Equal(t, int64(25), status.RecordCounts["sale_items"])
	assert.Equal(t, int64(35), status.RecordCounts["stock_movements"])
	assert.Equal(t, int64(210), status.RecordCounts["store_menu_inventories"])

	var mismatched int64
	ledgerQuery := "SELECT COUNT(*) FROM inventory_lots l " +
		"LEFT JOIN (SELECT inventory_lot_id, SUM(quantity_delta) AS balance " +
		"FROM stock_movements GROUP BY inventory_lot_id) m ON m.inventory_lot_id = l.id " +
		"WHERE l.available_quantity <> COALESCE(m.balance, 0)"
	require.NoError(t, db.Raw(ledgerQuery).Scan(&mismatched).Error)
	assert.Zero(t, mismatched)

	var invalidTraceCodes int64
	require.NoError(t, db.Model(&domain.SaleItem{}).
		Where("trace_code <> product_id OR inventory_lot_id = ''").
		Count(&invalidTraceCodes).Error)
	assert.Zero(t, invalidTraceCodes)

	t.Run("concurrent sale does not oversell", func(t *testing.T) {
		now := time.Now().UTC()
		storeID := uuid.NewString()
		menuID := "MENU-PG-" + uuid.NewString()
		menuItemID := "MI-PG-" + uuid.NewString()
		lotID := "LOT-PG-" + uuid.NewString()
		sku := "SKU-PG-" + uuid.NewString()
		require.NoError(t, db.Create(&domain.Store{ID: storeID, Code: "STORE-PG-" + uuid.NewString()[:8], Name: "Postgres Concurrency Store", City: "Test", Address: "Test", Status: domain.StoreStatusActive}).Error)
		require.NoError(t, db.Create(&domain.Menu{ID: menuID, Code: menuID, Name: "Postgres Test Menu", Status: domain.MenuStatusActive}).Error)
		require.NoError(t, db.Create(&domain.MenuItem{ID: menuItemID, MenuID: menuID, Code: menuItemID, ProductGroupCode: "PG-TEST", CategoryCode: "TEST", Size: "S", Name: "Concurrency Cup", StockSKU: sku, CoffeeType: "ARABICA", Price: 10, ConsumptionQuantity: 18, ConsumptionUnit: "GRAM", Active: true}).Error)
		require.NoError(t, db.Create(&domain.InventoryLot{ID: lotID, StoreID: storeID, StockSKU: sku, SourceBatchID: "BATCH-PG", SourceHarvestID: "HARVEST-PG", SourceWarehouseID: "WAREHOUSE-PG", ReceivedQuantity: 18, AvailableQuantity: 18, Unit: "GRAM", Status: domain.InventoryLotStatusAvailable, ReceivedAt: now}).Error)
		require.NoError(t, db.Create(&domain.StockMovement{ID: uuid.NewString(), StoreID: storeID, InventoryLotID: lotID, StockSKU: sku, MovementType: domain.StockMovementReceived, QuantityDelta: 18, Unit: "GRAM", ReferenceType: domain.StockReferenceDemo, ReferenceID: "PG-CONCURRENCY", OccurredAt: now}).Error)
		require.NoError(t, db.Create(&domain.StoreMenuInventory{StoreID: storeID, MenuItemID: menuItemID, AvailableUnits: 1, ActiveLotCount: 1, Version: 1, UpdatedAt: now}).Error)

		start := make(chan struct{})
		errs := make(chan error, 2)
		var wg sync.WaitGroup
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				<-start
				_, err := service.CreateSale(context.Background(), CreateSaleRequest{StoreID: storeID, Items: []CreateSaleItemRequest{{MenuItemID: menuItemID, ProductID: uuid.NewString()}}}, fmt.Sprintf("pg-concurrent-%d-%s", index, uuid.NewString()))
				errs <- err
			}(i)
		}
		close(start)
		wg.Wait()
		close(errs)
		var successes, conflicts int
		for err := range errs {
			if err == nil {
				successes++
			} else if ErrorKindOf(err) == ErrorConflict {
				conflicts++
			} else {
				require.NoError(t, err)
			}
		}
		assert.Equal(t, 1, successes)
		assert.Equal(t, 1, conflicts)
		var lot domain.InventoryLot
		require.NoError(t, db.Where("id = ?", lotID).Take(&lot).Error)
		assert.Zero(t, lot.AvailableQuantity)
	})
}
