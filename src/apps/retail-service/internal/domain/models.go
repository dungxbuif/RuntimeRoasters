package domain

import (
	"time"
)

type StoreStatus string
type OrderStatus string
type OutboxStatus string

const (
	StoreStatusActive   StoreStatus = "ACTIVE"
	StoreStatusInactive StoreStatus = "INACTIVE"

	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusShipping  OrderStatus = "SHIPPING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusRejected  OrderStatus = "REJECTED"

	OutboxStatusPending   OutboxStatus = "PENDING"
	OutboxStatusCompleted OutboxStatus = "COMPLETED"
	OutboxStatusFailed    OutboxStatus = "FAILED"
)

type Store struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string      `gorm:"size:120;not null" json:"name"`
	City         string      `gorm:"size:80;not null;index" json:"city"`
	Address      string      `gorm:"size:255;not null" json:"address"`
	ManagerID    string      `gorm:"type:varchar(64);index" json:"manager_id"`
	ManagerEmail string      `gorm:"size:160;index" json:"manager_email"`
	Status       StoreStatus `gorm:"size:20;not null;default:'ACTIVE'" json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
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
