package casbin_tests

import (
	"context"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
)

// MockAuthClient implements AuthSnapshotClient for testing
type MockAuthClient struct {
	policies []string
}

func (m *MockAuthClient) GetFullSnapshot(ctx context.Context) ([]string, error) {
	return m.policies, nil
}

func TestResilientReader_Sync(t *testing.T) {
	modelText := `
[request_definition]
r = sub, obj, act
[policy_definition]
p = sub, obj, act
[role_definition]
g = _, _
[policy_effect]
e = some(where (p.eft == allow))
[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	mock := &MockAuthClient{
		policies: []string{
			"p, FARM_MANAGER, /v1/demo, read",
			"g, ADMIN, FARM_MANAGER",
		},
	}

	reader, err := casbin.NewResilientReader(mock, casbin.ReaderOptions{
		ModelText:     modelText,
		SyncInterval:  1 * time.Minute,
		RetryInterval: 1 * time.Second,
	})
	if err != nil {
		t.Fatalf("failed to create reader: %v", err)
	}

	// 1. Initial state should be empty
	allowed, _ := reader.Enforce("FARM_MANAGER", "/v1/demo", "read")
	if allowed {
		t.Error("expected access denied before sync")
	}

	// 2. Perform Sync
	if err := reader.Sync(context.Background()); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	// 3. Check permissions after sync
	t.Run("Check direct policy", func(t *testing.T) {
		allowed, _ := reader.Enforce("FARM_MANAGER", "/v1/demo", "read")
		if !allowed {
			t.Error("expected access granted for farm manager after sync")
		}
	})

	t.Run("Check role hierarchy", func(t *testing.T) {
		allowed, _ := reader.Enforce("ADMIN", "/v1/demo", "read")
		if !allowed {
			t.Error("expected access granted for admin via hierarchy")
		}
	})

	t.Run("Check unauthorized", func(t *testing.T) {
		allowed, _ := reader.Enforce("guest", "/v1/demo", "read")
		if allowed {
			t.Error("expected access denied for guest")
		}
	})
}
