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
