package grpc

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
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

func (h *FarmHandler) CreateFarm(ctx context.Context, req *farmv1.CreateFarmRequest) (*farmv1.CreateFarmResponse, error) {
	log := logger.FromContext(ctx)

	f := &domain.Farm{
		Name:       req.Name,
		Location:   req.Location,
		Area:       req.Area,
		CoffeeType: req.FarmType,
		OwnerID:    req.OwnerId,
	}

	res, err := h.usecase.CreateFarm(ctx, f)
	if err != nil {
		log.Error("failed to create farm", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.CreateFarmResponse{
		Farm: mapToProto(res),
	}, nil
}

func (h *FarmHandler) GetFarm(ctx context.Context, req *farmv1.GetFarmRequest) (*farmv1.GetFarmResponse, error) {
	log := logger.FromContext(ctx)

	res, err := h.usecase.GetFarm(ctx, req.Id)
	if err != nil {
		log.Error("failed to get farm", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.GetFarmResponse{
		Farm: mapToProto(res),
	}, nil
}

func (h *FarmHandler) ListFarms(ctx context.Context, req *farmv1.ListFarmsRequest) (*farmv1.ListFarmsResponse, error) {
	log := logger.FromContext(ctx)

	res, err := h.usecase.ListFarms(ctx)
	if err != nil {
		log.Error("failed to list farms", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	farms := make([]*farmv1.Farm, len(res))
	for i, f := range res {
		farms[i] = mapToProto(f)
	}

	return &farmv1.ListFarmsResponse{
		Farms: farms,
	}, nil
}

func (h *FarmHandler) UpdateFarm(ctx context.Context, req *farmv1.UpdateFarmRequest) (*farmv1.UpdateFarmResponse, error) {
	log := logger.FromContext(ctx)

	f := &domain.Farm{
		ID:         req.Id,
		Name:       req.Name,
		Location:   req.Location,
		Area:       req.Area,
		CoffeeType: req.FarmType,
	}

	res, err := h.usecase.UpdateFarm(ctx, f)
	if err != nil {
		log.Error("failed to update farm", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.UpdateFarmResponse{
		Farm: mapToProto(res),
	}, nil
}

func (h *FarmHandler) DeleteFarm(ctx context.Context, req *farmv1.DeleteFarmRequest) (*farmv1.DeleteFarmResponse, error) {
	log := logger.FromContext(ctx)

	err := h.usecase.DeleteFarm(ctx, req.Id)
	if err != nil {
		log.Error("failed to delete farm", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.DeleteFarmResponse{
		Success: true,
	}, nil
}

func mapToProto(f *domain.Farm) *farmv1.Farm {
	return &farmv1.Farm{
		Id:        f.ID,
		Name:      f.Name,
		Location:  f.Location,
		Area:      f.Area,
		FarmType:  f.CoffeeType,
		OwnerId:   f.OwnerID,
		CreatedAt: f.CreatedAt.Unix(),
		UpdatedAt: f.UpdatedAt.Unix(),
	}
}
