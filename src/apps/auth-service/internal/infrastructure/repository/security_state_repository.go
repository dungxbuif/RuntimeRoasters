package repository

import (
	"context"
	"errors"
	"time"

	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/pkg/database"
	"gorm.io/gorm"
)

type securityStateRepository struct {
	db *database.DB
}

func NewSecurityStateRepository(db *database.DB) domain.SecurityStateRepository {
	return &securityStateRepository{db: db}
}

func (r *securityStateRepository) GetBySubject(ctx context.Context, subject string) (*domain.UserSecurityState, error) {
	var state domain.UserSecurityState
	tx := r.db.GetTx(ctx)
	err := tx.Where("subject = ?", subject).First(&state).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &domain.UserSecurityState{
				Subject:          subject,
				TokensValidAfter: time.Unix(0, 0),
				Version:          0,
			}, nil
		}
		return nil, err
	}
	return &state, nil
}

func (r *securityStateRepository) MonotonicUpdate(ctx context.Context, subject string, reason string, updatedBy string) (*domain.UserSecurityState, error) {
	var state domain.UserSecurityState
	tx := r.db.GetTx(ctx)

	query := `
		INSERT INTO user_security_states (subject, tokens_valid_after, reason, updated_by, updated_at, version)
		VALUES (?, CURRENT_TIMESTAMP, ?, ?, CURRENT_TIMESTAMP, 1)
		ON CONFLICT (subject) DO UPDATE
		SET tokens_valid_after = GREATEST(user_security_states.tokens_valid_after, CURRENT_TIMESTAMP),
			reason = EXCLUDED.reason,
			updated_by = EXCLUDED.updated_by,
			updated_at = CURRENT_TIMESTAMP,
			version = user_security_states.version + 1
		RETURNING subject, tokens_valid_after, reason, updated_by, updated_at, version;
	`
	err := tx.Raw(query, subject, reason, updatedBy).Scan(&state).Error
	if err != nil {
		return nil, err
	}
	return &state, nil
}
