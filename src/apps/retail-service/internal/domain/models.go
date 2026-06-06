package domain

import (
	"time"
)

type StoreStatus string
type OrderStatus string
type OutboxStatus string
type MenuStatus string
type InventoryLotStatus string
type SaleStatus string
type StockMovementType string
type StockReferenceType string

const (
	StoreStatusActive   StoreStatus = "ACTIVE"
	StoreStatusInactive StoreStatus = "INACTIVE"

	OrderStatusPending           OrderStatus = "PENDING"
	OrderStatusPaymentPending    OrderStatus = "PAYMENT_PENDING"
	OrderStatusPaymentCompleted  OrderStatus = "PAYMENT_COMPLETED"
	OrderStatusReserved          OrderStatus = "RESERVED"
	OrderStatusDispatchRequested OrderStatus = "DISPATCH_REQUESTED"
	OrderStatusPreparing         OrderStatus = "PREPARING"
	OrderStatusShipping          OrderStatus = "SHIPPING"
	OrderStatusDelivered         OrderStatus = "DELIVERED"
	OrderStatusCompleted         OrderStatus = "COMPLETED"
	OrderStatusRejected          OrderStatus = "REJECTED"
	OrderStatusFailed            OrderStatus = "FAILED"

	OutboxStatusPending   OutboxStatus = "PENDING"
	OutboxStatusCompleted OutboxStatus = "COMPLETED"
	OutboxStatusFailed    OutboxStatus = "FAILED"

	MenuStatusDraft    MenuStatus = "DRAFT"
	MenuStatusActive   MenuStatus = "ACTIVE"
	MenuStatusInactive MenuStatus = "INACTIVE"

	InventoryLotStatusAvailable InventoryLotStatus = "AVAILABLE"
	InventoryLotStatusDepleted  InventoryLotStatus = "DEPLETED"
	InventoryLotStatusBlocked   InventoryLotStatus = "BLOCKED"

	SaleStatusCompleted SaleStatus = "COMPLETED"
	SaleStatusVoided    SaleStatus = "VOIDED"

	StockMovementReceived   StockMovementType = "RECEIVED"
	StockMovementSold       StockMovementType = "SOLD"
	StockMovementAdjustment StockMovementType = "ADJUSTMENT"

	StockReferenceSupplyOrder StockReferenceType = "SUPPLY_ORDER"
	StockReferenceSaleItem    StockReferenceType = "SALE_ITEM"
	StockReferenceDemo        StockReferenceType = "DEMO_SCENARIO"
	StockReferenceAdjustment  StockReferenceType = "ADJUSTMENT"
)

