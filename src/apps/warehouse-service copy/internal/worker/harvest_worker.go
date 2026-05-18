package worker

import (
	"context"
	"fmt"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	kafka_go "github.com/segmentio/kafka-go"
)

type HarvestWorker struct {
	consumer kafka.Consumer
	usecase  *usecase.IntakeUseCase
}

func NewHarvestWorker(consumer kafka.Consumer, usecase *usecase.IntakeUseCase) *HarvestWorker {
	return &HarvestWorker{
		consumer: consumer,
		usecase:  usecase,
	}
}

func (w *HarvestWorker) Start(ctx context.Context) error {
	handler := func(ctx context.Context, msg kafka_go.Message) error {
		msgID := fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset)
		
		return w.usecase.ProcessHarvestEvent(ctx, msgID, msg.Value)
	}

	return w.consumer.Listen(ctx, handler)
}
