package usecase

import (
	"context"
	"os"
	"testing"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/pkg/events"
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
}
