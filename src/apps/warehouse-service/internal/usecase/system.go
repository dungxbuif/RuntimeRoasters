package usecase

import (
	"context"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"

	"gorm.io/gorm"
)

type SystemUseCase struct {
	db *gorm.DB
}

func NewSystemUseCase(db *gorm.DB) *SystemUseCase {
	return &SystemUseCase{db: db}
}

func (u *SystemUseCase) GetStatus(ctx context.Context) (bool, map[string]int64, error) {
	var pickupCount, inventoryCount int64
	u.db.WithContext(ctx).Model(&domain.PickupRequest{}).Count(&pickupCount)
	u.db.WithContext(ctx).Model(&domain.Inventory{}).Count(&inventoryCount)

	return inventoryCount > 0, map[string]int64{
		"pickup_requests": pickupCount,
		"inventories":      inventoryCount,
	}, nil
}

func (u *SystemUseCase) SeedData(ctx context.Context, force bool, usersMap map[string]string) (int64, error) {
	seeded, _, _ := u.GetStatus(ctx)
	if seeded && !force {
		return 0, nil
	}

	inventories := []domain.Inventory{
		{
			ID:                "INV_ARABICA_001",
			CoffeeType:        "ARABICA",
			OriginCode:        "VN-LD",
			WarehouseID:       "WAREHOUSE-HN-001",
			SKU:               "SKU-AR-VN-LD-001",
			AvailableQuantity: 1000.0,
		},
		{
			ID:                "INV_ROBUSTA_001",
			CoffeeType:        "ROBUSTA",
			OriginCode:        "VN-DL",
			WarehouseID:       "WAREHOUSE-HN-001",
			SKU:               "SKU-RB-VN-DL-001",
			AvailableQuantity: 2000.0,
		},
	}

	var created int64
	for _, inv := range inventories {
		if err := u.db.WithContext(ctx).Where("id = ?", inv.ID).FirstOrCreate(&inv).Error; err == nil {
			created++
		}
	}

	return created, nil
}
