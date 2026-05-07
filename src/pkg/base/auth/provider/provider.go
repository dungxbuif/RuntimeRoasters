package provider

import "crypto/rsa"

/* Interface để tách logic get JWKS (từ cached. từ sources,...)
   Provider sẽ được implement bởi JWKS, cached, v.v.
*/
type KeyProvider interface {
	GetPublicKey(kid string) (*rsa.PublicKey, error)
}