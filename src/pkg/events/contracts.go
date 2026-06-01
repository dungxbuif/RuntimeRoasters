package events

import "time"

const (
	TopicFarmHarvestCreated                = "farm.harvest.created"
	TopicRetailOrderCreated                = "retail.order.created"
	TopicPaymentIntentCreated              = "payment.intent.created"
	TopicPaymentCompleted                  = "payment.completed"
	TopicPaymentSimulatedCompleted         = "payment.simulated_completed"
	TopicPaymentFailed                     = "payment.failed"
	TopicPaymentRefunded                   = "payment.refunded"
	TopicWarehouseStockReserved            = "warehouse.stock.reserved"
	TopicWarehouseStockReservationFailed   = "warehouse.stock.reservation_failed"
	TopicWarehouseInventoryUpdated         = "warehouse.inventory.updated"
	TopicWarehousePickupRequested          = "warehouse.pickup.requested"
	TopicWarehousePickupReceived           = "warehouse.pickup.received"
	TopicWarehouseDispatchRequested        = "warehouse.dispatch.requested"
	TopicWarehouseIntakeCreated            = "warehouse.intake.created"
	TopicLogisticsDeliveryAssigned         = "logistics.delivery.assigned"
	TopicLogisticsDeliveryDeparted         = "logistics.delivery.departed"
	TopicLogisticsDeliveryArrivedAtStore   = "logistics.delivery.arrived_at_store"
	TopicLogisticsDeliveryDriverConfirmed  = "logistics.delivery.driver_confirmed"
	TopicLogisticsDeliveryCompleted        = "logistics.delivery.completed"
	TopicLogisticsPickupAssigned           = "logistics.pickup.assigned"
	TopicLogisticsPickupDeparted           = "logistics.pickup.departed"
	TopicLogisticsPickupArrivedAtFarm      = "logistics.pickup.arrived_at_farm"
	TopicLogisticsPickupLoadingConfirmed   = "logistics.pickup.loading_confirmed"
	TopicLogisticsPickupReturnStarted      = "logistics.pickup.return_started"
	TopicLogisticsPickupArrivedAtWarehouse = "logistics.pickup.arrived_at_warehouse"
	TopicLogisticsPickupCompleted          = "logistics.pickup.completed"
	TopicLogisticsDriverReturnStarted      = "logistics.driver.return_started"
	TopicLogisticsDriverReturnCompleted    = "logistics.driver.return_completed"
	TopicLogisticsDriverReturnedToBase     = "logistics.driver.returned_to_base"
	TopicLogisticsGPSUpdated               = "logistics.gps.updated"
	TopicLogisticsShipmentStatusChanged    = "logistics.shipment.status_changed"
	TopicNotificationCreated               = "notification.created"
	TopicNotificationAcknowledged          = "notification.acknowledged"
	TopicSocketBroadcastRequested          = "socket.broadcast.requested"
)

var TraceableTopics = []string{
	TopicFarmHarvestCreated,
	TopicRetailOrderCreated,
	TopicPaymentIntentCreated,
	TopicPaymentCompleted,
	TopicPaymentSimulatedCompleted,
	TopicPaymentFailed,
	TopicPaymentRefunded,
	TopicWarehouseStockReserved,
	TopicWarehouseStockReservationFailed,
	TopicWarehouseInventoryUpdated,
	TopicWarehousePickupRequested,
	TopicWarehousePickupReceived,
	TopicWarehouseDispatchRequested,
	TopicWarehouseIntakeCreated,
	TopicLogisticsDeliveryAssigned,
	TopicLogisticsDeliveryDeparted,
	TopicLogisticsDeliveryArrivedAtStore,
	TopicLogisticsDeliveryDriverConfirmed,
	TopicLogisticsDeliveryCompleted,
	TopicLogisticsPickupAssigned,
	TopicLogisticsPickupDeparted,
	TopicLogisticsPickupArrivedAtFarm,
	TopicLogisticsPickupLoadingConfirmed,
	TopicLogisticsPickupReturnStarted,
	TopicLogisticsPickupArrivedAtWarehouse,
	TopicLogisticsPickupCompleted,
	TopicLogisticsDriverReturnStarted,
	TopicLogisticsDriverReturnCompleted,
	TopicLogisticsDriverReturnedToBase,
	TopicLogisticsGPSUpdated,
	TopicLogisticsShipmentStatusChanged,
	TopicNotificationCreated,
	TopicNotificationAcknowledged,
	TopicSocketBroadcastRequested,
}

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
	EventID     string    `json:"event_id"`
	OrderID     string    `json:"order_id"`
	StoreID     string    `json:"store_id"`
	WarehouseID string    `json:"warehouse_id"`
	SKU         string    `json:"sku"`
	Quantity    float64   `json:"quantity"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehouseStockReservationFailed struct {
	EventID     string    `json:"event_id"`
	OrderID     string    `json:"order_id"`
	StoreID     string    `json:"store_id"`
	WarehouseID string    `json:"warehouse_id"`
	Reason      string    `json:"reason"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehouseStockUpdated struct {
	BatchID   string    `json:"batch_id"`
	SKU       string    `json:"sku"`
	Quantity  float64   `json:"quantity"`
	Timestamp time.Time `json:"timestamp"`
}

