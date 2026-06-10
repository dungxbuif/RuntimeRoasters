package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/pkg/database"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	"github.com/ory/client-go"
	hydra "github.com/ory/hydra-client-go/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type SecurityUsecase interface {
	RevokeAllSessions(ctx context.Context, subject string, reason string, updatedBy string) (*domain.UserSecurityState, error)
	ValidateTokenValidity(ctx context.Context, subject string, iat time.Time) (bool, error)
	GetTokensValidAfter(ctx context.Context, subject string) (time.Time, error)
}

type securityUsecase struct {
	cfg          *config.Config
	repo         domain.SecurityStateRepository
	tx           database.TxManager
	redis        *redis.Client
	kratosClient *client.APIClient
	hydraClient  *hydra.APIClient
	producer     kafka.Producer
	sfGroup      singleflight.Group
}

type CachedState struct {
	ValidAfter int64 `json:"valid_after"`
	Version    int64 `json:"version"`
}

func NewSecurityUsecase(
	cfg *config.Config,
	repo domain.SecurityStateRepository,
	tx database.TxManager,
	rClient *redis.Client,
	producer kafka.Producer,
) SecurityUsecase {
	kratosCfg := client.NewConfiguration()
	kratosCfg.Servers = client.ServerConfigurations{{URL: cfg.KratosAdminURL}}
	if cfg.InternalSecret != "" {
		kratosCfg.DefaultHeader["X-Internal-Secret"] = cfg.InternalSecret
	}
	kratosClient := client.NewAPIClient(kratosCfg)

	hydraCfg := hydra.NewConfiguration()
	hydraCfg.Servers = hydra.ServerConfigurations{{URL: cfg.HydraAdminURL}}
	hydraClient := hydra.NewAPIClient(hydraCfg)

	return &securityUsecase{
		cfg:          cfg,
		repo:         repo,
		tx:           tx,
		redis:        rClient,
		kratosClient: kratosClient,
		hydraClient:  hydraClient,
		producer:     producer,
	}
}

func (u *securityUsecase) getCacheKey(subject string) string {
	h := sha256.New()
	h.Write([]byte(subject))
	hashed := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("auth:valid-after:v1:%s", hashed)
}

func (u *securityUsecase) RevokeAllSessions(ctx context.Context, subject string, reason string, updatedBy string) (*domain.UserSecurityState, error) {
	log := logger.FromContext(ctx).With(zap.String("subject", subject))
	log.Info("Revoking all sessions for subject", zap.String("reason", reason), zap.String("updated_by", updatedBy))

	var state *domain.UserSecurityState
	err := u.tx.WithTx(ctx, func(txCtx context.Context) error {
		var txErr error
		state, txErr = u.repo.MonotonicUpdate(txCtx, subject, reason, updatedBy)
		if txErr != nil {
			return fmt.Errorf("failed database revocation update: %w", txErr)
		}
		return nil
	})

	if err != nil {
		log.Error("Failed to commit database revocation boundary", zap.Error(err))
		return nil, err
	}

	// Write-through to Valkey cache
	cacheKey := u.getCacheKey(subject)
	payload := CachedState{
		ValidAfter: state.TokensValidAfter.Unix(),
		Version:    state.Version,
	}
	payloadJSON, _ := json.Marshal(payload)

	// Configurable TTL with jitter (e.g. 15 minutes + 0 to 60 seconds jitter)
	jitter := time.Duration(rand.Intn(60)) * time.Second
	ttl := 15*time.Minute + jitter

	if u.redis != nil {
		if err := u.redis.Set(ctx, cacheKey, payloadJSON, ttl).Err(); err != nil {
			log.Warn("Failed to update Valkey cache on revocation", zap.Error(err))
		} else {
			log.Debug("Successfully updated Valkey revocation cache")
		}
	}

	// Invalidate Kratos session
	log.Debug("Deleting Kratos identity sessions")
	_, kratosErr := u.kratosClient.IdentityAPI.DeleteIdentitySessions(ctx, subject).Execute()
	if kratosErr != nil {
		log.Warn("Failed to invalidate Kratos sessions", zap.Error(kratosErr))
	}

	// Invalidate Hydra login & consent sessions
	log.Debug("Revoking Hydra login sessions")
	_, hydraLoginErr := u.hydraClient.OAuth2API.RevokeOAuth2LoginSessions(ctx).Subject(subject).Execute()
	if hydraLoginErr != nil {
		log.Warn("Failed to invalidate Hydra login sessions", zap.Error(hydraLoginErr))
	}

	log.Debug("Revoking Hydra consent sessions")
	_, hydraConsentErr := u.hydraClient.OAuth2API.RevokeOAuth2ConsentSessions(ctx).Subject(subject).All(true).Execute()
	if hydraConsentErr != nil {
		log.Warn("Failed to invalidate Hydra consent sessions", zap.Error(hydraConsentErr))
	}

	// Publish security audit event to Kafka
	u.publishRevocationEvent(ctx, state, reason)

	return state, nil
}

