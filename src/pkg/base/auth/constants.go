package auth

const (
	HeaderAuthorization = "Authorization"
	BearerPrefix        = "bearer"
	GRPCAuthMetadataKey = "authorization"
)

const (
	JWTHeaderKeyID = "kid"
	JWTHeaderAlg   = "alg"
)

const (
	ClaimRole  = "role"
	ClaimOrgID = "org_id"
	ClaimJTI   = "jti"
)

const (
	TraceAttrUserID = "user.id"
)
