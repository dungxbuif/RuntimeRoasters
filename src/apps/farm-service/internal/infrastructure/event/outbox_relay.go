package event

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

type OutboxRelay struct {
	outboxRepo usecase.OutboxRepository
	publisher  usecase.EventPublisher
	interval   time.Duration
}

func NewOutboxRelay(
	outboxRepo usecase.OutboxRepository,
	publisher usecase.EventPublisher,
	interval time.Duration,
) *OutboxRelay {
	return &OutboxRelay{
		outboxRepo: outboxRepo,
		publisher:  publisher,
		interval:   interval,
	}
}

func (w *OutboxRelay) Start(ctx context.Context) {
	log := logger.GetLogger().With(zap.String("component", "outbox-relay"))
	log.Info("Outbox relay started", zap.Duration("interval", w.interval))

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Outbox relay stopping...")
			return
		case <-ticker.C:
			w.processEvents(ctx, log)
		}
	}
}

func (w *OutboxRelay) processEvents(ctx context.Context, log *zap.Logger) {
	events, err := w.outboxRepo.ListPending(ctx, 10)
	if err != nil {
		log.Error("failed to list pending events", zap.Error(err))
		return
	}

	for _, event := range events {
		log.Info("Processing outbox event", zap.String("id", event.ID), zap.String("type", event.EventType))

		// Mark as processing
		event.Status = domain.OutboxStatusProcessing
		_ = w.outboxRepo.Update(ctx, event)

		// Mock Publish
		err = w.publisher.Publish(ctx, event)
		if err != nil {
			log.Error("failed to publish event", zap.Error(err), zap.String("id", event.ID))
			event.RetryCount++
			if event.RetryCount >= domain.MaxRetryCount {
				event.Status = domain.OutboxStatusFailed
			} else {
				event.Status = domain.OutboxStatusPending
			}
			_ = w.outboxRepo.Update(ctx, event)
			continue
		}

		// Mark as completed
		now := time.Now()
		event.Status = domain.OutboxStatusCompleted
		event.ProcessedAt = &now
		if err := w.outboxRepo.Update(ctx, event); err != nil {
			log.Error("failed to update event status", zap.Error(err))
		} else {
			log.Info("Successfully processed outbox event", zap.String("id", event.ID))
		}
	}
}
