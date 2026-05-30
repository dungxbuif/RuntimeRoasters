package grpc

import (
	"context"

	"RuntimeRoasters/apps/retail-service/internal/usecase"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SystemHandler struct {
	systemv1.UnimplementedSystemServiceServer
	service *usecase.Service
}

func NewSystemHandler(service *usecase.Service) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) GetStatus(ctx context.Context, req *systemv1.GetStatusRequest) (*systemv1.GetStatusResponse, error) {
	status, err := h.service.GetSystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	res := &systemv1.GetStatusResponse{
		Seeded:      status.Seeded,
		ServiceName: status.ServiceName,
		RecordCounts: status.RecordCounts,
	}

	if status.LastSeededAt != nil {
		res.LastSeededAt = timestamppb.New(*status.LastSeededAt)
	}

	return res, nil
}

func (h *SystemHandler) SeedData(ctx context.Context, req *systemv1.SeedDataRequest) (*systemv1.SeedDataResponse, error) {
	result, err := h.service.SeedData(ctx, req.Force, req.UsersMap)
	if err != nil {
		return nil, err
	}

	return &systemv1.SeedDataResponse{
		Success:        result.Success,
		Message:        result.Message,
		RecordsCreated: result.RecordsCreated,
	}, nil
}
