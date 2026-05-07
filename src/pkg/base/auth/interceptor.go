package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicRoutes = map[string]bool{
	"/grpc.health.v1.Health/Check": true,
}

// AddPublicRoute cho phép microservice đăng ký thêm các API công khai vào Whitelist.
func AddPublicRoute(fullMethod string) {
	publicRoutes[fullMethod] = true
}

// extractTokenFromMetadata retrieves the Bearer token from gRPC metadata.
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get(GRPCAuthMetadataKey)
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization metadata")
	}

	parts := strings.Split(authHeaders[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != BearerPrefix {
		return "", status.Error(codes.Unauthenticated, "invalid authorization metadata format")
	}

	return parts[1], nil
}

// GRPCUnaryInterceptor creates a gRPC Unary Server Interceptor for JWT authentication.
func GRPCUnaryInterceptor(keyProvider KeyProvider, expectedIssuer string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 1. Check Whitelist (Bypass Public APIs)
		if publicRoutes[info.FullMethod] {
			return handler(ctx, req)
		}

		jwtToken, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return nil, err
		}

		token, err := VerifyAndParseJWT(jwtToken, keyProvider, expectedIssuer)
		if err != nil {
			// Prevent Error Information Leakage
			// TODO: Log the detailed `err` internally here (e.g., logger.ErrorContext(ctx, "jwt validation failed", "error", err))
			
			// Return a generic error to the client
			return nil, status.Error(codes.Unauthenticated, "unauthorized or invalid token")
		}

		ctx = SetIdentityInContext(ctx, *token)

		RecordTracingData(ctx, token.Subject)

		return handler(ctx, req)
	}
}
