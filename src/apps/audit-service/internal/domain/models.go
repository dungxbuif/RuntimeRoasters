package domain

import "time"

type AuditLog struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	PartitionKey string    `gorm:"size:120;index;not null" json:"partition_key"`
	MessageID    string    `gorm:"uniqueIndex;not null" json:"message_id"`
	Topic        string    `gorm:"size:160;index;not null" json:"topic"`
	StoreID      string    `gorm:"size:100;index" json:"store_id"`
	Payload      string    `gorm:"type:jsonb;not null" json:"payload"`
	PreviousHash string    `gorm:"size:64" json:"previous_hash"`
	CurrentHash  string    `gorm:"size:64;not null" json:"current_hash"`
	OccurredAt   time.Time `gorm:"index;not null" json:"occurred_at"`
	CreatedAt    time.Time `json:"created_at"`
}
