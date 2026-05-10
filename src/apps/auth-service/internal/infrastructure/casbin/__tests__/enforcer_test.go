package casbin_tests

import (
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	casbininfra "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
)

func TestSyncDefaultPolicies(t *testing.T) {
	// 1. Setup a memory-based enforcer for speed
	// Use embedded model text for testing
	m, _ := model.NewModelFromString(`
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`)
	e, _ := casbin.NewEnforcer(m) // Memory adapter by default

	// 2. Override the embedded CSV for testing
	casbininfra.DefaultPoliciesCSV = `
p, tester, /test, read
g, lead_tester, tester
`

	// 3. Run Sync
	if err := casbininfra.SyncDefaultPolicies(e); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	// 4. Verify policies
	t.Run("Verify p policy", func(t *testing.T) {
		ok, _ := e.Enforce("tester", "/test", "read")
		if !ok {
			t.Error("expected tester to have read access to /test")
		}
	})

	t.Run("Verify g policy", func(t *testing.T) {
		ok, _ := e.Enforce("lead_tester", "/test", "read")
		if !ok {
			t.Error("expected lead_tester to inherit access")
		}
	})

	// 5. Run Sync again - should be idempotent (no duplicates)
	if err := casbininfra.SyncDefaultPolicies(e); err != nil {
		t.Fatalf("second sync failed: %v", err)
	}

	policies, _ := e.GetPolicy()
	if len(policies) != 1 {
		t.Errorf("expected 1 policy rule after idempotent sync, got %d", len(policies))
	}
}
