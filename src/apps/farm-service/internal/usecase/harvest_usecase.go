package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type HarvestUsecase interface {
	CreateHarvest(ctx context.Context, harvest *domain.Harvest) (*domain.Harvest, error)
	GetHarvest(ctx context.Context, id uint64) (*domain.Harvest, error)
	ListFarmHarvests(ctx context.Context, farmID uint64) ([]*domain.Harvest, error)
	DeleteHarvest(ctx context.Context, id uint64) error
}

type harvestUsecase struct {
	db         *database.DB
	repo       HarvestRepository
	farmRepo   FarmRepository
	outboxRepo OutboxRepository
}

func NewHarvestUsecase(
	db *database.DB,
	repo HarvestRepository,
	farmRepo FarmRepository,
	outboxRepo OutboxRepository,
) HarvestUsecase {
	return &harvestUsecase{
		db:         db,
		repo:       repo,
		farmRepo:   farmRepo,
		outboxRepo: outboxRepo,
	}
}

func (u *harvestUsecase) CreateHarvest(ctx context.Context, harvest *domain.Harvest) (*domain.Harvest, error) {
	_, err := u.farmRepo.GetByID(ctx, harvest.FarmID)
	if err != nil {
		return nil, err
	}

	user, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}
	if (user.Role == domain.RoleAdmin || user.Role == domain.RoleFarmAdmin) && harvest.OwnerID != "" {
		// Admin can specify owner
	} else {
		harvest.OwnerID = user.Subject
	}

	if err := harvest.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrValidation, err)
	}

	harvest.Status = domain.StatusNew

	err = u.db.WithTx(ctx, func(txCtx context.Context) error {
		if err := u.repo.Create(txCtx, harvest); err != nil {
			return err
		}

		traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
		eventID := uuid.NewString()

		payload, _ := json.Marshal(harvest)

		metadata := map[string]interface{}{
			"trace_id": traceID,
			"user_id":  harvest.OwnerID,
		}
		metadataBytes, _ := json.Marshal(metadata)

		event := &domain.OutboxEvent{
			ID:        eventID,
			EventType: domain.EventTypeHarvestBatchCreated,
			Payload:   payload,
			Metadata:  metadataBytes,
			Status:    domain.OutboxStatusPending,
			CreatedAt: time.Now(),
		}

		if err := u.outboxRepo.Create(txCtx, event); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
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
