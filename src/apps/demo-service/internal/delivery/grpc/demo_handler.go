package grpc

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/usecase"
	demov1 "github.com/dungxbuif/RuntimeRoasters/runtime/demo/v1"
)

type DemoHandler struct {
	demov1.UnimplementedDemoServiceServer
	usecase usecase.DemoUsecase
}

func NewDemoHandler(u usecase.DemoUsecase) *DemoHandler {
	return &DemoHandler{
		usecase: u,
	}
}

func (h *DemoHandler) GetDemo(ctx context.Context, req *demov1.GetDemoRequest) (*demov1.GetDemoResponse, error) {
	demo, err := h.usecase.GetDemo(ctx)
	if err != nil {
		return nil, err
	}

	return &demov1.GetDemoResponse{
		Id:      demo.ID,
		Message: demo.Message,
	}, nil
}
