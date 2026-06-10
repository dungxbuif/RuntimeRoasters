package identity

import "time"

type Claims struct {
	Subject      string
	Email        string
	Role         string
	OrgID        string
	StoreIDs     []string
	WarehouseIDs []string
	JTI          string
	IssuedAt     time.Time
	ExpiresAt    time.Time
}
