package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"gorm.io/gorm"
)

type AggregationUseCase struct {
	db *gorm.DB
}

func NewAggregationUseCase(db *gorm.DB) *AggregationUseCase {
	return &AggregationUseCase{db: db}
}

func (uc *AggregationUseCase) CreateProductionBatch(ctx context.Context, intakeIDs []string) (string, error) {
	var batchID string
	err := uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Validate and fetch Intakes
		var intakes []domain.Intake
		if err := tx.Where("id IN ? AND status = ?", intakeIDs, domain.IntakeStatusUnassigned).Find(&intakes).Error; err != nil {
			return err
		}

		if len(intakes) != len(intakeIDs) {
			return fmt.Errorf("some intakes are not available or already assigned")
		}

		// 2. Calculate totals
		var totalWeight float64
		origin := intakes[0].OriginCode
		warehouseID := intakes[0].WarehouseID
		for _, it := range intakes {
			if it.WarehouseID != warehouseID {
				return fmt.Errorf("all intakes must belong to the same warehouse")
			}
			totalWeight += it.Quantity
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(ctx); !allWarehouses && !claims.CanAccessWarehouse(warehouseID) {
			return fmt.Errorf("warehouse access denied")
		}

		// 3. Create Production Batch
		batchID = fmt.Sprintf("BATCH-%s-%s-%d", origin, time.Now().Format("0102"), rand.Intn(999))
		batch := domain.ProductionBatch{
			ID:               fmt.Sprintf("%d", time.Now().UnixNano()), // Simplified UUID for demo
			BatchID:          batchID,
			WarehouseID:      warehouseID,
			Status:           domain.BatchStatusDraft,
			TotalInputWeight: totalWeight,
		}

		if err := tx.Create(&batch).Error; err != nil {
			return err
		}

		// 4. Update Intakes to ASSIGNED
		if err := tx.Model(&domain.Intake{}).Where("id IN ?", intakeIDs).Updates(map[string]interface{}{
			"status":   domain.IntakeStatusAssigned,
			"batch_id": batch.ID,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	return batchID, err
}

func (uc *AggregationUseCase) GetUnassignedIntakes(ctx context.Context) ([]domain.Intake, error) {
	var intakes []domain.Intake
	query := uc.db.WithContext(ctx).Where("status = ?", domain.IntakeStatusUnassigned)
	_, warehouseIDs, allWarehouses := identity.WarehouseScopeFromContext(ctx)
	if !allWarehouses {
		if len(warehouseIDs) == 0 {
			return []domain.Intake{}, nil
		}
		query = query.Where("warehouse_id IN ?", warehouseIDs)
	}
	err := query.Find(&intakes).Error
	return intakes, err
}
