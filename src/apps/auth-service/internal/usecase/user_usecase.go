package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	"github.com/ory/client-go"
	hydra "github.com/ory/hydra-client-go/v2"
	"go.uber.org/zap"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	AcceptHydraLogin(ctx context.Context, req domain.AcceptLoginRequest) (*domain.AcceptLoginResponse, error)
	SyncCasbinWithKratos(ctx context.Context) error
}

type userUsecase struct {
	cfg          *config.Config
	enforcer     *casbin.Enforcer
	kratosClient *client.APIClient
	hydraClient  *hydra.APIClient
	producer     kafka.Producer
}

func NewUserUsecase(cfg *config.Config, enforcer *casbin.Enforcer, producer kafka.Producer) UserUsecase {
	kratosCfg := client.NewConfiguration()
	kratosCfg.Servers = client.ServerConfigurations{{URL: cfg.KratosAdminURL}}
	if cfg.InternalSecret != "" {
		kratosCfg.DefaultHeader["X-Internal-Secret"] = cfg.InternalSecret
	}
	kratosClient := client.NewAPIClient(kratosCfg)

	hydraCfg := hydra.NewConfiguration()
	// Using direct localhost port for dev environment as per .env standard
	hydraAdminURL := "http://localhost:4445"
	hydraCfg.Servers = hydra.ServerConfigurations{{URL: hydraAdminURL}}
	// Hydra Admin 4445 usually doesn't have Nginx proxy in this setup, but let's be consistent if needed.
	// In docker-compose, hydra 4445 is exposed directly.
	hydraClient := hydra.NewAPIClient(hydraCfg)

	return &userUsecase{
		cfg:          cfg,
		enforcer:     enforcer,
		kratosClient: kratosClient,
		hydraClient:  hydraClient,
		producer:     producer,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Info("Creating new user", zap.String("email", req.Email), zap.String("role", req.Role))

	identityBody := *client.NewCreateIdentityBody(
		"default",
		map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
			"role":  req.Role,
		},
	)
	identityBody.Credentials = &client.IdentityWithCredentials{
		Password: &client.IdentityWithCredentialsPassword{
			Config: &client.IdentityWithCredentialsPasswordConfig{
				Password: &req.Password,
			},
		},
	}

	createdIdentity, _, err := u.kratosClient.IdentityAPI.CreateIdentity(ctx).CreateIdentityBody(identityBody).Execute()
	if err != nil {
		log.Error("failed to create identity in kratos", zap.Error(err))
		return nil, fmt.Errorf("failed to create identity in kratos: %w", err)
	}

	userId := createdIdentity.Id
	log.Debug("Identity created in Kratos", zap.String("kratos_id", userId))

	_, err = u.enforcer.AddGroupingPolicy(userId, req.Role)
	if err != nil {
		log.Error("failed to assign role in casbin", zap.Error(err))
		return nil, fmt.Errorf("failed to assign role in casbin: %w", err)
	}
	log.Debug("Role assigned in Casbin", zap.String("role", req.Role))

	user := &domain.User{
		ID:    userId,
		Email: req.Email,
		Name:  req.Name,
		Role:  req.Role,
	}

	event := map[string]interface{}{
		"event_id":    uuid.New().String(),
		"event_type":  "USER_CREATED",
		"payload":     user,
		"occurred_at": time.Now().Format(time.RFC3339),
	}

	err = u.producer.Publish(ctx, u.cfg.KafkaTopic, userId, event)
	if err != nil {
		log.Warn("failed to publish user.created event", zap.Error(err))
	} else {
		log.Info("Published user.created event to Kafka", zap.String("topic", u.cfg.KafkaTopic))
	}

	return user, nil
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	log := logger.FromContext(ctx)
	log.Debug("Listing users from Kratos")

	identities, _, err := u.kratosClient.IdentityAPI.ListIdentities(ctx).Execute()
	if err != nil {
		log.Error("failed to list identities from kratos", zap.Error(err))
		return nil, fmt.Errorf("failed to list identities from kratos: %w", err)
	}

	log.Debug("Successfully fetched identities", zap.Int("count", len(identities)))

	users := make([]*domain.User, len(identities))
	for i, id := range identities {
		traits := id.Traits.(map[string]interface{})
		email, _ := traits["email"].(string)
		name, _ := traits["name"].(string)

		roles, _ := u.enforcer.GetRolesForUser(id.Id)
		role := ""
		if len(roles) > 0 {
			role = roles[0]
		} else {
			// Fallback to traits if not found in Casbin (useful for seeded users)
			if r, ok := traits["role"].(string); ok {
				role = r
				log.Debug("Role not found in Casbin, falling back to Kratos trait", zap.String("user_id", id.Id), zap.String("role", role))
			} else {
				role = "GUEST"
			}
		}

		users[i] = &domain.User{
			ID:    id.Id,
			Email: email,
			Name:  name,
			Role:  role,
		}
	}

	return users, nil
}