type WarehousePickupRequested struct {
	EventID     string    `json:"event_id"`
	PickupID    string    `json:"pickup_id"`
	HarvestID   string    `json:"harvest_id"`
	FarmID      string    `json:"farm_id"`
	WarehouseID string    `json:"warehouse_id"`
	Quantity    float64   `json:"quantity"`
	Status      string    `json:"status,omitempty"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehousePickupReceived struct {
	EventID     string    `json:"event_id"`
	PickupID    string    `json:"pickup_id"`
	HarvestID   string    `json:"harvest_id"`
	FarmID      string    `json:"farm_id"`
	WarehouseID string    `json:"warehouse_id"`
	DriverID    string    `json:"driver_id"`
	VehicleID   string    `json:"vehicle_id"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehouseDispatchRequested struct {
	EventID     string    `json:"event_id"`
	DispatchID  string    `json:"dispatch_id"`
	OrderID     string    `json:"order_id,omitempty"`
	StoreID     string    `json:"store_id,omitempty"`
	WarehouseID string    `json:"warehouse_id"`
	SKU         string    `json:"sku,omitempty"`
	Quantity    float64   `json:"quantity,omitempty"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type WarehouseIntakeCreated struct {
	EventID     string    `json:"event_id"`
	IntakeID    string    `json:"intake_id"`
	HarvestID   string    `json:"harvest_id"`
	BatchID     string    `json:"batch_id,omitempty"`
	FarmID      string    `json:"farm_id"`
	WarehouseID string    `json:"warehouse_id"`
	Quantity    float64   `json:"quantity"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type LogisticsShipmentAssigned struct {
	EventID    string    `json:"event_id"`
	ShipmentID string    `json:"shipment_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	DriverID   string    `json:"driver_id"`
	VehicleID  string    `json:"vehicle_id,omitempty"`
	OccurredAt time.Time `json:"occurred_at"`
}

type LogisticsShipmentDelivered struct {
	EventID    string    `json:"event_id"`
	ShipmentID string    `json:"shipment_id"`
	OrderID    string    `json:"order_id"`
	StoreID    string    `json:"store_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type LogisticsPickupStatusChanged struct {
	EventID     string    `json:"event_id"`
	ShipmentID  string    `json:"shipment_id"`
	HarvestID   string    `json:"harvest_id"`
	FarmID      string    `json:"farm_id"`
	WarehouseID string    `json:"warehouse_id"`
	DriverID    string    `json:"driver_id"`
	VehicleID   string    `json:"vehicle_id"`
	Status      string    `json:"status"`
	OccurredAt  time.Time `json:"occurred_at"`
}

type LogisticsDeliveryStatusChanged struct {
	EventID     string    `json:"event_id"`
	ShipmentID  string    `json:"shipment_id"`
	OrderID     string    `json:"order_id"`
	StoreID     string    `json:"store_id"`
	WarehouseID string    `json:"warehouse_id,omitempty"`
	DriverID    string    `json:"driver_id"`
	VehicleID   string    `json:"vehicle_id"`
	Status      string    `json:"status"`
	OccurredAt  time.Time `json:"occurred_at"`
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

type LogisticsShipmentStatusChanged struct {
	EventID    string    `json:"event_id"`
	ShipmentID string    `json:"shipment_id"`
	OrderID    string    `json:"order_id,omitempty"`
	HarvestID  string    `json:"harvest_id,omitempty"`
	StoreID    string    `json:"store_id,omitempty"`
	FarmID     string    `json:"farm_id,omitempty"`
	DriverID   string    `json:"driver_id"`
	VehicleID  string    `json:"vehicle_id"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

type NotificationCreated struct {
	EventID     string                 `json:"event_id"`
	UserID      string                 `json:"user_id,omitempty"`
	Role        string                 `json:"role,omitempty"`
	StoreID     string                 `json:"store_id,omitempty"`
	FarmID      string                 `json:"farm_id,omitempty"`
	WarehouseID string                 `json:"warehouse_id,omitempty"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Severity    string                 `json:"severity"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	OccurredAt  time.Time              `json:"occurred_at"`
}

type NotificationAcknowledged struct {
	EventID        string    `json:"event_id"`
	NotificationID string    `json:"notification_id"`
	UserID         string    `json:"user_id"`
	OccurredAt     time.Time `json:"occurred_at"`
}

type SocketBroadcastRequested struct {
	EventID       string                 `json:"event_id"`
	Channel       string                 `json:"channel"`
	FlowID        string                 `json:"flow_id,omitempty"`
	NodeID        string                 `json:"node_id,omitempty"`
	EdgeID        string                 `json:"edge_id,omitempty"`
	Visibility    string                 `json:"visibility,omitempty"`
	Status        string                 `json:"status,omitempty"`
	SourceService string                 `json:"source_service,omitempty"`
	Role          string                 `json:"role,omitempty"`
	StoreID       string                 `json:"store_id,omitempty"`
	FarmID        string                 `json:"farm_id,omitempty"`
	WarehouseID   string                 `json:"warehouse_id,omitempty"`
	EventType     string                 `json:"event_type"`
	Payload       map[string]interface{} `json:"payload"`
	OccurredAt    time.Time              `json:"occurred_at"`
}
