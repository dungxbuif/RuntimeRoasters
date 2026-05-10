package grpc

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	farmv1 "github.com/dungxbuif/RuntimeRoasters/runtime/farm/v1"
	"go.uber.org/zap"
)

type FarmHandler struct {
	farmv1.UnimplementedFarmServiceServer
	usecase usecase.FarmUsecase
}

func NewFarmHandler(u usecase.FarmUsecase) *FarmHandler {
	return &FarmHandler{
		usecase: u,
	}
}

func (h *FarmHandler) GetFarm(ctx context.Context, req *farmv1.GetFarmRequest) (*farmv1.GetFarmResponse, error) {
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

	farm, err := h.usecase.GetFarm(ctx)
	if err != nil {
		log.Error("failed to get farm", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.GetFarmResponse{
		Id:      farm.ID,
		Message: farm.Message,
	}, nil
}