func (u *userUsecase) AcceptHydraLogin(ctx context.Context, req domain.AcceptLoginRequest) (*domain.AcceptLoginResponse, error) {
	log := logger.FromContext(ctx)
	log.Info("Accepting Hydra login request", zap.String("subject", req.Subject))

	// 1. Fetch Role (Priority: Casbin -> Kratos Traits)
	roles, _ := u.enforcer.GetRolesForUser(req.Subject)
	role := ""
	if len(roles) > 0 {
		role = roles[0]
	} else {
		// Fetch from Kratos to find the trait
		id, _, err := u.kratosClient.IdentityAPI.GetIdentity(ctx, req.Subject).Execute()
		if err == nil {
			traits := id.Traits.(map[string]interface{})
			if r, ok := traits["role"].(string); ok {
				role = r
			}
		}
	}

	if role == "" {
		role = "GUEST"
	}

	log.Info("Mapping role to Hydra session", zap.String("subject", req.Subject), zap.String("role", role))

	// 2. Accept Login with Session Claims
	accept := *hydra.NewAcceptOAuth2LoginRequest(req.Subject)
	/* 
	   TODO: Fix custom claims for Hydra v2 Go SDK. 
	   In v2, custom claims might need to be passed differently.
	   For now, we just accept the login with the subject.
	*/

	res, _, err := u.hydraClient.OAuth2API.AcceptOAuth2LoginRequest(ctx).
		LoginChallenge(req.LoginChallenge).
		AcceptOAuth2LoginRequest(accept).
		Execute()

	if err != nil {
		log.Error("failed to accept hydra login", zap.Error(err))
		return nil, fmt.Errorf("failed to accept hydra login: %w", err)
	}

	log.Info("Hydra login accepted", zap.String("redirect_to", res.RedirectTo))
	return &domain.AcceptLoginResponse{
		RedirectTo: res.RedirectTo,
	}, nil
}

func (u *userUsecase) SyncCasbinWithKratos(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("Starting Casbin synchronization with Kratos identities")

	identities, _, err := u.kratosClient.IdentityAPI.ListIdentities(ctx).Execute()
	if err != nil {
		return fmt.Errorf("failed to list identities from kratos: %w", err)
	}

	syncedCount := 0
	for _, id := range identities {
		traits := id.Traits.(map[string]interface{})
		role, ok := traits["role"].(string)
		if !ok || role == "" {
			continue
		}

		// Check if user already has a role assigned in Casbin
		currentRoles, _ := u.enforcer.GetRolesForUser(id.Id)
		hasRole := false
		for _, r := range currentRoles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			log.Info("Syncing missing role for user", zap.String("user_id", id.Id), zap.String("role", role))
			_, err = u.enforcer.AddGroupingPolicy(id.Id, role)
			if err != nil {
				log.Error("failed to sync role for user", zap.String("user_id", id.Id), zap.Error(err))
				continue
			}
			syncedCount++
		}
	}

	log.Info("Casbin synchronization complete", zap.Int("total_identities", len(identities)), zap.Int("synced_new", syncedCount))
	return nil
}