type Store struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	Code         string      `gorm:"size:32;uniqueIndex;not null" json:"code"`
	Name         string      `gorm:"size:120;not null" json:"name"`
	City         string      `gorm:"size:80;not null;index" json:"city"`
	Address      string      `gorm:"size:255;not null" json:"address"`
	ManagerID    string      `gorm:"type:varchar(64);index" json:"manager_id"`
	ManagerEmail string      `gorm:"size:160;index" json:"manager_email"`
	Status       StoreStatus `gorm:"size:20;not null;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type Menu struct {
	ID            string     `gorm:"size:64;primaryKey" json:"id"`
	Code          string     `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name          string     `gorm:"size:120;not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description,omitempty"`
	Status        MenuStatus `gorm:"size:20;not null;index" json:"status"`
	EffectiveFrom *time.Time `json:"effective_from,omitempty"`
	EffectiveTo   *time.Time `json:"effective_to,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type MenuItem struct {
	ID                  string  `gorm:"size:64;primaryKey" json:"id"`
	MenuID              string  `gorm:"size:64;not null;index:idx_menu_items_menu_active_order,priority:1" json:"menu_id"`
	Menu                Menu    `gorm:"foreignKey:MenuID" json:"-"`
	Code                string  `gorm:"size:64;uniqueIndex;not null" json:"code"`
	ProductGroupCode    string  `gorm:"size:64;not null;index" json:"product_group_code"`
	CategoryCode        string  `gorm:"size:32;not null;index" json:"category_code"`
	Size                string  `gorm:"size:8;not null" json:"size"`
	Name                string  `gorm:"size:120;not null" json:"name"`
	Description         string  `gorm:"type:text" json:"description,omitempty"`
	StockSKU            string  `gorm:"size:80;not null;index" json:"stock_sku"`
	CoffeeType          string  `gorm:"size:30;not null" json:"coffee_type"`
	Price               float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	ConsumptionQuantity float64 `gorm:"type:decimal(12,3);not null" json:"consumption_quantity"`
	ConsumptionUnit     string  `gorm:"size:16;not null" json:"consumption_unit"`
	ImageURL            string  `gorm:"type:text" json:"image_url,omitempty"`
	Active              bool    `gorm:"not null;index:idx_menu_items_menu_active_order,priority:2" json:"active"`
	DisplayOrder        int     `gorm:"not null;index:idx_menu_items_menu_active_order,priority:3" json:"display_order"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type InventoryLot struct {
	ID                string             `gorm:"size:64;primaryKey" json:"id"`
	StoreID           string             `gorm:"type:uuid;not null;index:idx_inventory_lots_lookup,priority:1" json:"store_id"`
	Store             Store              `gorm:"foreignKey:StoreID" json:"-"`
	StockSKU          string             `gorm:"size:80;not null;index:idx_inventory_lots_lookup,priority:2" json:"stock_sku"`
	SourceBatchID     string             `gorm:"size:128;not null;index" json:"source_batch_id"`
	SourceHarvestID   string             `gorm:"size:128;not null;index" json:"source_harvest_id"`
	SourceWarehouseID string             `gorm:"size:80;not null;index" json:"source_warehouse_id"`
	SourceOrderID     *string            `gorm:"type:uuid;index" json:"source_order_id,omitempty"`
	ReceivedQuantity  float64            `gorm:"type:decimal(12,3);not null" json:"received_quantity"`
	AvailableQuantity float64            `gorm:"type:decimal(12,3);not null" json:"available_quantity"`
	Unit              string             `gorm:"size:16;not null" json:"unit"`
	Status            InventoryLotStatus `gorm:"size:20;not null;index:idx_inventory_lots_lookup,priority:3" json:"status"`
	ReceivedAt        time.Time          `gorm:"not null;index:idx_inventory_lots_lookup,priority:4" json:"received_at"`
	ExpiresAt         *time.Time         `gorm:"index" json:"expires_at,omitempty"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Sale struct {
	ID             string     `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID        string     `gorm:"type:uuid;not null;index:idx_sales_store_sold,priority:1" json:"store_id"`
	Store          Store      `gorm:"foreignKey:StoreID" json:"-"`
	InvoiceNo      string     `gorm:"size:64;uniqueIndex;not null" json:"invoice_no"`
	Status         SaleStatus `gorm:"size:20;not null;index" json:"status"`
	Subtotal       float64    `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	TotalAmount    float64    `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	IdempotencyKey string     `gorm:"size:128;uniqueIndex;not null" json:"idempotency_key"`
	SoldAt         time.Time  `gorm:"not null;index:idx_sales_store_sold,priority:2,sort:desc" json:"sold_at"`
	CreatedAt      time.Time
}

