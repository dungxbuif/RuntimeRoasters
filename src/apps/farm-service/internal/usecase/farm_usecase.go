package usecase

import (
	"context"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
)

type FarmUsecase interface {
	GetFarm(ctx context.Context) (*domain.Farm, error)
}

type farmUsecase struct{}

func NewFarmUsecase() FarmUsecase {
	return &farmUsecase{}
}

func (u *farmUsecase) GetFarm(ctx context.Context) (*domain.Farm, error) {
	d := &domain.Farm{
		ID:      "farm-1",
		Message: "Pong! Farm Service is alive.",
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}
