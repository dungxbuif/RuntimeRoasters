package app

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	svcconfig "RuntimeRoasters/apps/retail-service/config"
	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/apps/retail-service/internal/usecase"
	"RuntimeRoasters/pkg/events"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPublicSaleAndInternalSoldItemRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&domain.Store{}, &domain.Menu{}, &domain.MenuItem{}, &domain.Order{},
		&domain.InventoryLot{}, &domain.Sale{}, &domain.SaleItem{},
		&domain.StockMovement{}, &domain.StoreMenuInventory{},
		&domain.OutboxEvent{}, &domain.InboxEvent{},
	))
	service := usecase.NewService(db, nil, events.TopicRetailOrderCreated)
	_, err = service.SeedData(context.Background(), false, map[string]string{})
	require.NoError(t, err)

	application := &App{Cfg: &svcconfig.Config{InternalSecret: "test-secret"}, Service: service}
	router := gin.New()
	application.routes(router)

	menuResponse := httptest.NewRecorder()
	router.ServeHTTP(menuResponse, httptest.NewRequest(http.MethodGet, "/v1/public/retail/stores/11111111-1111-1111-1111-111111111101/menu-items", nil))
	assert.Equal(t, http.StatusOK, menuResponse.Code)
	assert.Contains(t, menuResponse.Body.String(), "available_units")

	requestBody, err := json.Marshal(usecase.CreateSaleRequest{StoreID: "11111111-1111-1111-1111-111111111101", Items: []usecase.CreateSaleItemRequest{{MenuItemID: "MI-ESPRESSO-DOUBLE-S", ProductID: "RR-CUP-HTTP-001"}}})
	require.NoError(t, err)
	saleRequest := httptest.NewRequest(http.MethodPost, "/v1/public/retail/sales", bytes.NewReader(requestBody))
	saleRequest.Header.Set("Content-Type", "application/json")
	saleRequest.Header.Set("Idempotency-Key", "http-sale-key-001")
	saleResponse := httptest.NewRecorder()
	router.ServeHTTP(saleResponse, saleRequest)
	assert.Equal(t, http.StatusCreated, saleResponse.Code, saleResponse.Body.String())
	assert.Contains(t, saleResponse.Body.String(), "RR-CUP-HTTP-001")

	unauthorized := httptest.NewRecorder()
	router.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/v1/internal/retail/sold-items/RR-CUP-HTTP-001", nil))
	assert.Equal(t, http.StatusUnauthorized, unauthorized.Code)

	internalRequest := httptest.NewRequest(http.MethodGet, "/v1/internal/retail/sold-items/RR-CUP-HTTP-001", nil)
	internalRequest.Header.Set("X-Internal-Secret", "test-secret")
	internalResponse := httptest.NewRecorder()
	router.ServeHTTP(internalResponse, internalRequest)
	assert.Equal(t, http.StatusOK, internalResponse.Code, internalResponse.Body.String())
	assert.Contains(t, internalResponse.Body.String(), "BATCH-")
}
