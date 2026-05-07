package auth

import "time"

type Config struct {
	JWKSURL        string
	InternalSecret string
	CacheTTL       time.Duration
	ExpectedIssuer string
}