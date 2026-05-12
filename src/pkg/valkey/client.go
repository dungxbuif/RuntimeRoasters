package valkey

import "github.com/redis/go-redis/v9"

type Config struct {
	Addr          string   // single node (dev): "localhost:6379"
	SentinelAddrs []string // sentinel nodes (prod): ["s1:26379", "s2:26379"]
	MasterName    string   // sentinel master name: "mymaster"
	Password      string
}

// NewClient returns a redis.Client (compatible with Valkey) configured for single node or Sentinel (HA) mode.
func NewClient(cfg Config) *redis.Client {
	if len(cfg.SentinelAddrs) > 0 {
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			SentinelAddrs: cfg.SentinelAddrs,
			Password:      cfg.Password,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
	})
}
