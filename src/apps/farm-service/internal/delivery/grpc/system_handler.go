package grpc

import (
	"context"

	"RuntimeRoasters/apps/farm-service/internal/usecase"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SystemHandler struct {
	systemv1.UnimplementedSystemServiceServer
	usecase usecase.SystemUsecase
}

func NewSystemHandler(u usecase.SystemUsecase) *SystemHandler {
	return &SystemHandler{
		usecase: u,
	}
}

func (h *SystemHandler) GetStatus(ctx context.Context, req *systemv1.GetStatusRequest) (*systemv1.GetStatusResponse, error) {
	seeded, counts, err := h.usecase.GetStatus(ctx)
	if err != nil {
		return nil, err
	}

	return &systemv1.GetStatusResponse{
		Seeded:         seeded,
		ServiceName:    "farm-service",
		RecordCounts:   counts,
		LastSeededAt:   timestamppb.Now(), // Simplified for now
	}, nil
}

func (h *SystemHandler) SeedData(ctx context.Context, req *systemv1.SeedDataRequest) (*systemv1.SeedDataResponse, error) {
	created, err := h.usecase.SeedData(ctx, req.UsersMap)
	if err != nil {
		return &systemv1.SeedDataResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &systemv1.SeedDataResponse{
		Success:        true,
		Message:        "Data seeded successfully",
		RecordsCreated: created,
	}, nil
}
