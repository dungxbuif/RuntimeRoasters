package grpc

import (
	"context"
	"fmt"
	"strings"

	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	authv1 "RuntimeRoasters/runtime/auth/v1"
)

type Handler struct {
	authv1.UnimplementedAuthServiceServer
	enforcer *casbin.Enforcer
	usecase  usecase.UserUsecase
}

func NewHandler(enforcer *casbin.Enforcer, u usecase.UserUsecase) *Handler {
	return &Handler{
		enforcer: enforcer,
		usecase:  u,
	}
}

func (h *Handler) GetFullSnapshot(ctx context.Context, req *authv1.GetFullSnapshotRequest) (*authv1.GetFullSnapshotResponse, error) {
	policies, _ := h.enforcer.GetPolicy()
	groupingPolicies, _ := h.enforcer.GetGroupingPolicy()

	var result []string

	// Format p policies: p, role, obj, act
	for _, p := range policies {
		result = append(result, fmt.Sprintf("p, %s", strings.Join(p, ", ")))
	}

	// Format g policies: g, user, role
	for _, g := range groupingPolicies {
		result = append(result, fmt.Sprintf("g, %s", strings.Join(g, ", ")))
	}

	return &authv1.GetFullSnapshotResponse{
		Policies: result,
	}, nil
}

func (h *Handler) AcceptLogin(ctx context.Context, req *authv1.AcceptLoginRequest) (*authv1.AcceptLoginResponse, error) {
	res, err := h.usecase.AcceptHydraLogin(ctx, domain.AcceptLoginRequest{
		LoginChallenge: req.LoginChallenge,
		Subject:        req.Subject,
		Remember:       req.Remember,
		RememberFor:    req.RememberFor,
	})
	if err != nil {
		return nil, err
	}

	return &authv1.AcceptLoginResponse{
		RedirectTo: res.RedirectTo,
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	user, err := h.usecase.CreateUser(ctx, domain.CreateUserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Name:     req.GetName(),
		Role:     req.GetRole(),
	})
	if err != nil {
		return nil, err
	}

	return &authv1.CreateUserResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *Handler) ListUsers(ctx context.Context, _ *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	users, err := h.usecase.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	res := &authv1.ListUsersResponse{
		Users: make([]*authv1.User, 0, len(users)),
	}
	for _, user := range users {
		res.Users = append(res.Users, toProtoUser(user))
	}

	return res, nil
}

func toProtoUser(user *domain.User) *authv1.User {
	if user == nil {
		return nil
	}

	return &authv1.User{
		Id:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}
