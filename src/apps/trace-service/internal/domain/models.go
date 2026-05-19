package domain

import "time"

type TraceEvent struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	MessageID  string    `gorm:"uniqueIndex;not null" json:"message_id"`
	Topic      string    `gorm:"size:160;index;not null" json:"topic"`
	BatchID    string    `gorm:"size:100;index" json:"batch_id"`
	OrderID    string    `gorm:"size:100;index" json:"order_id"`
	ShipmentID string    `gorm:"size:100;index" json:"shipment_id"`
	StoreID    string    `gorm:"size:100;index" json:"store_id"`
	Payload    string    `gorm:"type:jsonb;not null" json:"payload"`
	OccurredAt time.Time `gorm:"index;not null" json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type TraceDocument struct {
	ID         string    `gorm:"type:uuid;primaryKey" json:"id"`
	EntityID   string    `gorm:"size:120;uniqueIndex;not null" json:"entity_id"`
	EntityType string    `gorm:"size:40;index;not null" json:"entity_type"`
	StoreID    string    `gorm:"size:100;index" json:"store_id"`
	Document   string    `gorm:"type:jsonb;not null" json:"document"`
	UpdatedAt  time.Time `json:"updated_at"`
}
