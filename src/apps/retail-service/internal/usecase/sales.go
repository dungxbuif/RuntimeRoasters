package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var productIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)

type ErrorKind string

const (
	ErrorValidation ErrorKind = "validation"
	ErrorNotFound   ErrorKind = "not_found"
	ErrorConflict   ErrorKind = "conflict"
)

type ServiceError struct {
	Kind    ErrorKind
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

func ErrorKindOf(err error) ErrorKind {
	var serviceErr *ServiceError
	if errors.As(err, &serviceErr) {
		return serviceErr.Kind
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorNotFound
	}
	return ""
}

type MenuItemView struct {
	ID                  string  `json:"id"`
	MenuID              string  `json:"menu_id"`
	Code                string  `json:"code"`
	ProductGroupCode    string  `json:"product_group_code"`
	CategoryCode        string  `json:"category_code"`
	Size                string  `json:"size"`
	Name                string  `json:"name"`
	Description         string  `json:"description,omitempty"`
	StockSKU            string  `json:"stock_sku"`
	CoffeeType          string  `json:"coffee_type"`
	Price               float64 `json:"price"`
	ConsumptionQuantity float64 `json:"consumption_quantity"`
	ConsumptionUnit     string  `json:"consumption_unit"`
	ImageURL            string  `json:"image_url,omitempty"`
	DisplayOrder        int     `json:"display_order"`
	Available           bool    `json:"available,omitempty"`
	AvailableUnits      int64   `json:"available_units,omitempty"`
}

type MenuView struct {
	ID          string         `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Items       []MenuItemView `json:"items"`
}

type CreateSaleItemRequest struct {
	MenuItemID string `json:"menu_item_id"`
	ProductID  string `json:"product_id"`
}

type CreateSaleRequest struct {
	StoreID string                  `json:"store_id"`
	Items   []CreateSaleItemRequest `json:"items"`
}

type SoldItemView struct {
	domain.SaleItem
	InvoiceNo         string `json:"invoice_no"`
	StoreName         string `json:"store_name"`
	StoreCity         string `json:"store_city"`
	StoreAddress      string `json:"store_address"`
	SourceBatchID     string `json:"source_batch_id"`
	SourceHarvestID   string `json:"source_harvest_id"`
	SourceWarehouseID string `json:"source_warehouse_id"`
}

type SaleView struct {
	domain.Sale
	Items []SoldItemView `json:"items"`
}

func (s *Service) ListPublicStores(ctx context.Context) ([]domain.Store, error) {
	var stores []domain.Store
	err := s.db.WithContext(ctx).Where("status = ?", domain.StoreStatusActive).Order("city, name").Find(&stores).Error
	return stores, err
}

func (s *Service) ListPublicMenu(ctx context.Context) ([]MenuView, error) {
	var menus []domain.Menu
	if err := s.db.WithContext(ctx).Where("status = ?", domain.MenuStatusActive).Order("name").Find(&menus).Error; err != nil {
		return nil, err
	}
	var items []domain.MenuItem
	if err := s.db.WithContext(ctx).Where("active = ?", true).Order("display_order, name, size").Find(&items).Error; err != nil {
		return nil, err
	}
	return buildMenuViews(menus, items, nil), nil
}

func (s *Service) ListPublicStoreMenu(ctx context.Context, storeID string) ([]MenuView, error) {
	var store domain.Store
	if err := s.db.WithContext(ctx).Where("id = ? AND status = ?", storeID, domain.StoreStatusActive).Take(&store).Error; err != nil {
		return nil, err
	}
	var menus []domain.Menu
	if err := s.db.WithContext(ctx).Where("status = ?", domain.MenuStatusActive).Order("name").Find(&menus).Error; err != nil {
		return nil, err
	}
	var items []domain.MenuItem
	if err := s.db.WithContext(ctx).Where("active = ?", true).Order("display_order, name, size").Find(&items).Error; err != nil {
		return nil, err
	}
	var inventory []domain.StoreMenuInventory
	if err := s.db.WithContext(ctx).Where("store_id = ?", storeID).Find(&inventory).Error; err != nil {
		return nil, err
	}
	byItem := make(map[string]domain.StoreMenuInventory, len(inventory))
	for _, row := range inventory {
		byItem[row.MenuItemID] = row
	}
	return buildMenuViews(menus, items, byItem), nil
}

func buildMenuViews(menus []domain.Menu, items []domain.MenuItem, inventory map[string]domain.StoreMenuInventory) []MenuView {
	views := make([]MenuView, 0, len(menus))
	index := make(map[string]int, len(menus))
	for _, menu := range menus {
		index[menu.ID] = len(views)
		views = append(views, MenuView{ID: menu.ID, Code: menu.Code, Name: menu.Name, Description: menu.Description, Items: []MenuItemView{}})
	}
	for _, item := range items {
		menuIndex, ok := index[item.MenuID]
		if !ok {
			continue
		}
		view := MenuItemView{ID: item.ID, MenuID: item.MenuID, Code: item.Code, ProductGroupCode: item.ProductGroupCode, CategoryCode: item.CategoryCode, Size: item.Size, Name: item.Name, Description: item.Description, StockSKU: item.StockSKU, CoffeeType: item.CoffeeType, Price: item.Price, ConsumptionQuantity: item.ConsumptionQuantity, ConsumptionUnit: item.ConsumptionUnit, ImageURL: item.ImageURL, DisplayOrder: item.DisplayOrder}
		if inventory != nil {
			row := inventory[item.ID]
			view.AvailableUnits = row.AvailableUnits
			view.Available = row.AvailableUnits > 0
		}
		views[menuIndex].Items = append(views[menuIndex].Items, view)
	}
	return views
}

func (s *Service) CreateSale(ctx context.Context, req CreateSaleRequest, idempotencyKey string) (*SaleView, error) {
	req.StoreID = strings.TrimSpace(req.StoreID)
	idempotencyKey = strings.TrimSpace(idempotencyKey)
	if req.StoreID == "" {
		return nil, &ServiceError{Kind: ErrorValidation, Message: "store_id is required"}
	}
	if idempotencyKey == "" {
		return nil, &ServiceError{Kind: ErrorValidation, Message: "Idempotency-Key header is required"}
	}
	if len(idempotencyKey) > 128 {
		return nil, &ServiceError{Kind: ErrorValidation, Message: "Idempotency-Key must be at most 128 characters"}
	}
	if len(req.Items) == 0 {
		return nil, &ServiceError{Kind: ErrorValidation, Message: "at least one sale item is required"}
	}
	seenProducts := make(map[string]struct{}, len(req.Items))
	for i := range req.Items {
		req.Items[i].MenuItemID = strings.TrimSpace(req.Items[i].MenuItemID)
		req.Items[i].ProductID = strings.TrimSpace(req.Items[i].ProductID)
		item := req.Items[i]
		if item.MenuItemID == "" || item.ProductID == "" {
			return nil, &ServiceError{Kind: ErrorValidation, Message: fmt.Sprintf("items[%d] requires menu_item_id and product_id", i)}
		}
		if len(item.MenuItemID) > 64 || !productIDPattern.MatchString(item.ProductID) {
			return nil, &ServiceError{Kind: ErrorValidation, Message: fmt.Sprintf("items[%d] has an invalid menu_item_id or product_id", i)}
		}
		if _, exists := seenProducts[item.ProductID]; exists {
			return nil, &ServiceError{Kind: ErrorValidation, Message: "duplicate product_id in request: " + item.ProductID}
		}
		seenProducts[item.ProductID] = struct{}{}
	}

	if existing, err := s.getSaleByIdempotencyKey(ctx, idempotencyKey); err == nil {
		return existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var createdSaleID string
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var store domain.Store
		if err := tx.Where("id = ? AND status = ?", req.StoreID, domain.StoreStatusActive).Take(&store).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &ServiceError{Kind: ErrorNotFound, Message: "active store not found"}
			}
			return err
		}

		menuIDs := make([]string, 0, len(req.Items))
		for _, item := range req.Items {
			menuIDs = append(menuIDs, item.MenuItemID)
		}
		var menuItems []domain.MenuItem
		if err := tx.Joins("JOIN menus ON menus.id = menu_items.menu_id").Where("menu_items.id IN ? AND menu_items.active = ? AND menus.status = ?", menuIDs, true, domain.MenuStatusActive).Find(&menuItems).Error; err != nil {
			return err
		}
		byID := make(map[string]domain.MenuItem, len(menuItems))
		for _, item := range menuItems {
			byID[item.ID] = item
		}
		for _, item := range req.Items {
			if _, ok := byID[item.MenuItemID]; !ok {
				return &ServiceError{Kind: ErrorNotFound, Message: "active menu item not found: " + item.MenuItemID}
			}
		}

		work := append([]CreateSaleItemRequest(nil), req.Items...)
		sort.SliceStable(work, func(i, j int) bool {
			left, right := byID[work[i].MenuItemID], byID[work[j].MenuItemID]
			if left.StockSKU == right.StockSKU {
				return work[i].ProductID < work[j].ProductID
			}
			return left.StockSKU < right.StockSKU
		})
		now := time.Now().UTC()
		sale := domain.Sale{ID: uuid.NewString(), StoreID: req.StoreID, InvoiceNo: "INV-" + now.Format("20060102-150405") + "-" + strings.ToUpper(uuid.NewString()[:8]), Status: domain.SaleStatusCompleted, IdempotencyKey: idempotencyKey, SoldAt: now}
		if err := tx.Create(&sale).Error; err != nil {
			return mapSaleConflict(err)
		}

		eventItems := make([]events.RetailSaleCompletedItem, 0, len(work))
		affectedSKUs := make(map[string]struct{})
		for _, requested := range work {
			menuItem := byID[requested.MenuItemID]
			var lot domain.InventoryLot
			query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("store_id = ? AND stock_sku = ? AND status = ? AND available_quantity >= ? AND (expires_at IS NULL OR expires_at > ?)", req.StoreID, menuItem.StockSKU, domain.InventoryLotStatusAvailable, menuItem.ConsumptionQuantity, now).Order("received_at, id")
			if err := query.Take(&lot).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &ServiceError{Kind: ErrorConflict, Message: "inventory unavailable for menu item: " + menuItem.ID}
				}
				return err
			}
			itemID := uuid.NewString()
			item := domain.SaleItem{ID: itemID, SaleID: sale.ID, StoreID: req.StoreID, MenuItemID: menuItem.ID, InventoryLotID: lot.ID, ProductID: requested.ProductID, TraceCode: requested.ProductID, ProductName: menuItem.Name, SKU: menuItem.StockSKU, Size: menuItem.Size, UnitPrice: menuItem.Price, ConsumedQuantity: menuItem.ConsumptionQuantity, ConsumedUnit: menuItem.ConsumptionUnit, PublicURL: "/trace/" + requested.ProductID, SoldAt: now}
			if err := tx.Create(&item).Error; err != nil {
				return mapSaleConflict(err)
			}
			result := tx.Model(&domain.InventoryLot{}).Where("id = ? AND available_quantity >= ?", lot.ID, menuItem.ConsumptionQuantity).Updates(map[string]interface{}{"available_quantity": gorm.Expr("available_quantity - ?", menuItem.ConsumptionQuantity), "updated_at": now})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return &ServiceError{Kind: ErrorConflict, Message: "inventory changed while creating sale"}
			}
			if err := tx.Model(&domain.InventoryLot{}).Where("id = ? AND available_quantity = 0", lot.ID).Update("status", domain.InventoryLotStatusDepleted).Error; err != nil {
				return err
			}
			movement := domain.StockMovement{ID: uuid.NewString(), StoreID: req.StoreID, InventoryLotID: lot.ID, StockSKU: menuItem.StockSKU, MovementType: domain.StockMovementSold, QuantityDelta: -menuItem.ConsumptionQuantity, Unit: menuItem.ConsumptionUnit, ReferenceType: domain.StockReferenceSaleItem, ReferenceID: item.ID, OccurredAt: now}
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
			sale.Subtotal += menuItem.Price
			affectedSKUs[menuItem.StockSKU] = struct{}{}
			eventItems = append(eventItems, events.RetailSaleCompletedItem{SaleItemID: item.ID, MenuItemID: item.MenuItemID, ProductID: item.ProductID, SKU: item.SKU, ProductName: item.ProductName, UnitPrice: item.UnitPrice, InventoryLotID: lot.ID, TraceCode: item.TraceCode, PublicURL: item.PublicURL, SourceBatchID: lot.SourceBatchID, SourceHarvestID: lot.SourceHarvestID, SourceWarehouseID: lot.SourceWarehouseID})
		}
		sale.TotalAmount = sale.Subtotal
		if err := tx.Model(&domain.Sale{}).Where("id = ?", sale.ID).Updates(map[string]interface{}{"subtotal": sale.Subtotal, "total_amount": sale.TotalAmount}).Error; err != nil {
			return err
		}
		if err := rebuildStoreMenuInventoryForSKUs(ctx, tx, req.StoreID, affectedSKUs, now); err != nil {
			return err
		}
		if err := createSaleOutbox(ctx, tx, sale, eventItems); err != nil {
			return err
		}
		createdSaleID = sale.ID
		return nil
	})
	if err != nil {
		if ErrorKindOf(err) == ErrorConflict {
			if existing, lookupErr := s.getSaleByIdempotencyKey(ctx, idempotencyKey); lookupErr == nil {
				return existing, nil
			}
		}
		return nil, err
	}
	return s.GetSale(ctx, createdSaleID)
}

func mapSaleConflict(err error) error {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unique") || strings.Contains(message, "duplicate") {
		return &ServiceError{Kind: ErrorConflict, Message: "idempotency key or product_id already exists"}
	}
	return err
}

func createSaleOutbox(ctx context.Context, tx *gorm.DB, sale domain.Sale, items []events.RetailSaleCompletedItem) error {
	eventID := uuid.NewString()
	data := events.RetailSaleCompleted{EventID: eventID, SaleID: sale.ID, InvoiceNo: sale.InvoiceNo, StoreID: sale.StoreID, Items: items, TotalAmount: sale.TotalAmount, SoldAt: sale.SoldAt, OccurredAt: sale.SoldAt}
	cloudEvent, err := events.NewCloudEvent(ctx, events.TopicRetailSaleCompleted, events.SourceRetailService, "sales/"+sale.ID, data, events.Metadata{EventID: eventID, CorrelationID: sale.ID, StoreID: sale.StoreID, OccurredAt: sale.SoldAt})
	if err != nil {
		return err
	}
	payload, err := json.Marshal(cloudEvent)
	if err != nil {
		return err
	}
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return tx.Create(&domain.OutboxEvent{ID: uuid.NewString(), EventType: events.TopicRetailSaleCompleted, Topic: events.TopicRetailSaleCompleted, Key: sale.ID, Payload: string(payload), TraceParent: carrier.Get("traceparent"), TraceState: carrier.Get("tracestate"), Status: domain.OutboxStatusPending}).Error
}

func rebuildStoreMenuInventoryForSKUs(ctx context.Context, tx *gorm.DB, storeID string, skus map[string]struct{}, now time.Time) error {
	if len(skus) == 0 {
		return nil
	}
	skuList := make([]string, 0, len(skus))
	for sku := range skus {
		skuList = append(skuList, sku)
	}
	var items []domain.MenuItem
	if err := tx.WithContext(ctx).Joins("JOIN menus ON menus.id = menu_items.menu_id").Where("menu_items.stock_sku IN ? AND menu_items.active = ? AND menus.status = ?", skuList, true, domain.MenuStatusActive).Find(&items).Error; err != nil {
		return err
	}
	var lots []domain.InventoryLot
	if err := tx.WithContext(ctx).Where("store_id = ? AND stock_sku IN ? AND status = ? AND available_quantity > 0 AND (expires_at IS NULL OR expires_at > ?)", storeID, skuList, domain.InventoryLotStatusAvailable, now).Find(&lots).Error; err != nil {
		return err
	}
	for _, item := range items {
		var units, lotCount int64
		for _, lot := range lots {
			if lot.StockSKU != item.StockSKU {
				continue
			}
			lotUnits := int64(math.Floor(lot.AvailableQuantity / item.ConsumptionQuantity))
			if lotUnits > 0 {
				units += lotUnits
				lotCount++
			}
		}
		row := domain.StoreMenuInventory{StoreID: storeID, MenuItemID: item.ID, AvailableUnits: units, ActiveLotCount: lotCount, Version: 1, UpdatedAt: now}
		if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "store_id"}, {Name: "menu_item_id"}}, DoUpdates: clause.Assignments(map[string]interface{}{"available_units": units, "active_lot_count": lotCount, "version": gorm.Expr("store_menu_inventories.version + 1"), "updated_at": now})}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) getSaleByIdempotencyKey(ctx context.Context, key string) (*SaleView, error) {
	var sale domain.Sale
	if err := s.db.WithContext(ctx).Where("idempotency_key = ?", key).Take(&sale).Error; err != nil {
		return nil, err
	}
	return s.loadSaleView(ctx, sale)
}

func (s *Service) GetSale(ctx context.Context, id string) (*SaleView, error) {
	var sale domain.Sale
	query := s.db.WithContext(ctx).Where("id = ?", id)
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	if !allStores {
		if len(storeIDs) == 0 {
			return nil, gorm.ErrRecordNotFound
		}
		query = query.Where("store_id IN ?", storeIDs)
	}
	if err := query.Take(&sale).Error; err != nil {
		return nil, err
	}
	return s.loadSaleView(ctx, sale)
}

func (s *Service) loadSaleView(ctx context.Context, sale domain.Sale) (*SaleView, error) {
	items, err := s.listSoldItemsQuery(s.db.WithContext(ctx).Where("sale_items.sale_id = ?", sale.ID))
	if err != nil {
		return nil, err
	}
	return &SaleView{Sale: sale, Items: items}, nil
}

func (s *Service) ListSales(ctx context.Context) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := s.db.WithContext(ctx).Order("sold_at DESC")
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	if !allStores {
		if len(storeIDs) == 0 {
			return sales, nil
		}
		query = query.Where("store_id IN ?", storeIDs)
	}
	return sales, query.Find(&sales).Error
}

func (s *Service) ListSoldItems(ctx context.Context, storeID string) ([]SoldItemView, error) {
	if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(storeID) {
		return []SoldItemView{}, nil
	}
	return s.listSoldItemsQuery(s.db.WithContext(ctx).Where("sale_items.store_id = ?", storeID).Order("sale_items.sold_at DESC"))
}

func (s *Service) GetSoldItemByProductID(ctx context.Context, productID string) (*SoldItemView, error) {
	items, err := s.listSoldItemsQuery(s.db.WithContext(ctx).Where("sale_items.product_id = ?", productID).Limit(1))
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &items[0], nil
}

func (s *Service) ListInternalSoldItems(ctx context.Context, limit int) ([]SoldItemView, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}
	return s.listSoldItemsQuery(s.db.WithContext(ctx).Order("sale_items.sold_at DESC").Limit(limit))
}

func (s *Service) listSoldItemsQuery(query *gorm.DB) ([]SoldItemView, error) {
	var items []SoldItemView
	err := query.Model(&domain.SaleItem{}).Select("sale_items.*, sales.invoice_no, stores.name AS store_name, stores.city AS store_city, stores.address AS store_address, inventory_lots.source_batch_id, inventory_lots.source_harvest_id, inventory_lots.source_warehouse_id").Joins("JOIN sales ON sales.id = sale_items.sale_id").Joins("JOIN stores ON stores.id = sale_items.store_id").Joins("JOIN inventory_lots ON inventory_lots.id = sale_items.inventory_lot_id").Scan(&items).Error
	return items, err
}
