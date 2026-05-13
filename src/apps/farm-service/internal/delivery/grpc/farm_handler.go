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
	usecase        usecase.FarmUsecase
	harvestUsecase usecase.HarvestUsecase
}

func NewFarmHandler(u usecase.FarmUsecase, h usecase.HarvestUsecase) *FarmHandler {
	return &FarmHandler{
		usecase:        u,
		harvestUsecase: h,
	}
}

func (h *FarmHandler) CreateFarm(ctx context.Context, req *farmv1.CreateFarmRequest) (*farmv1.CreateFarmResponse, error) {
	log := logger.FromContext(ctx)

	f := &domain.Farm{
		Name:       req.Name,
		Location:   domain.Location(req.Location),
		Area:       req.Area,
		CoffeeType: domain.CoffeeType(req.FarmType),
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
		Location:   domain.Location(req.Location),
		Area:       req.Area,
		CoffeeType: domain.CoffeeType(req.FarmType),
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


func (h *FarmHandler) CreateHarvest(ctx context.Context, req *farmv1.CreateHarvestRequest) (*farmv1.CreateHarvestResponse, error) {
	log := logger.FromContext(ctx)

	harvest := &domain.Harvest{
		FarmID:     req.FarmId,
		CoffeeType: domain.CoffeeType(req.CoffeeType),
		Quantity:   req.Quantity,
		Notes:      req.Notes,
		Status:     domain.StatusNew,
	}

	res, err := h.harvestUsecase.CreateHarvest(ctx, harvest)
	if err != nil {
		log.Error("failed to create harvest", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.CreateHarvestResponse{
		Harvest: mapHarvestToProto(res),
	}, nil
}

func (h *FarmHandler) GetHarvest(ctx context.Context, req *farmv1.GetHarvestRequest) (*farmv1.GetHarvestResponse, error) {
	log := logger.FromContext(ctx)

	res, err := h.harvestUsecase.GetHarvest(ctx, req.Id)
	if err != nil {
		log.Error("failed to get harvest", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	return &farmv1.GetHarvestResponse{
		Harvest: mapHarvestToProto(res),
	}, nil
}

func (h *FarmHandler) ListFarmHarvests(ctx context.Context, req *farmv1.ListFarmHarvestsRequest) (*farmv1.ListFarmHarvestsResponse, error) {
	log := logger.FromContext(ctx)

	res, err := h.harvestUsecase.ListFarmHarvests(ctx, req.FarmId)
	if err != nil {
		log.Error("failed to list harvests", zap.Error(err))
		return nil, errs.ToGRPCError(err)
	}

	harvests := make([]*farmv1.Harvest, len(res))
	for i, hv := range res {
		harvests[i] = mapHarvestToProto(hv)
	}

	return &farmv1.ListFarmHarvestsResponse{
		Harvests: harvests,
	}, nil
}


func mapToProto(f *domain.Farm) *farmv1.Farm {
	return &farmv1.Farm{
		Id:        f.ID,
		Name:      f.Name,
		Location:  string(f.Location),
		Area:      f.Area,
		FarmType:  string(f.CoffeeType),
		OwnerId:   f.OwnerID,
		CreatedAt: f.CreatedAt.Unix(),
		UpdatedAt: f.UpdatedAt.Unix(),
	}
}

func mapHarvestToProto(h *domain.Harvest) *farmv1.Harvest {
	return &farmv1.Harvest{
		Id:          h.ID,
		FarmId:      h.FarmID,
		OwnerId:     h.OwnerID,
		CoffeeType:  string(h.CoffeeType),
		Quantity:    h.Quantity,
		HarvestDate: h.HarvestDate.Unix(),
		Status:      string(h.Status),
		Notes:       h.Notes,
		CreatedAt:   h.CreatedAt.Unix(),
		UpdatedAt:   h.UpdatedAt.Unix(),
	}
}
