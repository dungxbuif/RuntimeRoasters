package repository

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/usecase"
	rr_casbin "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"gorm.io/gorm"
)

type FarmModel struct {
	ID         uint64 `gorm:"primaryKey;autoIncrement"`
	Name       string
	Location   domain.Location
	Area       float64
	CoffeeType domain.CoffeeType
	OwnerID    string `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (FarmModel) TableName() string {
	return "farms"
}

func (m *FarmModel) ToDomain() *domain.Farm {
	return &domain.Farm{
		ID:         m.ID,
		Name:       m.Name,
		Location:   m.Location,
		Area:       m.Area,
		CoffeeType: m.CoffeeType,
		OwnerID:    m.OwnerID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func FromDomain(f *domain.Farm) *FarmModel {
	return &FarmModel{
		ID:         f.ID,
		Name:       f.Name,
		Location:   f.Location,
		Area:       f.Area,
		CoffeeType: f.CoffeeType,
		OwnerID:    f.OwnerID,
		CreatedAt:  f.CreatedAt,
		UpdatedAt:  f.UpdatedAt,
	}
}

type farmRepository struct {
	db     *database.DB
	scoper *rr_casbin.GormScoper
}

func NewFarmRepository(db *database.DB, enforcer rr_casbin.PolicyEnforcer) usecase.FarmRepository {
	return &farmRepository{
		db:     db,
		scoper: rr_casbin.NewGormScoper(enforcer),
	}
}

func (r *farmRepository) Create(ctx context.Context, farm *domain.Farm) error {
	model := FromDomain(farm)
	err := r.db.WithContext(ctx).Create(model).Error
	if err == nil {
		farm.ID = model.ID
	}
	return err
}

func (r *farmRepository) GetByID(ctx context.Context, id uint64) (*domain.Farm, error) {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var model FarmModel
	err := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "farm", "read", "owner_id")).
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *farmRepository) List(ctx context.Context) ([]*domain.Farm, error) {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var models []FarmModel
	err := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "farm", "read", "owner_id")).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	farms := make([]*domain.Farm, len(models))
	for i, m := range models {
		farms[i] = m.ToDomain()
	}

	return farms, nil
}

func (r *farmRepository) Update(ctx context.Context, farm *domain.Farm) error {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	model := FromDomain(farm)
	result := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "farm", "write", "owner_id")).
		Where("id = ?", model.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *farmRepository) Delete(ctx context.Context, id uint64) error {
	user, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	result := r.db.WithContext(ctx).
		Scopes(r.scoper.ApplyScope(user.Subject, user.Role, "farm", "delete", "owner_id")).
		Where("id = ?", id).
		Delete(&FarmModel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
