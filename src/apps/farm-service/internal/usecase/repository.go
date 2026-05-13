package usecase

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
)

type FarmRepository interface {
	Create(ctx context.Context, farm *domain.Farm) error
	GetByID(ctx context.Context, id uint64) (*domain.Farm, error)
	List(ctx context.Context) ([]*domain.Farm, error)
	Update(ctx context.Context, farm *domain.Farm) error
	Delete(ctx context.Context, id uint64) error
}

type HarvestRepository interface {
	Create(ctx context.Context, harvest *domain.Harvest) error
	GetByID(ctx context.Context, id uint64) (*domain.Harvest, error)
	ListByFarm(ctx context.Context, farmID uint64) ([]*domain.Harvest, error)
	Update(ctx context.Context, harvest *domain.Harvest) error
	Delete(ctx context.Context, id uint64) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event *domain.OutboxEvent) error
}

type OutboxRepository interface {
	Create(ctx context.Context, event *domain.OutboxEvent) error
	ListPending(ctx context.Context, limit int) ([]*domain.OutboxEvent, error)
	Update(ctx context.Context, event *domain.OutboxEvent) error
}



