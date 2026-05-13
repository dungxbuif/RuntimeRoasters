package usecase

import (
	"context"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
)

type HarvestUsecase interface {
	CreateHarvest(ctx context.Context, harvest *domain.Harvest) (*domain.Harvest, error)
	GetHarvest(ctx context.Context, id uint64) (*domain.Harvest, error)
	ListFarmHarvests(ctx context.Context, farmID uint64) ([]*domain.Harvest, error)
	DeleteHarvest(ctx context.Context, id uint64) error
}

type harvestUsecase struct {
	repo     HarvestRepository
	farmRepo FarmRepository
}

func NewHarvestUsecase(repo HarvestRepository, farmRepo FarmRepository) HarvestUsecase {
	return &harvestUsecase{
		repo:     repo,
		farmRepo: farmRepo,
	}
}

func (u *harvestUsecase) CreateHarvest(ctx context.Context, harvest *domain.Harvest) (*domain.Harvest, error) {
	_, err := u.farmRepo.GetByID(ctx, harvest.FarmID)
	if err != nil {
		return nil, err
	}
	if err := harvest.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrValidation, err)
	}
	if err := u.repo.Create(ctx, harvest); err != nil {
		return nil, err
	}

	return harvest, nil
}

func (u *harvestUsecase) GetHarvest(ctx context.Context, id uint64) (*domain.Harvest, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *harvestUsecase) ListFarmHarvests(ctx context.Context, farmID uint64) ([]*domain.Harvest, error) {
	return u.repo.ListByFarm(ctx, farmID)
}

func (u *harvestUsecase) DeleteHarvest(ctx context.Context, id uint64) error {
	return u.repo.Delete(ctx, id)
}
