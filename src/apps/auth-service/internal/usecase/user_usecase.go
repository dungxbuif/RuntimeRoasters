package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/config"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	"github.com/ory/client-go"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
}

type userUsecase struct {
	cfg          *config.Config
	enforcer     *casbin.Enforcer
	kratosClient *client.APIClient
	producer     kafka.Producer
}

func NewUserUsecase(cfg *config.Config, enforcer *casbin.Enforcer, producer kafka.Producer) UserUsecase {
	kratosCfg := client.NewConfiguration()
	kratosCfg.Servers = client.ServerConfigurations{
		{
			URL: cfg.KratosAdminURL,
		},
	}
	kratosClient := client.NewAPIClient(kratosCfg)

	return &userUsecase{
		cfg:          cfg,
		enforcer:     enforcer,
		kratosClient: kratosClient,
		producer:     producer,
	}
}

func (u *userUsecase) CreateUser(ctx context.Context, req domain.CreateUserRequest) (*domain.User, error) {
	identityBody := *client.NewCreateIdentityBody(
		"default",
		map[string]interface{}{
			"email": req.Email,
			"name":  req.Name,
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
		return nil, fmt.Errorf("failed to create identity in kratos: %w", err)
	}

	userId := createdIdentity.Id

	_, err = u.enforcer.AddGroupingPolicy(userId, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to assign role in casbin: %w", err)
	}

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
		fmt.Printf("Warning: failed to publish user.created event: %v\n", err)
	}

	return user, nil
}

func (u *userUsecase) ListUsers(ctx context.Context) ([]*domain.User, error) {
	identities, _, err := u.kratosClient.IdentityAPI.ListIdentities(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list identities from kratos: %w", err)
	}

	users := make([]*domain.User, len(identities))
	for i, id := range identities {
		traits := id.Traits.(map[string]interface{})
		email, _ := traits["email"].(string)
		name, _ := traits["name"].(string)

		roles, _ := u.enforcer.GetRolesForUser(id.Id)
		role := "unknown"
		if len(roles) > 0 {
			role = roles[0]
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
