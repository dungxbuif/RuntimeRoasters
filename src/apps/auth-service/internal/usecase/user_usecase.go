package usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
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
	SeedUsers(ctx context.Context) error
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
	hydraCfg.Servers = hydra.ServerConfigurations{{URL: cfg.HydraAdminURL}}
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

	storeIDs := req.StoreIDs
	if storeIDs == nil {
		storeIDs = []string{}
	}
	warehouseIDs := req.WarehouseIDs
	if warehouseIDs == nil {
		warehouseIDs = []string{}
	}
	identityBody := *client.NewCreateIdentityBody(
		"default",
		map[string]interface{}{
			"email":         req.Email,
			"name":          req.Name,
			"role":          req.Role,
			"org_id":        req.OrgID,
			"store_ids":     storeIDs,
			"warehouse_ids": warehouseIDs,
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
	u.publishPolicyChanged(ctx, userId, "user_role_assigned")

	user := &domain.User{
		ID:           userId,
		Email:        req.Email,
		Name:         req.Name,
		Role:         req.Role,
		OrgID:        req.OrgID,
		StoreIDs:     storeIDs,
		WarehouseIDs: warehouseIDs,
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

	identities, _, err := u.kratosClient.IdentityAPI.ListIdentities(ctx).PerPage(250).Execute()
	if err != nil {
		log.Error("failed to list identities from kratos", zap.Error(err))
		return nil, fmt.Errorf("failed to list identities from kratos: %w", err)
	}

	log.Debug("Successfully fetched identities", zap.Int("count", len(identities)))

	users := make([]*domain.User, len(identities))
	for i, id := range identities {
		var traits map[string]interface{}
		if t, ok := id.Traits.(map[string]interface{}); ok {
			traits = t
		} else {
			traits = make(map[string]interface{})
		}

		email, _ := traits["email"].(string)
		name, _ := traits["name"].(string)
		orgID, _ := traits["org_id"].(string)
		storeIDs := traitStringSlice(traits["store_ids"])
		warehouseIDs := traitStringSlice(traits["warehouse_ids"])

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
			ID:           id.Id,
			Email:        email,
			Name:         name,
			Role:         role,
			OrgID:        orgID,
			StoreIDs:     storeIDs,
			WarehouseIDs: warehouseIDs,
		}
	}

	return users, nil
}

func (u *userUsecase) AcceptHydraLogin(ctx context.Context, req domain.AcceptLoginRequest) (*domain.AcceptLoginResponse, error) {
	log := logger.FromContext(ctx)
	log.Info("Accepting Hydra login request", zap.String("subject", req.Subject))

	var traits map[string]interface{}
	id, _, identityErr := u.kratosClient.IdentityAPI.GetIdentity(ctx, req.Subject).Execute()
	if identityErr == nil {
		if parsed, ok := id.Traits.(map[string]interface{}); ok {
			traits = parsed
		}
	}

	// 1. Fetch Role (Priority: Casbin -> Kratos Traits)
	roles, _ := u.enforcer.GetRolesForUser(req.Subject)
	role := ""
	if len(roles) > 0 {
		role = roles[0]
	} else if traits != nil {
		if r, ok := traits["role"].(string); ok {
			role = r
		}
	}

	if role == "" {
		role = "GUEST"
	}

	log.Info("Mapping role to Hydra session", zap.String("subject", req.Subject), zap.String("role", role))

	// 2. Accept Login with Session Claims
	accept := *hydra.NewAcceptOAuth2LoginRequest(req.Subject)
	accept.SetContext(map[string]interface{}{
		"email":         stringTrait(traits, "email"),
		"role":          role,
		"org_id":        stringTrait(traits, "org_id"),
		"store_ids":     traitStringSlice(traits["store_ids"]),
		"warehouse_ids": traitStringSlice(traits["warehouse_ids"]),
	})
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

func stringTrait(traits map[string]interface{}, key string) string {
	if traits == nil {
		return ""
	}
	value, _ := traits[key].(string)
	return value
}

func traitStringSlice(raw interface{}) []string {
	switch values := raw.(type) {
	case []string:
		return append([]string(nil), values...)
	case []interface{}:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if str, ok := value.(string); ok && str != "" {
				out = append(out, str)
			}
		}
		return out
	default:
		return nil
	}
}

func (u *userUsecase) SeedUsers(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("Starting Master User seeding in Kratos")

	// 1. Get existing users to avoid duplicates
	existingUsers, err := u.ListUsers(ctx)
	if err != nil {
		log.Warn("Could not fetch existing users for duplicate check, proceeding with caution", zap.Error(err))
	}
	existingMap := make(map[string]bool)
	for _, user := range existingUsers {
		existingMap[user.Email] = true
	}

	// 2. Define master users
	masterUsers := []domain.CreateUserRequest{
		// Farm Managers
		{Email: "mgr.farm.kho@runtimeroasters.com", Password: "Hello@123", Name: "K'Ho Farm Manager", Role: "FARM_MANAGER"},
		{Email: "mgr.farm.caudat@runtimeroasters.com", Password: "Hello@123", Name: "Cau Dat Manager", Role: "FARM_MANAGER"},
		{Email: "mgr.farm.sonpacamara@runtimeroasters.com", Password: "Hello@123", Name: "Son Pacamara Manager", Role: "FARM_MANAGER"},
		{Email: "mgr.farm.aeroco@runtimeroasters.com", Password: "Hello@123", Name: "Aeroco Manager", Role: "FARM_MANAGER"},
		{Email: "mgr.farm.trungnguyen@runtimeroasters.com", Password: "Hello@123", Name: "Trung Nguyen Manager", Role: "FARM_MANAGER"},
		{Email: "mgr.farm.chuse@runtimeroasters.com", Password: "Hello@123", Name: "Chu Se Manager", Role: "FARM_MANAGER"},

		// Store Managers
		{Email: "mgr.hn.hoankiem@runtimeroasters.com", Password: "Hello@123", Name: "Hoan Kiem Manager", Role: "STORE_MGR"},
		{Email: "mgr.hn.caugiay@runtimeroasters.com", Password: "Hello@123", Name: "Cau Giay Manager", Role: "STORE_MGR"},
		{Email: "mgr.hcm.d1@runtimeroasters.com", Password: "Hello@123", Name: "District 1 Manager", Role: "STORE_MGR"},
		{Email: "mgr.hcm.d7@runtimeroasters.com", Password: "Hello@123", Name: "District 7 Manager", Role: "STORE_MGR"},
		{Email: "mgr.dn.haichau@runtimeroasters.com", Password: "Hello@123", Name: "Hai Chau Manager", Role: "STORE_MGR"},

		// Warehouse Managers
		{Email: "mgr.wh.hn@runtimeroasters.com", Password: "Hello@123", Name: "HN Warehouse Manager", Role: "WAREHOUSE_MGR"},
		{Email: "mgr.wh.hcm@runtimeroasters.com", Password: "Hello@123", Name: "HCM Warehouse Manager", Role: "WAREHOUSE_MGR"},
		{Email: "mgr.wh.dn@runtimeroasters.com", Password: "Hello@123", Name: "DN Warehouse Manager", Role: "WAREHOUSE_MGR"},

		// Logistics Drivers (HN)
		{Email: "driver.hn.01@runtimeroasters.com", Password: "Hello@123", Name: "HN Driver 01", Role: "DRIVER"},
		{Email: "driver.hn.02@runtimeroasters.com", Password: "Hello@123", Name: "HN Driver 02", Role: "DRIVER"},
		{Email: "driver.hn.03@runtimeroasters.com", Password: "Hello@123", Name: "HN Driver 03", Role: "DRIVER"},
		{Email: "driver.hn.04@runtimeroasters.com", Password: "Hello@123", Name: "HN Driver 04", Role: "DRIVER"},
		// Logistics Drivers (HCM)
		{Email: "driver.hcm.01@runtimeroasters.com", Password: "Hello@123", Name: "HCM Driver 01", Role: "DRIVER"},
		{Email: "driver.hcm.02@runtimeroasters.com", Password: "Hello@123", Name: "HCM Driver 02", Role: "DRIVER"},
		{Email: "driver.hcm.03@runtimeroasters.com", Password: "Hello@123", Name: "HCM Driver 03", Role: "DRIVER"},
		{Email: "driver.hcm.04@runtimeroasters.com", Password: "Hello@123", Name: "HCM Driver 04", Role: "DRIVER"},
		// Logistics Drivers (DN)
		{Email: "driver.dn.01@runtimeroasters.com", Password: "Hello@123", Name: "DN Driver 01", Role: "DRIVER"},
		{Email: "driver.dn.02@runtimeroasters.com", Password: "Hello@123", Name: "DN Driver 02", Role: "DRIVER"},
		{Email: "driver.dn.03@runtimeroasters.com", Password: "Hello@123", Name: "DN Driver 03", Role: "DRIVER"},
	}

	var wg sync.WaitGroup
	for _, req := range masterUsers {
		if existingMap[req.Email] {
			log.Debug("User already exists, skipping", zap.String("email", req.Email))
			continue
		}

		wg.Add(1)
		go func(r domain.CreateUserRequest) {
			defer wg.Done()
			// Use a separate context with a longer timeout for each user creation
			// to avoid parent context cancellation affecting individuals too early
			createCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := u.CreateUser(createCtx, r)
			if err != nil {
				log.Error("Failed to create master user", zap.String("email", r.Email), zap.Error(err))
			} else {
				log.Info("Successfully seeded master user", zap.String("email", r.Email), zap.String("role", r.Role))
			}
		}(req)
	}
	wg.Wait()

	return nil
}

func (u *userUsecase) SyncCasbinWithKratos(ctx context.Context) error {
	log := logger.FromContext(ctx)
	log.Info("Starting Casbin synchronization with Kratos identities")

	identities, _, err := u.kratosClient.IdentityAPI.ListIdentities(ctx).PerPage(250).Execute()
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
			u.publishPolicyChanged(ctx, id.Id, "kratos_role_synced")
			syncedCount++
		}
	}

	log.Info("Casbin synchronization complete", zap.Int("total_identities", len(identities)), zap.Int("synced_new", syncedCount))
	return nil
}

func (u *userUsecase) publishPolicyChanged(ctx context.Context, key string, reason string) {
	if u.producer == nil {
		return
	}

	event := map[string]interface{}{
		"event_type":  "AUTH_POLICY_CHANGED",
		"reason":      reason,
		"occurred_at": time.Now().Format(time.RFC3339),
	}
	if err := u.producer.Publish(ctx, u.cfg.KafkaPolicyTopic, key, event); err != nil {
		logger.FromContext(ctx).Warn("failed to publish auth policy change event", zap.Error(err))
	}
}
