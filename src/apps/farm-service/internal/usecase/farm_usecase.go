package usecase

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/infrastructure/repository"
	"github.com/google/uuid"
)

type FarmUsecase interface {
	CreateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error)
	GetFarm(ctx context.Context, id string) (*domain.Farm, error)
	ListFarms(ctx context.Context) ([]*domain.Farm, error)
	UpdateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error)
	DeleteFarm(ctx context.Context, id string) error
}

type farmUsecase struct {
	repo repository.FarmRepository
}

func NewFarmUsecase(repo repository.FarmRepository) FarmUsecase {
	return &farmUsecase{
		repo: repo,
	}
}

func (u *farmUsecase) CreateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error) {
	farm.ID = uuid.New().String()
	if err := farm.Validate(); err != nil {
		return nil, err
	}

	if err := u.repo.Create(ctx, farm); err != nil {
		return nil, err
	}

	return farm, nil
}

func (u *farmUsecase) GetFarm(ctx context.Context, id string) (*domain.Farm, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *farmUsecase) ListFarms(ctx context.Context) ([]*domain.Farm, error) {
	return u.repo.List(ctx)
}

func (u *farmUsecase) UpdateFarm(ctx context.Context, farm *domain.Farm) (*domain.Farm, error) {
	if err := farm.Validate(); err != nil {
		return nil, err
	}

	if err := u.repo.Update(ctx, farm); err != nil {
		return nil, err
	}

	return farm, nil
}

func (u *farmUsecase) DeleteFarm(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}
