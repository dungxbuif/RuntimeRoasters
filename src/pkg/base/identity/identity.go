package identity

type Claims struct {
	Subject      string
	Email        string
	Role         string
	OrgID        string
	StoreIDs     []string
	WarehouseIDs []string
	JTI          string
}
