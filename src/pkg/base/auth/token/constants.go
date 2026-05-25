package token

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
	ClaimEmail        = "email"
	ClaimRole         = "role"
	ClaimOrgID        = "org_id"
	ClaimStoreIDs     = "store_ids"
	ClaimWarehouseIDs = "warehouse_ids"
	ClaimJTI          = "jti"
)

const (
	TraceAttrUserID = "user.id"
)
