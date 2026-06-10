package domain

import (
	"context"
	"time"
)

type UserSecurityState struct {
	Subject          string    `json:"subject" gorm:"primaryKey"`
	TokensValidAfter time.Time `json:"tokens_valid_after"`
	Reason           string    `json:"reason"`
	UpdatedBy        string    `json:"updated_by"`
	UpdatedAt        time.Time `json:"updated_at"`
	Version          int64     `json:"version"`
}

type SecurityStateRepository interface {
	GetBySubject(ctx context.Context, subject string) (*UserSecurityState, error)
	MonotonicUpdate(ctx context.Context, subject string, reason string, updatedBy string) (*UserSecurityState, error)
}
