package token

import "context"

// BlacklistChecker defines the interface for checking if a token (JTI) has been revoked.
type BlacklistChecker interface {
	IsRevoked(ctx context.Context, jti string) (bool, error)
}
