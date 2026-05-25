package identity_tests

import (
	"testing"

	"RuntimeRoasters/pkg/base/identity"
)

func TestClaimsCanAccessStore(t *testing.T) {
	tests := []struct {
		name    string
		claims  identity.Claims
		storeID string
		want    bool
	}{
		{name: "admin sees all records", claims: identity.Claims{Role: identity.RoleAdmin}, storeID: "store-a", want: true},
		{name: "store manager sees assigned store", claims: identity.Claims{Role: identity.RoleStoreMgr, StoreIDs: []string{"store-a"}}, storeID: "store-a", want: true},
		{name: "store manager cannot see other stores", claims: identity.Claims{Role: identity.RoleStoreMgr, StoreIDs: []string{"store-a"}}, storeID: "store-b", want: false},
		{name: "store manager without stores fails closed", claims: identity.Claims{Role: identity.RoleStoreMgr}, storeID: "store-a", want: false},
		{name: "records without store are admin only", claims: identity.Claims{Role: identity.RoleStoreMgr, StoreIDs: []string{"store-a"}}, storeID: "", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.claims.CanAccessStore(tt.storeID); got != tt.want {
				t.Fatalf("CanAccessStore() = %v, want %v", got, tt.want)
			}
		})
	}
}
