package events

import "time"

const (
	TopicRetailOrderCreated              = "retail.order.created"
	TopicPaymentIntentCreated            = "payment.intent.created"
	TopicPaymentCompleted                = "payment.completed"
	TopicPaymentFailed                   = "payment.failed"
	TopicPaymentRefunded                 = "payment.refunded"
	TopicWarehouseStockReserved          = "warehouse.stock.reserved"
	TopicWarehouseStockReservationFailed = "warehouse.stock.reservation_failed"
	TopicWarehouseStockUpdated           = "warehouse.stock.updated"
	TopicLogisticsShipmentAssigned       = "logistics.shipment.assigned"
	TopicLogisticsShipmentDelivered      = "logistics.shipment.delivered"
	TopicLogisticsGPSUpdated             = "logistics.gps.updated"
)

type RetailOrderItem struct {
	SKU      string  `json:"sku"`
	Quantity float64 `json:"quantity"`
}

type RetailOrderCreated struct {
	EventID       string            `json:"event_id"`
	OrderID       string            `json:"order_id"`
	StoreID       string            `json:"store_id"`
	Items         []RetailOrderItem `json:"items"`
	TotalAmount   float64           `json:"total_amount"`
	PaymentMethod string            `json:"payment_method,omitempty"`
	OccurredAt    time.Time         `json:"occurred_at"`
}

type PaymentIntentCreated struct {
	EventID     string            `json:"event_id"`
	PaymentID   string            `json:"payment_id"`
	OrderID     string            `json:"order_id"`
	StoreID     string            `json:"store_id"`
	Provider    string            `json:"provider"`
	ProviderRef string            `json:"provider_ref"`
	Items       []RetailOrderItem `json:"items"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	CheckoutURL string            `json:"checkout_url"`
	Simulated   bool              `json:"simulated"`
	OccurredAt  time.Time         `json:"occurred_at"`
}

type PaymentCompleted struct {
	EventID     string            `json:"event_id"`
	PaymentID   string            `json:"payment_id"`
	OrderID     string            `json:"order_id"`
	StoreID     string            `json:"store_id"`
	Provider    string            `json:"provider"`
	ProviderRef string            `json:"provider_ref"`
	Items       []RetailOrderItem `json:"items"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	OccurredAt  time.Time         `json:"occurred_at"`
}

type PaymentFailed struct {
	EventID     string    `json:"event_id"`
	PaymentID   string    `json:"payment_id"`
	OrderID     string    `json:"order_id"`
	StoreID     string    `json:"store_id"`
	Provider    string    `json:"provider"`
	ProviderRef string    `json:"provider_ref"`
	Reason      string    `json:"reason"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type PaymentRefunded struct {
	EventID     string    `json:"event_id"`
	PaymentID   string    `json:"payment_id"`
	OrderID     string    `json:"order_id"`
	StoreID     string    `json:"store_id"`
	Provider    string    `json:"provider"`
	ProviderRef string    `json:"provider_ref"`
	RefundRef   string    `json:"refund_ref"`
	Reason      string    `json:"reason"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehouseStockReserved struct {
	EventID    string    `json:"event_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	SKU        string    `json:"sku"`
	Quantity   float64   `json:"quantity"`
	OccurredAt time.Time `json:"occurred_at"`
}

type WarehouseStockReservationFailed struct {
	EventID    string    `json:"event_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	Reason     string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

type WarehouseStockUpdated struct {
	BatchID   string    `json:"batch_id"`
	SKU       string    `json:"sku"`
	Quantity  float64   `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

type LogisticsShipmentAssigned struct {
	EventID    string    `json:"event_id"`
	ShipmentID string    `json:"shipment_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	DriverID   string    `json:"driver_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type LogisticsShipmentDelivered struct {
	EventID    string    `json:"event_id"`
	ShipmentID string    `json:"shipment_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type LogisticsGPSUpdated struct {
	EventID    string    `json:"event_id"`
	DriverID   string    `json:"driver_id"`
	ShipmentID string    `json:"shipment_id"`
	StoreID    string    `json:"store_id,omitempty"`
	Latitude   float64   `json:"lat"`
	Longitude  float64   `json:"long"`
	OccurredAt time.Time `json:"occurred_at"`
}
