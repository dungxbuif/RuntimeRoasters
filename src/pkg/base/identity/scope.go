package identity

import "context"

const (
	RoleAdmin    = "ADMIN"
	RoleStoreMgr = "STORE_MGR"
)

func (c Claims) IsAdmin() bool {
	return c.Role == RoleAdmin
}

func (c Claims) CanAccessStore(storeID string) bool {
	if c.IsAdmin() {
		return true
	}
	if c.Role != RoleStoreMgr || storeID == "" {
		return false
	}
	for _, allowed := range c.StoreIDs {
		if allowed == storeID {
			return true
		}
	}
	return false
}

func (c Claims) ScopedStoreIDs() []string {
	if c.IsAdmin() {
		return nil
	}
	if c.Role != RoleStoreMgr || len(c.StoreIDs) == 0 {
		return []string{}
	}
	return append([]string(nil), c.StoreIDs...)
}

func StoreScopeFromContext(ctx context.Context) (Claims, []string, bool) {
	claims, ok := FromContext(ctx)
	if !ok {
		return Claims{}, nil, true
	}
	if claims.IsAdmin() {
		return claims, nil, true
	}
	if claims.Role != RoleStoreMgr || len(claims.StoreIDs) == 0 {
		return claims, []string{}, false
	}
	return claims, claims.StoreIDs, false
}
