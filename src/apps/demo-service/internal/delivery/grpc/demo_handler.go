package grpc

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	demov1 "github.com/dungxbuif/RuntimeRoasters/runtime/demo/v1"
	"go.uber.org/zap"
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
	log := logger.FromContext(ctx)

	// Debug log: Check if identity is propagated
	if id, ok := identity.FromContext(ctx); ok {
		log.Info("Request authenticated",
			zap.String("user_id", id.Subject),
			zap.String("role", id.Role),
		)
	} else {
		log.Warn("Request NOT authenticated (missing context identity)")
	}

	demo, err := h.usecase.GetDemo(ctx)
	if err != nil {
		return nil, errs.ToGRPCError(err)
	}

	return &demov1.GetDemoResponse{
		Id:      demo.ID,
		Message: demo.Message,
	}, nil
}