type SaleItem struct {
	ID               string    `gorm:"type:uuid;primaryKey" json:"id"`
	SaleID           string    `gorm:"type:uuid;not null;index" json:"sale_id"`
	Sale             Sale      `gorm:"foreignKey:SaleID" json:"-"`
	StoreID          string    `gorm:"type:uuid;not null;index:idx_sale_items_store_sold,priority:1" json:"store_id"`
	MenuItemID       string    `gorm:"size:64;not null;index" json:"menu_item_id"`
	InventoryLotID   string    `gorm:"size:64;not null;index" json:"inventory_lot_id"`
	ProductID        string    `gorm:"size:128;uniqueIndex;not null" json:"product_id"`
	TraceCode        string    `gorm:"size:128;uniqueIndex;not null" json:"trace_code"`
	ProductName      string    `gorm:"size:120;not null" json:"product_name"`
	SKU              string    `gorm:"size:80;not null" json:"sku"`
	Size             string    `gorm:"size:8;not null" json:"size"`
	UnitPrice        float64   `gorm:"type:decimal(12,2);not null" json:"unit_price"`
	ConsumedQuantity float64   `gorm:"type:decimal(12,3);not null" json:"consumed_quantity"`
	ConsumedUnit     string    `gorm:"size:16;not null" json:"consumed_unit"`
	PublicURL        string    `gorm:"type:text;not null" json:"public_url"`
	SoldAt           time.Time `gorm:"not null;index:idx_sale_items_store_sold,priority:2,sort:desc" json:"sold_at"`
	CreatedAt        time.Time
	MenuItem         MenuItem     `gorm:"foreignKey:MenuItemID" json:"-"`
	InventoryLot     InventoryLot `gorm:"foreignKey:InventoryLotID" json:"-"`
}

type StockMovement struct {
	ID             string             `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID        string             `gorm:"type:uuid;not null;index" json:"store_id"`
	InventoryLotID string             `gorm:"size:64;not null;uniqueIndex:uq_stock_movement_reference,priority:4;index:idx_stock_movements_lot_time,priority:1" json:"inventory_lot_id"`
	StockSKU       string             `gorm:"size:80;not null;index" json:"stock_sku"`
	MovementType   StockMovementType  `gorm:"size:20;not null;uniqueIndex:uq_stock_movement_reference,priority:1" json:"movement_type"`
	QuantityDelta  float64            `gorm:"type:decimal(12,3);not null" json:"quantity_delta"`
	Unit           string             `gorm:"size:16;not null" json:"unit"`
	ReferenceType  StockReferenceType `gorm:"size:30;not null;uniqueIndex:uq_stock_movement_reference,priority:2" json:"reference_type"`
	ReferenceID    string             `gorm:"size:128;not null;uniqueIndex:uq_stock_movement_reference,priority:3" json:"reference_id"`
	OccurredAt     time.Time          `gorm:"not null;index:idx_stock_movements_lot_time,priority:2" json:"occurred_at"`
	CreatedAt      time.Time
}

type StoreMenuInventory struct {
	StoreID        string    `gorm:"type:uuid;primaryKey" json:"store_id"`
	MenuItemID     string    `gorm:"size:64;primaryKey" json:"menu_item_id"`
	AvailableUnits int64     `gorm:"not null" json:"available_units"`
	ActiveLotCount int64     `gorm:"not null" json:"active_lot_count"`
	Version        int64     `gorm:"not null;default:1" json:"version"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Order struct {
	ID             string      `gorm:"type:uuid;primaryKey" json:"id"`
	StoreID        string      `gorm:"type:uuid;not null;index" json:"store_id"`
	Store          Store       `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Items          string      `gorm:"type:jsonb;not null" json:"items"`
	TotalAmount    float64     `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	Status         OrderStatus `gorm:"size:20;not null;default:'PENDING';index" json:"status"`
	IdempotencyKey string      `gorm:"size:128;uniqueIndex" json:"idempotency_key"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type OutboxEvent struct {
	ID          string       `gorm:"type:uuid;primaryKey"`
	EventType   string       `gorm:"size:120;not null;index"`
	Topic       string       `gorm:"size:160;not null;index"`
	Key         string       `gorm:"size:128;not null;index"`
	Payload     string       `gorm:"type:jsonb;not null"`
	TraceParent string       `gorm:"size:128" json:"traceparent,omitempty"`
	TraceState  string       `gorm:"size:512" json:"tracestate,omitempty"`
	Status      OutboxStatus `gorm:"size:20;not null;default:'PENDING';index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type InboxEvent struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	MessageID   string    `gorm:"uniqueIndex;not null"`
	EventType   string    `gorm:"size:120;not null"`
	ProcessedAt time.Time `gorm:"not null"`
}
