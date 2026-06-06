package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	svcconfig "RuntimeRoasters/apps/retail-service/config"
	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/apps/retail-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/base/security"
	"RuntimeRoasters/pkg/events"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestStoreManagerScoping(t *testing.T) {
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
	
	// Seed two stores and some sales
	storeA := "store-a"
	storeB := "store-b"
	require.NoError(t, db.Create(&domain.Store{ID: storeA, Code: "CODE-A", Name: "Store A", City: "HN", Status: domain.StoreStatusActive}).Error)
	require.NoError(t, db.Create(&domain.Store{ID: storeB, Code: "CODE-B", Name: "Store B", City: "HCM", Status: domain.StoreStatusActive}).Error)
	
	saleA := uuid.NewString()
	saleB := uuid.NewString()
	require.NoError(t, db.Create(&domain.Sale{ID: saleA, StoreID: storeA, InvoiceNo: "INV-A", Status: domain.SaleStatusCompleted, IdempotencyKey: "idemp-a"}).Error)
	require.NoError(t, db.Create(&domain.Sale{ID: saleB, StoreID: storeB, InvoiceNo: "INV-B", Status: domain.SaleStatusCompleted, IdempotencyKey: "idemp-b"}).Error)

	application := &App{
		Cfg: &svcconfig.Config{InternalSecret: "test-secret"}, 
		Service: service,
		Guards: &security.HTTPGuards{
			Authn: func(c *gin.Context) {
				// Mock Authn: Get claims from header if provided (for testing)
				role := c.GetHeader("X-Test-Role")
				storeID := c.GetHeader("X-Test-StoreID")
				if role != "" {
					claims := identity.Claims{
						Subject: "test-user",
						Role:    role,
					}
					if storeID != "" {
						claims.StoreIDs = []string{storeID}
					}
					// Use InjectContext to set claims in request context, then put it back into Gin
					ctx := identity.InjectContext(c.Request.Context(), claims)
					c.Request = c.Request.WithContext(ctx)
				}
				c.Next()
			},
			Authz: func(c *gin.Context) {
				// Mock Authz: Always allow for now, scoping is handled in usecase
				c.Next()
			},
		},
	}
	router := gin.New()
	application.routes(router)

	t.Run("Manager of Store A can see Store A sales", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/retail/sales", nil)
		req.Header.Set("X-Test-Role", identity.RoleStoreMgr)
		req.Header.Set("X-Test-StoreID", storeA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string][]domain.Sale
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp["sales"], 1)
		assert.Equal(t, storeA, resp["sales"][0].StoreID)
	})

	t.Run("Manager of Store A cannot see Store B sales via List", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/retail/sales", nil)
		req.Header.Set("X-Test-Role", identity.RoleStoreMgr)
		req.Header.Set("X-Test-StoreID", storeA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string][]domain.Sale
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		for _, s := range resp["sales"] {
			assert.NotEqual(t, storeB, s.StoreID)
		}
	})

	t.Run("Manager of Store A is blocked from Store B specific sale", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/retail/sales/"+saleB, nil)
		req.Header.Set("X-Test-Role", identity.RoleStoreMgr)
		req.Header.Set("X-Test-StoreID", storeA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		// Should return 404 because GetSale uses store-scope filter
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Manager of Store A is blocked from Store B sold items", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/retail/stores/"+storeB+"/sold-items", nil)
		req.Header.Set("X-Test-Role", identity.RoleStoreMgr)
		req.Header.Set("X-Test-StoreID", storeA)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string][]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Len(t, resp["sold_items"], 0)
	})
}