func (u *securityUsecase) ValidateTokenValidity(ctx context.Context, subject string, iat time.Time) (bool, error) {
	boundary, err := u.GetTokensValidAfter(ctx, subject)
	if err != nil {
		return false, err
	}
	return iat.Unix() > boundary.Unix(), nil
}

func (u *securityUsecase) GetTokensValidAfter(ctx context.Context, subject string) (time.Time, error) {
	log := logger.FromContext(ctx).With(zap.String("subject", subject))

	cacheKey := u.getCacheKey(subject)

	// 1. Valkey Cache Read
	if u.redis != nil {
		val, err := u.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var cached CachedState
			if err := json.Unmarshal([]byte(val), &cached); err == nil {
				log.Debug("Valkey revocation cache hit", zap.Int64("valid_after", cached.ValidAfter))
				return time.Unix(cached.ValidAfter, 0), nil
			}
		} else if !errors.Is(err, redis.Nil) {
			log.Warn("Valkey cache read error", zap.Error(err))
		}
	}

	// 2. Cache Miss fallback with singleflight coalescing
	res, err, _ := u.sfGroup.Do(subject, func() (interface{}, error) {
		// Read from database
		state, dbErr := u.repo.GetBySubject(ctx, subject)
		if dbErr != nil {
			return nil, fmt.Errorf("authoritative DB lookup failed: %w", dbErr)
		}

		// Repopulate Valkey
		if u.redis != nil {
			payload := CachedState{
				ValidAfter: state.TokensValidAfter.Unix(),
				Version:    state.Version,
			}
			payloadJSON, _ := json.Marshal(payload)
			jitter := time.Duration(rand.Intn(60)) * time.Second
			ttl := 15*time.Minute + jitter
			if err := u.redis.Set(ctx, cacheKey, payloadJSON, ttl).Err(); err != nil {
				log.Warn("Failed to update Valkey cache after miss", zap.Error(err))
			}
		}

		return state.TokensValidAfter, nil
	})

	if err != nil {
		// Fail closed on dependency errors
		log.Error("Get tokens valid after failed closed due to database lookup error", zap.Error(err))
		return time.Time{}, err
	}

	return res.(time.Time), nil
}

func (u *securityUsecase) publishRevocationEvent(ctx context.Context, state *domain.UserSecurityState, reason string) {
	if u.producer == nil {
		return
	}

	event := map[string]interface{}{
		"event_id":           uuid.New().String(),
		"event_type":         "security.token_revocation.completed",
		"occurred_at":        time.Now().Format(time.RFC3339),
		"subject":            state.Subject,
		"scope":              "ALL_SUBJECT_SESSIONS",
		"reason":             reason,
		"tokens_valid_after": state.TokensValidAfter.Format(time.RFC3339),
		"version":            state.Version,
		"source":             "auth-service",
	}

	if err := u.producer.Publish(ctx, u.cfg.KafkaTopic, state.Subject, event); err != nil {
		logger.FromContext(ctx).Warn("Failed to publish revocation audit event to Kafka", zap.Error(err))
	} else {
		logger.FromContext(ctx).Info("Published security.token_revocation.completed audit event to Kafka")
	}
}
