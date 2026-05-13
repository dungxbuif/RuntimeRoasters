package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "PENDING"
	OutboxStatusProcessing OutboxStatus = "PROCESSING"
	OutboxStatusCompleted  OutboxStatus = "COMPLETED"
	OutboxStatusFailed     OutboxStatus = "FAILED"
)

const MaxRetryCount = 5

type OutboxEvent struct {
	ID          string          `gorm:"type:uuid;primaryKey;"`
	EventType   string          `gorm:"size:100;not null"`
	Payload     json.RawMessage `gorm:"type:jsonb;not null"`
	Metadata    json.RawMessage `gorm:"type:jsonb"`
	RetryCount  int             `gorm:"default:0"`
	Status      OutboxStatus    `gorm:"size:20;default:'PENDING';index:idx_outbox_unprocessed,where:status = 'PENDING'"`
	ProcessedAt *time.Time      `gorm:"type:timestamptz;default:null"`
	CreatedAt   time.Time       `gorm:"type:timestamptz;default:now()"`
}

func (o *OutboxEvent) Validate() error {
	if o.ID == "" {
		return errors.New("id is required")
	}
	if o.EventType == "" {
		return errors.New("event type is required")
	}
	if o.Payload == nil {
		return errors.New("payload is required")
	}
	if o.Metadata == nil {
		return errors.New("metadata is required")
	}
	if o.RetryCount < 0 {
		return errors.New("retry count must be greater than or equal to 0")
	}
	if o.Status == "" {
		return errors.New("status is required")
	}
	if o.ProcessedAt == nil {
		return errors.New("processed at is required")
	}
	if o.CreatedAt.IsZero() {
		return errors.New("created at is required")
	}
	return nil
}