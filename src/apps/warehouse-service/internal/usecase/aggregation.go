package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
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
		for _, it := range intakes {
			totalWeight += it.Quantity
		}

		// 3. Create Production Batch
		batchID = fmt.Sprintf("BATCH-%s-%s-%d", origin, time.Now().Format("0102"), rand.Intn(999))
		batch := domain.ProductionBatch{
			ID:               fmt.Sprintf("%d", time.Now().UnixNano()), // Simplified UUID for demo
			BatchID:          batchID,
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
	err := uc.db.Where("status = ?", domain.IntakeStatusUnassigned).Find(&intakes).Error
	return intakes, err
}
