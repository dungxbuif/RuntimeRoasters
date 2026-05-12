package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/usecase"
	authv1 "github.com/dungxbuif/RuntimeRoasters/runtime/auth/v1"
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
