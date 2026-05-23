package casbin

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// GormScoper provides a way to apply Casbin-based filters to GORM queries
type PolicyEnforcer interface {
	Enforce(rvals ...interface{}) (bool, error)
}

type GormScoper struct {
	enforcer PolicyEnforcer
}

func NewGormScoper(e PolicyEnforcer) *GormScoper {
	return &GormScoper{enforcer: e}
}

// ApplyScope restricts a GORM query based on Casbin policies.
// It looks for policies that define constraints on the object.
func (s *GormScoper) ApplyScope(subject, role, object, action string, ownerField string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		normalizedRole := strings.ToUpper(role)
		if normalizedRole == "ADMIN" || normalizedRole == "FARM_ADMIN" {
			return db
		}

		// 1. Check if user has global access (*)
		allowed, err := s.enforcer.Enforce(role, object, action)
		if err != nil {
			return db.Where("1 = 0")
		}

		// 2. Determine IF we need to scope by owner.
		if allowed {
			return db.Where(fmt.Sprintf("%s = ?", ownerField), subject)
		}

		return db.Where("1 = 0")
	}
}

// NewGormAdapterEnforcer creates a local enforcer backed by a GORM table.
// modelSource can be a path to a .conf file or the raw model text.
func NewGormAdapterEnforcer(db *gorm.DB, modelSource string) (*casbin.SyncedEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	var m model.Model
	if len(modelSource) < 255 && (modelSource[len(modelSource)-5:] == ".conf" || modelSource[len(modelSource)-4:] == ".ini") {
		m, err = model.NewModelFromFile(modelSource)
	} else {
		m, err = model.NewModelFromString(modelSource)
	}

	if err != nil {
		return nil, err
	}

	e, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return e, nil
}
