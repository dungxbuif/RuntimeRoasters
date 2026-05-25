package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/errs"
	"RuntimeRoasters/pkg/events"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

type harvestCreatedEvent struct {
	HarvestID  string  `json:"harvest_id"`
	CoffeeType string  `json:"coffee_type"`
	OriginCode string  `json:"origin_code"`
	Quantity   float64 `json:"quantity"`
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
	farm, err := u.farmRepo.GetByID(ctx, harvest.FarmID)
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

		occurredAt := time.Now()
		payloadData := harvestCreatedEvent{
			HarvestID:  strconv.FormatUint(harvest.ID, 10),
			CoffeeType: string(harvest.CoffeeType),
			OriginCode: string(farm.Location),
			Quantity:   harvest.Quantity,
		}
		cloudEvent, err := events.NewCloudEvent(ctx, events.TopicFarmHarvestCreated, events.SourceFarmService, fmt.Sprintf("harvests/%d", harvest.ID), payloadData, events.Metadata{
			EventID:       eventID,
			CorrelationID: strconv.FormatUint(harvest.ID, 10),
			OccurredAt:    occurredAt,
			HarvestID:     strconv.FormatUint(harvest.ID, 10),
			FarmID:        strconv.FormatUint(harvest.FarmID, 10),
		})
		if err != nil {
			return err
		}
		payload, _ := json.Marshal(cloudEvent)

		traceHeaders := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(ctx, traceHeaders)
		metadata := map[string]interface{}{
			"trace_id": traceID,
			"user_id":  harvest.OwnerID,
		}
		if traceparent := traceHeaders.Get("traceparent"); traceparent != "" {
			metadata["traceparent"] = traceparent
		}
		if tracestate := traceHeaders.Get("tracestate"); tracestate != "" {
			metadata["tracestate"] = tracestate
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
