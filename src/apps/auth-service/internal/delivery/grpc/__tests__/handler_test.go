package grpc_tests

import (
	"context"
	"testing"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/delivery/grpc"
	infra_casbin "github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFullSnapshot(t *testing.T) {
	// 1. Setup Casbin Enforcer in memory
	m, err := model.NewModelFromString(`
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
`)
	require.NoError(t, err)

	e, err := casbin.NewEnforcer(m)
	require.NoError(t, err)

	// Add some policies
	_, err = e.AddPolicy("admin", "farm", "write")
	require.NoError(t, err)
	_, err = e.AddGroupingPolicy("alice", "admin")
	require.NoError(t, err)

	infraEnforcer := &infra_casbin.Enforcer{Enforcer: e}
	handler := grpc.NewHandler(infraEnforcer)

	// 2. Test GetFullSnapshot
	resp, err := handler.GetFullSnapshot(context.Background(), &authv1.GetFullSnapshotRequest{})
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	
	// Verify policies are formatted as CSV
	// Expected: ["p, admin, farm, write", "g, alice, admin"]
	assert.Contains(t, resp.Policies, "p, admin, farm, write")
	assert.Contains(t, resp.Policies, "g, alice, admin")
	assert.Equal(t, 2, len(resp.Policies))
}
