package token

import (
	"context"
	"time"
)

// TokenValidityChecker defines the interface for verifying user token validity boundary.
type TokenValidityChecker interface {
	IsValid(ctx context.Context, subject string, iat time.Time) (bool, error)
}
