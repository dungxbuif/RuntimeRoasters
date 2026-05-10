package casbin

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

//go:embed model.conf
var modelConf string

//go:embed default_policies.csv
var DefaultPoliciesCSV string

type Enforcer struct {
	*casbin.Enforcer
}

// NewEnforcer initializes the centralized enforcer for auth-service (The Writer).
// It uses GORM adapter to persist rules in Postgres and seeds default policies.
func NewEnforcer(dbURL string) (*Enforcer, error) {
	log := logger.GetLogger()

	// 1. Initialize GORM Adapter
	adapter, err := gormadapter.NewAdapter("postgres", dbURL, true)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize casbin gorm adapter: %w", err)
	}

	// 2. Load Model from embedded string
	m, err := model.NewModelFromString(modelConf)
	if err != nil {
		return nil, fmt.Errorf("failed to load casbin model: %w", err)
	}

	// 3. Create Enforcer
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 4. Load policies from DB
	if err := e.LoadPolicy(); err != nil {
		log.Error("failed to load policies from database", zap.Error(err))
	}

	// 5. Sync Default Policies (Auto-Seed)
	if err := SyncDefaultPolicies(e); err != nil {
		log.Error("failed to sync default policies", zap.Error(err))
	}

	log.Info("Casbin Enforcer initialized successfully (Centralized Writer)")
	return &Enforcer{e}, nil
}

// SyncDefaultPolicies parses the embedded CSV and seeds missing rules into the DB
func SyncDefaultPolicies(e *casbin.Enforcer) error {
	lines := strings.Split(strings.TrimSpace(DefaultPoliciesCSV), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		if len(parts) < 3 {
			continue
		}

		ptype := parts[0]
		params := convertToInterface(parts[1:])

		var exists bool
		var err error

		if ptype == "p" {
			exists, err = e.HasPolicy(params...)
			if err == nil && !exists {
				_, err = e.AddPolicy(params...)
			}
		} else if ptype == "g" {
			exists, err = e.HasGroupingPolicy(params...)
			if err == nil && !exists {
				_, err = e.AddGroupingPolicy(params...)
			}
		}

		if err != nil {
			return fmt.Errorf("failed to sync rule [%s]: %w", line, err)
		}
	}
	return nil
}

func convertToInterface(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, s := range strs {
		res[i] = s
	}
	return res
}
