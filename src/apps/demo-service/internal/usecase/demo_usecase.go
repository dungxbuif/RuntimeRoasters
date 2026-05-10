package usecase

import (
	"context"
	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/domain"
)

type DemoUsecase interface {
	GetDemo(ctx context.Context) (*domain.Demo, error)
}

type demoUsecase struct{}

func NewDemoUsecase() DemoUsecase {
	return &demoUsecase{}
}

func (u *demoUsecase) GetDemo(ctx context.Context) (*domain.Demo, error) {
	d := &domain.Demo{
		ID:      "demo-1",
		Message: "Pong! Demo Service is alive.",
	}

	if err := d.Validate(); err != nil {
		return nil, err
	}

	return d, nil
}
