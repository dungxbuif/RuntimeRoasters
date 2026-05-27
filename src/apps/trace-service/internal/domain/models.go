package domain

import "time"

type TraceEvent struct {
	ID             string    `gorm:"type:uuid;primaryKey" json:"id"`
	MessageID      string    `gorm:"uniqueIndex;not null" json:"message_id"`
	Topic          string    `gorm:"size:160;index;not null" json:"topic"`
	BatchID        string    `gorm:"size:100;index" json:"batch_id"`
	OrderID        string    `gorm:"size:100;index" json:"order_id"`
	ShipmentID     string    `gorm:"size:100;index" json:"shipment_id"`
	StoreID        string    `gorm:"size:100;index" json:"store_id"`
	HarvestID      string    `gorm:"size:100;index" json:"harvest_id"`
	FarmID         string    `gorm:"size:100;index" json:"farm_id"`
	WarehouseID    string    `gorm:"size:100;index" json:"warehouse_id"`
	DriverID       string    `gorm:"size:100;index" json:"driver_id"`
	VehicleID      string    `gorm:"size:100;index" json:"vehicle_id"`
	TraceID        string    `gorm:"size:64;index" json:"trace_id,omitempty"`
	FlowID         string    `gorm:"size:100;index" json:"flow_id,omitempty"`
	NodeID         string    `gorm:"size:100;index" json:"node_id,omitempty"`
	EdgeID         string    `gorm:"size:100;index" json:"edge_id,omitempty"`
	Pattern        string    `gorm:"size:100;index" json:"pattern,omitempty"`
	SourceService  string    `gorm:"size:100;index" json:"source_service,omitempty"`
	Visibility     string    `gorm:"size:40;index;default:'public'" json:"visibility,omitempty"`
	DisplayPayload string    `gorm:"type:jsonb;not null;default:'{}'" json:"display_payload,omitempty"`
	Payload        string    `gorm:"type:jsonb;not null" json:"payload"`
	OccurredAt     time.Time `gorm:"index;not null" json:"occurred_at"`
	CreatedAt      time.Time `json:"created_at"`
}

type TraceDocument struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	EntityID   string    `gorm:"size:120;uniqueIndex;not null" json:"entity_id"`
	EntityType string    `gorm:"size:40;index;not null" json:"entity_type"`
	StoreID    string    `gorm:"size:100;index" json:"store_id"`
	TraceIDs   string    `gorm:"type:jsonb;not null;default:'[]'" json:"trace_ids"`
	Document   string    `gorm:"type:jsonb;not null" json:"document"`
	UpdatedAt  time.Time `json:"updated_at"`
}
