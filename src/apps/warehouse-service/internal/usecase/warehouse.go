package usecase

import (
	"context"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WarehouseUseCase struct {
	db *gorm.DB
}

func NewWarehouseUseCase(db *gorm.DB) *WarehouseUseCase {
	return &WarehouseUseCase{db: db}
}

func (u *WarehouseUseCase) CreateWarehouse(ctx context.Context, wh domain.Warehouse) (domain.Warehouse, error) {
	if wh.ID == "" {
		wh.ID = uuid.NewString()
	}
	now := time.Now()
	wh.CreatedAt = now
	wh.UpdatedAt = now
	
	err := u.db.WithContext(ctx).Where("code = ?", wh.Code).FirstOrCreate(&wh).Error
	return wh, err
}

func (u *WarehouseUseCase) ListWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	var warehouses []domain.Warehouse
	err := u.db.WithContext(ctx).Order("created_at DESC").Find(&warehouses).Error
	return warehouses, err
}

func (u *WarehouseUseCase) GetWarehouse(ctx context.Context, id string) (domain.Warehouse, error) {
	var wh domain.Warehouse
	err := u.db.WithContext(ctx).Where("id = ? OR code = ?", id, id).First(&wh).Error
	return wh, err
}

func (u *WarehouseUseCase) UpdateWarehouse(ctx context.Context, wh domain.Warehouse) (domain.Warehouse, error) {
	wh.UpdatedAt = time.Now()
	err := u.db.WithContext(ctx).Save(&wh).Error
	return wh, err
}

func (u *WarehouseUseCase) DeleteWarehouse(ctx context.Context, id string) error {
	return u.db.WithContext(ctx).Delete(&domain.Warehouse{}, "id = ?", id).Error
}
