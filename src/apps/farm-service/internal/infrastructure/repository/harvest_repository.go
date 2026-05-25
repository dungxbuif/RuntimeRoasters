package repository

import (
	"context"
	"time"

	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/apps/farm-service/internal/usecase"
	rr_casbin "RuntimeRoasters/pkg/base/casbin"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/errs"
)

type HarvestModel struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	FarmID      uint64 `gorm:"not null"`
	OwnerID     string `gorm:"type:uuid;not null;index"`
	CoffeeType  domain.CoffeeType
	Quantity    float64 `gorm:"type:decimal(10,2);not null"`
	HarvestDate time.Time
	Status      domain.HarvestStatus
	Notes       string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (HarvestModel) TableName() string {
	return "harvests"
}

type harvestRepository struct {
	db     *database.DB
	scoper *rr_casbin.GormScoper
}

func NewHarvestRepository(db *database.DB, enforcer rr_casbin.PolicyEnforcer) usecase.HarvestRepository {
	return &harvestRepository{
		db:     db,
		scoper: rr_casbin.NewGormScoper(enforcer),
	}
}

func (r *harvestRepository) Create(ctx context.Context, harvest *domain.Harvest) error {
	model := toHarvestModel(harvest)
	if err := r.db.GetTx(ctx).Create(model).Error; err != nil {
		return err
	}
	harvest.ID = model.ID
	return nil
}

func (r *harvestRepository) GetByID(ctx context.Context, id uint64) (*domain.Harvest, error) {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var model HarvestModel
	err := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "harvest", "read", "owner_id")).
		First(&model, id).Error

	if err != nil {
		return nil, err
	}
	return toHarvestDomain(&model), nil
}

func (r *harvestRepository) ListByFarm(ctx context.Context, farmID uint64) ([]*domain.Harvest, error) {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var models []HarvestModel
	err := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "harvest", "read", "owner_id")).
		Where("farm_id = ?", farmID).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	harvests := make([]*domain.Harvest, len(models))
	for i, m := range models {
		harvests[i] = toHarvestDomain(&m)
	}
	return harvests, nil
}

func (r *harvestRepository) Update(ctx context.Context, harvest *domain.Harvest) error {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	model := toHarvestModel(harvest)
	result := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "harvest", "write", "owner_id")).
		Where("id = ?", model.ID).
		Updates(model)
	return result.Error
}

func (r *harvestRepository) Delete(ctx context.Context, id uint64) error {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	result := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "harvest", "delete", "owner_id")).
		Where("id = ?", id).
		Delete(&HarvestModel{})
	return result.Error
}

func toHarvestModel(h *domain.Harvest) *HarvestModel {
	return &HarvestModel{
		ID:          h.ID,
		FarmID:      h.FarmID,
		OwnerID:     h.OwnerID,
		CoffeeType:  h.CoffeeType,
		Quantity:    h.Quantity,
		HarvestDate: h.HarvestDate,
		Status:      h.Status,
		Notes:       h.Notes,
		CreatedAt:   h.CreatedAt,
		UpdatedAt:   h.UpdatedAt,
	}
}

func toHarvestDomain(m *HarvestModel) *domain.Harvest {
	return &domain.Harvest{
		ID:          m.ID,
		FarmID:      m.FarmID,
		OwnerID:     m.OwnerID,
		CoffeeType:  m.CoffeeType,
		Quantity:    m.Quantity,
		HarvestDate: m.HarvestDate,
		Status:      m.Status,
		Notes:       m.Notes,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
