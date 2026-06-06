package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"RuntimeRoasters/apps/trace-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeRetailClient struct {
	item  *RetailSoldItem
	items []RetailSoldItem
	err   error
}

func (f fakeRetailClient) GetSoldItem(context.Context, string) (*RetailSoldItem, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.item == nil {
		return nil, ErrSoldItemNotFound
	}
	return f.item, nil
}
func (f fakeRetailClient) ListSoldItems(context.Context, int) ([]RetailSoldItem, error) {
	return f.items, f.err
}

func TestGetPublicTraceComposesOnlyPublicFields(t *testing.T) {
	db := setupTraceTestDB(t)
	now := time.Now().UTC()
	require.NoError(t, db.Create(&domain.TraceEvent{ID: "event-1", MessageID: "message-1", Topic: "farm.harvest.created", BatchID: "batch-1", HarvestID: "harvest-1", TraceID: "trace-1", SourceService: "/services/farm-service", Visibility: "public", DisplayPayload: `{}`, Payload: `{"data":{"harvest_id":"harvest-1","farm_id":"farm-1"}}`, OccurredAt: now.Add(-time.Hour)}).Error)
	service := NewService(db, nil)
	service.SetRetailClient(fakeRetailClient{item: &RetailSoldItem{ID: "sale-item-1", StoreID: "store-1", StoreName: "Hoan Kiem", StoreCity: "Hanoi", StoreAddress: "2 Ly Thai To", MenuItemID: "menu-item-1", InventoryLotID: "lot-1", ProductID: "RR-CUP-1", TraceCode: "RR-CUP-1", ProductName: "Espresso", SKU: "BEAN-ARABICA", Size: "S", InvoiceNo: "INV-1", SourceBatchID: "batch-1", SourceHarvestID: "harvest-1", SourceWarehouseID: "warehouse-1", SoldAt: now}})

	doc, err := service.GetPublicTrace(context.Background(), "RR-CUP-1")
	require.NoError(t, err)
	assert.True(t, doc.PublicSafe)
	assert.Equal(t, "RR-CUP-1", doc.ProductID)
	assert.Equal(t, "harvest-1", doc.HarvestID)
	assert.Len(t, doc.Milestones, 2)
	assert.Equal(t, []string{"farm-service", "retail-service"}, doc.ServiceParticipation)
	assert.Equal(t, []string{"trace-1"}, doc.TraceIDs)
	_, leaksPayload := doc.Product["payload"]
	assert.False(t, leaksPayload)
}

func TestPublicTraceDependencyAndNotFoundErrors(t *testing.T) {
	db := setupTraceTestDB(t)
	service := NewService(db, nil)
	service.SetRetailClient(fakeRetailClient{err: ErrRetailUnavailable})
	_, err := service.GetPublicTrace(context.Background(), "RR-CUP-1")
	assert.True(t, errors.Is(err, ErrRetailUnavailable))

	service.SetRetailClient(fakeRetailClient{})
	_, err = service.GetPublicTrace(context.Background(), "missing")
	assert.True(t, errors.Is(err, ErrSoldItemNotFound))
}

func TestListPublicTraceItemsContainsOnlySoldCupProjection(t *testing.T) {
	db := setupTraceTestDB(t)
	service := NewService(db, nil)
	service.SetRetailClient(fakeRetailClient{items: []RetailSoldItem{{ProductID: "cup-1", TraceCode: "cup-1", ProductName: "Latte", StoreName: "District 1", InvoiceNo: "INV-1", PublicURL: "/trace/cup-1"}}})
	items, err := service.ListPublicTraceItems(context.Background(), 25)
	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "cup-1", items[0].ProductID)
}

func setupTraceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.TraceEvent{}, &domain.TraceDocument{}))
	return db
}
