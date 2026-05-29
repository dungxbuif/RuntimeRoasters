package warehousegrpc

import (
	"context"

	"RuntimeRoasters/pkg/database"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SystemHandler struct {
	systemv1.UnimplementedSystemServiceServer
	db *database.DB
}

func NewSystemHandler(db *database.DB) *SystemHandler {
	return &SystemHandler{
		db: db,
	}
}

func (h *SystemHandler) GetStatus(ctx context.Context, req *systemv1.GetStatusRequest) (*systemv1.GetStatusResponse, error) {
	// For now, warehouse-service doesn't have a main Warehouse table, 
	// but we can count Intakes or Batches as a proxy for 'seeded' status.
	var intakeCount int64
	if err := h.db.Table("intakes").Count(&intakeCount).Error; err != nil {
		// Table might not exist yet if not migrated
		intakeCount = 0
	}

	return &systemv1.GetStatusResponse{
		Seeded:      intakeCount > 0,
		ServiceName: "warehouse-service",
		RecordCounts: map[string]int64{
			"intakes": intakeCount,
		},
	}, nil
}

func (h *SystemHandler) SeedData(ctx context.Context, req *systemv1.SeedDataRequest) (*systemv1.SeedDataResponse, error) {
	// Centralized seeding for warehouse-service.
	// Currently, warehouse-service relies on events to create data, 
	// but we could seed initial inventory or reference data here.
	
	// For this sprint, we'll just return success to allow the choreography to complete.
	return &systemv1.SeedDataResponse{
		Success: true,
		Message: "Warehouse service seeding acknowledged (no local entities to seed yet)",
	}, nil
}
