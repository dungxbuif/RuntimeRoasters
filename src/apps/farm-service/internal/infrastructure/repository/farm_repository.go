package repository

import (
	"context"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/farm-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/database"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"gorm.io/gorm"
)

type FarmModel struct {
	ID         string `gorm:"primaryKey"`
	Name       string
	Location   string
	Area       float64
	CoffeeType string
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

type FarmRepository interface {
	Create(ctx context.Context, farm *domain.Farm) error
	GetByID(ctx context.Context, id string) (*domain.Farm, error)
	List(ctx context.Context) ([]*domain.Farm, error)
	Update(ctx context.Context, farm *domain.Farm) error
	Delete(ctx context.Context, id string) error
}

type farmRepository struct {
	db *database.DB
}

func NewFarmRepository(db *database.DB) FarmRepository {
	return &farmRepository{db: db}
}

func (r *farmRepository) Create(ctx context.Context, farm *domain.Farm) error {
	id, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	// If caller is admin/farm_admin and OwnerID is explicitly provided, we honor it.
	// Otherwise, we force ownership to the caller.
	if (id.Role == "admin" || id.Role == "farm_admin") && farm.OwnerID != "" {
		// Use provided OwnerID
	} else {
		farm.OwnerID = id.Subject
	}

	model := FromDomain(farm)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *farmRepository) GetByID(ctx context.Context, id string) (*domain.Farm, error) {
	userId, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var model FarmModel
	err := r.db.WithContext(ctx).Where("id = ? AND owner_id = ?", id, userId.Subject).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *farmRepository) List(ctx context.Context) ([]*domain.Farm, error) {
	userId, ok := identity.FromContext(ctx)
	if !ok {
		return nil, errs.ErrUnauthorized
	}

	var models []FarmModel
	err := r.db.WithContext(ctx).Where("owner_id = ?", userId.Subject).Find(&models).Error
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
	userId, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	model := FromDomain(farm)
	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", model.ID, userId.Subject).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}

func (r *farmRepository) Delete(ctx context.Context, id string) error {
	userId, ok := identity.FromContext(ctx)
	if !ok {
		return errs.ErrUnauthorized
	}

	result := r.db.WithContext(ctx).
		Where("id = ? AND owner_id = ?", id, userId.Subject).
		Delete(&FarmModel{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errs.ErrNotFound
	}

	return nil
}
