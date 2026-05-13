package repository

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
)

type OutboxEventModel struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventType   string    `gorm:"not null"`
	Payload     []byte    `gorm:"type:jsonb;not null"`
	Metadata    []byte    `gorm:"type:jsonb"`
	RetryCount  int       `gorm:"default:0"`
	Status      domain.OutboxStatus `gorm:"type:outbox_status;default:'PENDING'"`
	ProcessedAt *time.Time
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (OutboxEventModel) TableName() string {
	return "outbox_events"
}

type outboxRepository struct {
	db *database.DB
}

func NewOutboxRepository(db *database.DB) usecase.OutboxRepository {
	return &outboxRepository{db: db}
}

func (r *outboxRepository) Create(ctx context.Context, event *domain.OutboxEvent) error {
	model := &OutboxEventModel{
		ID:         event.ID,
		EventType:  event.EventType,
		Payload:    event.Payload,
		Metadata:   event.Metadata,
		Status:     event.Status,
		CreatedAt:  event.CreatedAt,
	}
	return r.db.GetTx(ctx).Create(model).Error
}

func (r *outboxRepository) ListPending(ctx context.Context, limit int) ([]*domain.OutboxEvent, error) {
	var models []OutboxEventModel
	err := r.db.WithContext(ctx).
		Where("status = ? AND retry_count < ?", domain.OutboxStatusPending, domain.MaxRetryCount).
		Order("created_at ASC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	events := make([]*domain.OutboxEvent, len(models))
	for i, m := range models {
		events[i] = &domain.OutboxEvent{
			ID:          m.ID,
			EventType:   m.EventType,
			Payload:     m.Payload,
			Metadata:    m.Metadata,
			RetryCount:  m.RetryCount,
			Status:      m.Status,
			ProcessedAt: m.ProcessedAt,
			CreatedAt:   m.CreatedAt,
		}
	}
	return events, nil
}

func (r *outboxRepository) Update(ctx context.Context, event *domain.OutboxEvent) error {
	return r.db.WithContext(ctx).Model(&OutboxEventModel{}).
		Where("id = ?", event.ID).
		Updates(map[string]interface{}{
			"status":       event.Status,
			"processed_at": event.ProcessedAt,
			"retry_count":  event.RetryCount,
		}).Error
}
