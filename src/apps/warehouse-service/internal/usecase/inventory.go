package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"gorm.io/gorm"
)

type InventoryUseCase struct {
	db         *gorm.DB
	producer   kafka.Producer
	stockTopic string
}

func NewInventoryUseCase(db *gorm.DB, producer kafka.Producer, stockTopic string) *InventoryUseCase {
	return &InventoryUseCase{
		db:         db,
		producer:   producer,
		stockTopic: stockTopic,
	}
}

func (uc *InventoryUseCase) FinalizeBatch(ctx context.Context, internalID string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get Batch
		var batch domain.ProductionBatch
		if err := tx.Preload("Intakes").Where("id = ?", internalID).First(&batch).Error; err != nil {
			return err
		}

		if batch.Status != domain.BatchStatusReady {
			return fmt.Errorf("only READY batches can be finalized, current: %s", batch.Status)
		}

		// 2. Update Inventory
		sku := fmt.Sprintf("%s-%s-ROASTED", batch.Intakes[0].OriginCode, batch.Intakes[0].CoffeeType)
		var inv domain.Inventory
		if err := tx.Where("sku = ?", sku).First(&inv).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				inv = domain.Inventory{
					CoffeeType:        batch.Intakes[0].CoffeeType,
					OriginCode:        batch.Intakes[0].OriginCode,
					SKU:               sku,
					AvailableQuantity: batch.TotalOutputWeight,
				}
				tx.Create(&inv)
			} else {
				return err
			}
		} else {
			tx.Model(&inv).Update("available_quantity", gorm.Expr("available_quantity + ?", batch.TotalOutputWeight))
		}

		// 3. Finalize Batch Status
		if err := tx.Model(&batch).Update("status", domain.BatchStatusStocked).Error; err != nil {
			return err
		}

		// 4. Notify (Simulated Outbox - in real use, we write to outbox table)
		event := map[string]interface{}{
			"batch_id":  batch.BatchID,
			"sku":       sku,
			"quantity":  batch.TotalOutputWeight,
			"timestamp": time.Now(),
		}

		// For demo, we publish directly but note the intent
		return uc.producer.Publish(ctx, uc.stockTopic, batch.BatchID, event)
	})
}
