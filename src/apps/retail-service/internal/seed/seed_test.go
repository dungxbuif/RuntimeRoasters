package seed

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadStores(t *testing.T) {
	stores, err := LoadStores()
	if err != nil {
		t.Fatalf("LoadStores() error = %v", err)
	}
	if len(stores) != 5 {
		t.Fatalf("LoadStores() len = %d, want 5", len(stores))
	}
	if stores[0].ID == "" || stores[0].ManagerEmail == "" {
		t.Fatalf("LoadStores() returned incomplete first store: %+v", stores[0])
	}
	assert.Equal(t, "HK", stores[0].Code)
}

func TestLoadMenuMatchesSampleMenu(t *testing.T) {
	menu, items, err := LoadMenu()
	require.NoError(t, err)
	assert.Equal(t, DefaultMenuID, menu.ID)
	assert.Len(t, items, 42)

	byID := map[string]float64{}
	for _, item := range items {
		byID[item.ID] = item.Price
		assert.Positive(t, item.ConsumptionQuantity)
		assert.Equal(t, UnitGram, item.ConsumptionUnit)
	}
	assert.Equal(t, float64(29000), byID["MI-PHIN-SUA-DA-S"])
	assert.Equal(t, float64(85000), byID["MI-CARAMEL-MACCHIATO-L"])
	_, espressoLargeExists := byID["MI-ESPRESSO-DOUBLE-L"]
	assert.False(t, espressoLargeExists)
}

func TestBuildDemoDataIsDeterministicAndTraceReady(t *testing.T) {
	stores, err := LoadStores()
	require.NoError(t, err)
	_, items, err := LoadMenu()
	require.NoError(t, err)

	first, err := BuildDemoData(stores, items)
	require.NoError(t, err)
	second, err := BuildDemoData(stores, items)
	require.NoError(t, err)

	assert.Equal(t, first, second)
	assert.Len(t, first.Lots, 10)
	assert.Len(t, first.Sales, 25)
	assert.Len(t, first.SaleItems, 25)
	assert.Len(t, first.Movements, 35)
	for _, item := range first.SaleItems {
		assert.NotEmpty(t, item.ProductID)
		assert.Equal(t, item.ProductID, item.TraceCode)
		assert.NotEmpty(t, item.InventoryLotID)
	}
}
