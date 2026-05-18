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
	if stockTopic == "" {
		stockTopic = "warehouse.stock.updated"
	}
	return &InventoryUseCase{
		db:         db,
		producer:   producer,
		stockTopic: stockTopic,
	}
}

func (uc *InventoryUseCase) FinalizeBatch(ctx context.Context, batchID string) error {
	return uc.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get Batch
		var batch domain.ProductionBatch
		if err := tx.Where("batch_id = ?", batchID).First(&batch).Error; err != nil {
			return err
		}

		if batch.Status != domain.BatchStatusProcessing {
			return fmt.Errorf("only PROCESSING batches can be finalized")
		}

		// 2. Generate Final SKU and ID
		sku := fmt.Sprintf("%s-%s-ROASTED", batch.OriginCode, batch.CoffeeType)

		// 3. Update Inventory
		var inv domain.Inventory
		err := tx.Where("sku = ?", sku).First(&inv).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new inventory record
				inv = domain.Inventory{
					CoffeeType:        batch.CoffeeType,
					OriginCode:        batch.OriginCode,
					SKU:               sku,
					AvailableQuantity: batch.OutputWeight,
				}
				if err := tx.Create(&inv).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// Update existing
			if err := tx.Model(&inv).Update("available_quantity", gorm.Expr("available_quantity + ?", batch.OutputWeight)).Error; err != nil {
				return err
			}
		}

		// 4. Finalize Batch Status
		if err := tx.Model(&batch).Update("status", domain.BatchStatusStocked).Error; err != nil {
			return err
		}

		// 5. Notify via Kafka
		event := map[string]interface{}{
			"batch_id":  batchID,
			"sku":       sku,
			"quantity":  batch.OutputWeight,
			"timestamp": time.Now(),
		}

		return uc.producer.Publish(ctx, uc.stockTopic, batchID, event)
	})
}
