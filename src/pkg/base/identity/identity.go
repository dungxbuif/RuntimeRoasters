package identity

type Claims struct {
	Subject  string
	Email    string
	Role     string
	OrgID    string
	StoreIDs []string
	JTI      string
}
