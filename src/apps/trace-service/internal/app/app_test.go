package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"RuntimeRoasters/apps/trace-service/internal/domain"
	"RuntimeRoasters/apps/trace-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type routeRetailClient struct{ found bool }

func (c routeRetailClient) GetSoldItem(context.Context, string) (*usecase.RetailSoldItem, error) {
	if !c.found {
		return nil, usecase.ErrSoldItemNotFound
	}
	return &usecase.RetailSoldItem{ID: "item-1", ProductID: "cup-1", TraceCode: "cup-1", ProductName: "Espresso", StoreID: "store-1", StoreName: "Hoan Kiem", InventoryLotID: "lot-1", SourceBatchID: "batch-1", SourceHarvestID: "harvest-1", SourceWarehouseID: "warehouse-1", SoldAt: time.Now()}, nil
}
func (c routeRetailClient) ListSoldItems(context.Context, int) ([]usecase.RetailSoldItem, error) {
	return []usecase.RetailSoldItem{{ProductID: "cup-1", TraceCode: "cup-1", ProductName: "Espresso"}}, nil
}

func TestPublicTraceRoutesRequireNoAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.TraceEvent{}, &domain.TraceDocument{}))
	service := usecase.NewService(db, nil)
	service.SetRetailClient(routeRetailClient{found: true})
	router := gin.New()
	(&App{Service: service}).routes(router)

	detail := httptest.NewRecorder()
	router.ServeHTTP(detail, httptest.NewRequest(http.MethodGet, "/v1/public/trace/cup-1", nil))
	assert.Equal(t, http.StatusOK, detail.Code)
	assert.Contains(t, detail.Body.String(), `"public_safe":true`)

	list := httptest.NewRecorder()
	router.ServeHTTP(list, httptest.NewRequest(http.MethodGet, "/v1/public/trace-items", nil))
	assert.Equal(t, http.StatusOK, list.Code)
	assert.Contains(t, list.Body.String(), "cup-1")
}

func TestPublicTraceRouteReturnsNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	service := usecase.NewService(db, nil)
	service.SetRetailClient(routeRetailClient{})
	router := gin.New()
	(&App{Service: service}).routes(router)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/v1/public/trace/missing", nil))
	assert.Equal(t, http.StatusNotFound, response.Code)
}
