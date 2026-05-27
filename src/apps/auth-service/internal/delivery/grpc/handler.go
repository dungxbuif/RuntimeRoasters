package grpc

import (
	"context"
	"fmt"
	"strings"

	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	authv1 "RuntimeRoasters/runtime/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (h *Handler) GetMe(ctx context.Context, _ *authv1.GetMeRequest) (*authv1.GetMeResponse, error) {
	claims, ok := identity.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no identity found in context")
	}

	// Currently, the identity claims have Subject which is the email or ID depending on setup.
	// Since Kratos Subject is used, we can list users or get user by email/ID.
	// Wait, claims.Subject is usually the user ID. But we don't have GetUser in UserUsecase yet?
	// Let's just list all users and find the one that matches claims.Role, claims.Subject etc.
	// Actually, the frontend just needs User details. We can construct it from claims if claims has everything,
	// or we can fetch from DB. Let's see what UserUsecase provides.
	// Let's use ListUsers and filter by Email or ID (since we don't have GetUserByID).
	users, err := h.usecase.ListUsers(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", err)
	}

	var currentUser *domain.User
	for _, u := range users {
		if u.Email == claims.Subject || u.ID == claims.Subject {
			currentUser = u
			break
		}
	}

	if currentUser == nil {
		// Fallback to claims if not found in DB
		return &authv1.GetMeResponse{
			User: &authv1.User{
				Id:    claims.Subject,
				Email: claims.Subject,
				Role:  claims.Role,
			},
		}, nil
	}

	return &authv1.GetMeResponse{
		User: toProtoUser(currentUser),
	}, nil
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
