package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	authv1 "RuntimeRoasters/runtime/auth/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type securityTokenValidityChecker struct {
	redisClient *redis.Client
	authClient  authv1.AuthServiceClient
	conn        *grpc.ClientConn
}

type CachedState struct {
	ValidAfter int64 `json:"valid_after"`
	Version    int64 `json:"version"`
}

// NewTokenValidityChecker creates a TokenValidityChecker that reads from Valkey and falls back to auth-service via gRPC.
func NewTokenValidityChecker(redisClient *redis.Client, authServiceAddr string) (TokenValidityChecker, error) {
	conn, err := grpc.Dial(authServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth service: %w", err)
	}
	return &securityTokenValidityChecker{
		redisClient: redisClient,
		authClient:  authv1.NewAuthServiceClient(conn),
		conn:        conn,
	}, nil
}

func (c *securityTokenValidityChecker) getCacheKey(subject string) string {
	h := sha256.New()
	h.Write([]byte(subject))
	hashed := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("auth:valid-after:v1:%s", hashed)
}

func (c *securityTokenValidityChecker) IsValid(ctx context.Context, subject string, iat time.Time) (bool, error) {
	cacheKey := c.getCacheKey(subject)

	// 1. Check Valkey Cache
	if c.redisClient != nil {
		val, err := c.redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cached CachedState
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				return iat.Unix() > cached.ValidAfter, nil
			}
		} else if !errors.Is(err, redis.Nil) {
			// Fail-open cache read, fallback to gRPC
		}
	}

	// 2. Fallback to auth-service gRPC lookup
	gCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.authClient.GetTokenValidityBoundary(gCtx, &authv1.GetTokenValidityBoundaryRequest{
		Subject: subject,
	})
	if err != nil {
		// Fail closed on gRPC error
		return false, fmt.Errorf("auth service token boundary lookup failed: %w", err)
	}

	validAfterTime, err := time.Parse(time.RFC3339, resp.TokensValidAfter)
	if err != nil {
		return false, fmt.Errorf("invalid boundary timestamp format: %w", err)
	}

	return iat.Unix() > validAfterTime.Unix(), nil
}
