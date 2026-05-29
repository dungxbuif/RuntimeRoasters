package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/apps/auth-service/internal/infrastructure/casbin"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/logger"
	authv1 "RuntimeRoasters/runtime/auth/v1"
	systemv1 "RuntimeRoasters/runtime/system/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	authv1.UnimplementedAuthServiceServer
	cfg      *config.Config
	enforcer *casbin.Enforcer
	usecase  usecase.UserUsecase
}

func NewHandler(cfg *config.Config, enforcer *casbin.Enforcer, u usecase.UserUsecase) *Handler {
	return &Handler{
		cfg:      cfg,
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

func (h *Handler) SeedData(ctx context.Context, req *systemv1.SeedDataRequest) (*systemv1.SeedDataResponse, error) {
	log := logger.FromContext(ctx)
	log.Info("Starting centralized system seeding from auth-service", zap.Bool("force", req.Force))

	// 1. Fetch users from Kratos via ListUsers to construct the users map
	usersMap := make(map[string]string)
	if h.usecase != nil {
		users, err := h.usecase.ListUsers(ctx)
		if err != nil {
			log.Error("failed to fetch users from Kratos for seeding", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
		}
		for _, u := range users {
			if u.Email != "" && u.ID != "" {
				usersMap[u.Email] = u.ID
			}
		}
		log.Info("Successfully compiled users map from Kratos", zap.Int("user_count", len(usersMap)))
	}

	// 2. Define downstream services to propagate the seed command from configuration
	var targets []string
	if h.cfg != nil {
		if h.cfg.FarmServiceURL != "" {
			targets = append(targets, h.cfg.FarmServiceURL+"/v1/system/seed")
		}
		if h.cfg.RetailServiceURL != "" {
			targets = append(targets, h.cfg.RetailServiceURL+"/v1/system/seed")
		}
		if h.cfg.LogisticsServiceURL != "" {
			targets = append(targets, h.cfg.LogisticsServiceURL+"/v1/system/seed")
		}
		if h.cfg.WarehouseServiceURL != "" {
			targets = append(targets, h.cfg.WarehouseServiceURL+"/v1/system/seed")
		}
	}

	// Fallback to legacy hardcoded logic if no targets defined in config
	if len(targets) == 0 {
		farmURL := "http://localhost:8083"
		retailURL := "http://localhost:8084"
		logisticsURL := "http://localhost:8085"
		warehouseURL := "http://localhost:8089"

		if h.cfg != nil && (strings.Contains(h.cfg.KratosAdminURL, "rr-kratos") || strings.Contains(h.cfg.KratosAdminURL, "kratos")) && !strings.Contains(h.cfg.KratosAdminURL, "localhost") {
			farmURL = "http://farm-service:8083"
			retailURL = "http://retail-service:8084"
			logisticsURL = "http://logistics-service:8085"
			warehouseURL = "http://warehouse-service:8089"
		}

		targets = []string{
			farmURL + "/v1/system/seed",
			retailURL + "/v1/system/seed",
			logisticsURL + "/v1/system/seed",
			warehouseURL + "/v1/system/seed",
		}
	}

	// 3. Construct JSON propagation payload
	type PropagatePayload struct {
		Force    bool              `json:"force"`
		UsersMap map[string]string `json:"users_map"`
	}
	payload := PropagatePayload{
		Force:    req.Force,
		UsersMap: usersMap,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal payload: %v", err)
	}

	// 4. Call each downstream service's /v1/system/seed in parallel
	httpClient := &http.Client{Timeout: 15 * time.Second}
	var wg sync.WaitGroup
	errChan := make(chan error, len(targets))

	for _, target := range targets {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			log.Info("Propagating seed request to downstream", zap.String("url", url))
			httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBytes))
			if err != nil {
				errChan <- fmt.Errorf("failed to create request for %s: %w", url, err)
				return
			}
			httpReq.Header.Set("Content-Type", "application/json")
			if h.cfg != nil && h.cfg.InternalSecret != "" {
				httpReq.Header.Set("X-Internal-Secret", h.cfg.InternalSecret)
			}

			resp, err := httpClient.Do(httpReq)
			if err != nil {
				errChan <- fmt.Errorf("failed to call %s: %w", url, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				errChan <- fmt.Errorf("calling %s returned status %d", url, resp.StatusCode)
				return
			}
			log.Info("Successfully seeded downstream service", zap.String("url", url))
		}(target)
	}

	wg.Wait()
	close(errChan)

	var errs []string
	for err := range errChan {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		errMsg := strings.Join(errs, "; ")
		log.Error("Centralized seeding completed with errors", zap.String("errors", errMsg))
		return nil, status.Errorf(codes.Internal, "propagation errors: %s", errMsg)
	}

	log.Info("Centralized seeding completed successfully across all services")
	return &systemv1.SeedDataResponse{
		Success:        true,
		Message:        "Centralized system seeding completed successfully",
		RecordsCreated: int64(len(usersMap)),
	}, nil
}

