package worker

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
)

type OrderWorker struct {
	consumers []kafka.Consumer
	usecase   *usecase.OrderReservationUseCase
}

func NewOrderWorker(consumers []kafka.Consumer, usecase *usecase.OrderReservationUseCase) *OrderWorker {
	return &OrderWorker{consumers: consumers, usecase: usecase}
}

func (w *OrderWorker) Start(ctx context.Context) error {
	errCh := make(chan error, len(w.consumers))
	for _, consumer := range w.consumers {
		c := consumer
		go func() {
			errCh <- c.Listen(ctx, w.usecase.ProcessOrderCreated)
		}()
	}
	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}
