package worker

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
)

type OrderWorker struct {
	consumer kafka.Consumer
	usecase  *usecase.OrderReservationUseCase
}

func NewOrderWorker(consumer kafka.Consumer, usecase *usecase.OrderReservationUseCase) *OrderWorker {
	return &OrderWorker{consumer: consumer, usecase: usecase}
}

func (w *OrderWorker) Start(ctx context.Context) error {
	return w.consumer.Listen(ctx, w.usecase.ProcessOrderCreated)
}
