package blacklist

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	jtiPrefix = "blacklist:jti:"
)

type RedisBlacklist struct {
	client *redis.Client
}

func NewRedisBlacklist(client *redis.Client) *RedisBlacklist {
	return &RedisBlacklist{client: client}
}

// IsRevoked checks if the JTI exists in the Redis blacklist.
func (r *RedisBlacklist) IsRevoked(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}
	key := fmt.Sprintf("%s%s", jtiPrefix, jti)
	val, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist in redis: %w", err)
	}
	return val > 0, nil
}

// Revoke adds a JTI to the Redis blacklist with a specific TTL.
func (r *RedisBlacklist) Revoke(ctx context.Context, jti string, ttl time.Duration) error {
	key := fmt.Sprintf("%s%s", jtiPrefix, jti)
	err := r.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke token in redis: %w", err)
	}
	return nil
}
