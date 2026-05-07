package identity

type Claims struct {
	Subject string
	Role    string
	OrgID   string
	JTI     string
}