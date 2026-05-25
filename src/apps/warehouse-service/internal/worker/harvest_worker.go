package worker

import (
	"context"

	"RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"RuntimeRoasters/pkg/kafka"
	kafka_go "github.com/segmentio/kafka-go"
)

type HarvestWorker struct {
	consumer kafka.Consumer
	usecase  *usecase.PickupUseCase
}

func NewHarvestWorker(consumer kafka.Consumer, usecase *usecase.PickupUseCase) *HarvestWorker {
	return &HarvestWorker{
		consumer: consumer,
		usecase:  usecase,
	}
}

func (w *HarvestWorker) Start(ctx context.Context) error {
	handler := func(ctx context.Context, msg kafka_go.Message) error {
		msgID := kafka.MessageID(msg)

		return w.usecase.ProcessHarvestEvent(ctx, msgID, msg.Value)
	}

	return w.consumer.Listen(ctx, handler)
}
