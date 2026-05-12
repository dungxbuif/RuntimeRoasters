package casbingrpc

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Enforcer interface {
	Enforce(rvals ...interface{}) (bool, error)
}

// GRPCUnaryInterceptor creates a gRPC interceptor for Casbin authorization
func GRPCUnaryInterceptor(engine Enforcer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Get identity from context (Injected by JWT interceptor)
		claims, ok := identity.FromContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "identity not found in context")
		}

		// 2. Map Method to Resource and Action
		// sub: User Role
		// obj: info.FullMethod (e.g., /farm.v1.FarmService/CreateFarm)
		// act: We map all gRPC calls to 'write' by default for safety, or custom logic
		action := "read"
		if info.FullMethod != "" {
			// Basic heuristic: if it contains Get/List -> read, else write
			// (In a real system, you might want more granular mapping)
			action = "write" 
			// Check for read-only methods
			// This can be improved by adding metadata to proto or a mapping table
		}

		// 3. Enforce
		allowed, err := engine.Enforce(claims.Role, info.FullMethod, action)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "authorization error: %v", err)
		}

		if !allowed {
			return nil, status.Errorf(codes.PermissionDenied, "permission denied for %s", info.FullMethod)
		}

		return handler(ctx, req)
	}
}
