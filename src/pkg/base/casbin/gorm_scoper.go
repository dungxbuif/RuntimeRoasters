package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// GormScoper provides a way to apply Casbin-based filters to GORM queries
type GormScoper struct {
	enforcer *casbin.SyncedEnforcer
}

func NewGormScoper(e *casbin.SyncedEnforcer) *GormScoper {
	return &GormScoper{enforcer: e}
}

// ApplyScope restricts a GORM query based on Casbin policies.
// It looks for policies that define constraints on the object.
func (s *GormScoper) ApplyScope(user, object, action string, ownerField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 1. Check if user has global access (*)
		allowed, err := s.enforcer.Enforce(user, object, action)
		if err != nil {
			return db.Where("1 = 0")
		}

		// 2. Determine IF we need to scope by owner.
		isAdmin, _ := s.enforcer.HasRoleForUser(user, "admin")
		isFarmAdmin, _ := s.enforcer.HasRoleForUser(user, "farm_admin")
		
		if allowed && !isAdmin && !isFarmAdmin {
			return db.Where(fmt.Sprintf("%s = ?", ownerField), user)
		}

		if allowed {
			return db
		}

		return db.Where("1 = 0")
	}
}

// NewGormAdapterEnforcer creates a local enforcer backed by a GORM table.
func NewGormAdapterEnforcer(db *gorm.DB, modelText string) (*casbin.SyncedEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return e, nil
}
