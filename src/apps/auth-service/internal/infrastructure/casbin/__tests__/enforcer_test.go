package casbin_tests

import (
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	casbininfra "RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
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
	original := casbininfra.DefaultPoliciesCSV
	t.Cleanup(func() {
		casbininfra.DefaultPoliciesCSV = original
	})
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

func TestDefaultPoliciesCoverRepositoryScopes(t *testing.T) {
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
	e, _ := casbin.NewEnforcer(m)

	original := casbininfra.DefaultPoliciesCSV
	t.Cleanup(func() {
		casbininfra.DefaultPoliciesCSV = original
	})
	casbininfra.DefaultPoliciesCSV = original

	if err := casbininfra.SyncDefaultPolicies(e); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	cases := []struct {
		sub string
		obj string
		act string
	}{
		{sub: "FARM_ADMIN", obj: "farm", act: "read"},
		{sub: "FARM_ADMIN", obj: "harvest", act: "write"},
		{sub: "FARM_MANAGER", obj: "farm", act: "read"},
		{sub: "FARM_MANAGER", obj: "farm", act: "write"},
		{sub: "FARM_MANAGER", obj: "farm", act: "delete"},
		{sub: "FARM_MANAGER", obj: "harvest", act: "read"},
	}

	for _, tc := range cases {
		ok, _ := e.Enforce(tc.sub, tc.obj, tc.act)
		if !ok {
			t.Fatalf("expected policy for %s %s %s", tc.sub, tc.obj, tc.act)
		}
	}

	ok, _ := e.Enforce("ADMIN", "farm", "read")
	if !ok {
		t.Fatal("expected ADMIN to inherit farm read through FARM_ADMIN")
	}
}
