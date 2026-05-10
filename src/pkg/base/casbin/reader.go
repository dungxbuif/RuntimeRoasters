package casbin

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

type ResilientReader struct {
	enforcer   *casbin.Enforcer
	authClient AuthSnapshotClient // Will be defined in transport/grpc
	mu         sync.RWMutex
	options    ReaderOptions
}

type AuthSnapshotClient interface {
	GetFullSnapshot(ctx context.Context) ([]string, error)
}

type ReaderOptions struct {
	ModelText     string
	SyncInterval  time.Duration
	RetryInterval time.Duration
}

func NewResilientReader(client AuthSnapshotClient, opts ReaderOptions) (*ResilientReader, error) {
	m, err := model.NewModelFromString(opts.ModelText)
	if err != nil {
		return nil, err
	}

	// Initialize with Memory adapter as rules come from gRPC/Kafka
	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, err
	}

	return &ResilientReader{
		enforcer:   e,
		authClient: client,
		options:    opts,
	}, nil
}

func (r *ResilientReader) Enforce(rvals ...interface{}) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.enforcer.Enforce(rvals...)
}

func (r *ResilientReader) GetEnforcer() *casbin.Enforcer {
	return r.enforcer
}

// Sync performs a full snapshot sync from the Auth Service
func (r *ResilientReader) Sync(ctx context.Context) error {
	log := logger.GetLogger()

	rules, err := r.authClient.GetFullSnapshot(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch auth snapshot: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Clear existing and load new
	r.enforcer.ClearPolicy()
	for _, line := range rules {
		parts := strings.Split(line, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}

		if len(parts) < 3 {
			continue
		}

		if parts[0] == "p" {
			_, _ = r.enforcer.AddPolicy(convertToInterface(parts[1:])...)
		} else if parts[0] == "g" {
			_, _ = r.enforcer.AddGroupingPolicy(convertToInterface(parts[1:])...)
		}
	}

	log.Info("Auth policies synchronized from snapshot", zap.Int("count", len(rules)))
	return nil
}

func convertToInterface(strs []string) []interface{} {
	res := make([]interface{}, len(strs))
	for i, s := range strs {
		res[i] = s
	}
	return res
}

// StartBackgroundSync starts the resilience pillars
func (r *ResilientReader) StartBackgroundSync(ctx context.Context) {
	go r.bootstrap(ctx)
	go r.pollingLoop(ctx)
}

func (r *ResilientReader) bootstrap(ctx context.Context) {
	log := logger.GetLogger()
	for {
		if err := r.Sync(ctx); err == nil {
			break
		}
		log.Warn("Bootstrap auth sync failed, retrying...", zap.Duration("delay", r.options.RetryInterval))
		time.Sleep(r.options.RetryInterval)
	}
}

func (r *ResilientReader) pollingLoop(ctx context.Context) {
	ticker := time.NewTicker(r.options.SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.Sync(ctx); err != nil {
				logger.GetLogger().Error("Background auth sync failed", zap.Error(err))
			}
		}
	}
}
