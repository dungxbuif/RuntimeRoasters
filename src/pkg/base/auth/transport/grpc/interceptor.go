package authgrpc

import (
	"context"
	"strings"

	"RuntimeRoasters/pkg/base/auth/provider"
	"RuntimeRoasters/pkg/base/auth/token"
	"RuntimeRoasters/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type InterceptorOptions struct {
	PublicRoutes         []string
	TokenValidityChecker token.TokenValidityChecker
}

type InterceptorOption func(*InterceptorOptions)

func WithPublicRoutes(routes ...string) InterceptorOption {
	return func(o *InterceptorOptions) {
		o.PublicRoutes = append(o.PublicRoutes, routes...)
	}
}

func WithTokenValidityChecker(checker token.TokenValidityChecker) InterceptorOption {
	return func(o *InterceptorOptions) {
		o.TokenValidityChecker = checker
	}
}

// extractTokenFromMetadata retrieves the Bearer token from gRPC metadata.
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get(token.GRPCAuthMetadataKey)
	if len(authHeaders) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization metadata")
	}

	parts := strings.Split(authHeaders[0], " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != token.BearerPrefix {
		return "", status.Error(codes.Unauthenticated, "invalid authorization metadata format")
	}

	return parts[1], nil
}

// GRPCUnaryInterceptor creates a gRPC Unary Server Interceptor for JWT authentication.
func GRPCUnaryInterceptor(keyProvider provider.KeyProvider, expectedIssuer string, opts ...InterceptorOption) grpc.UnaryServerInterceptor {
	// 1. Cấu hình mặc định (Default configs)
	options := &InterceptorOptions{
		PublicRoutes: []string{
			"/grpc.health.v1.Health/Check", // Health Check luôn luôn public
		},
	}

	// 2. Ghi đè cấu hình dựa trên các Options được truyền vào
	for _, opt := range opts {
		opt(options)
	}

	// 3. Khởi tạo Map 1 lần duy nhất để check public routes cực nhanh (O(1)) và Thread-Safe
	publicRoutesMap := make(map[string]bool, len(options.PublicRoutes))
	for _, route := range options.PublicRoutes {
		publicRoutesMap[route] = true
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logger.FromContext(ctx)

		// Bước 1: Kiểm tra Whitelist
		if publicRoutesMap[info.FullMethod] {
			return handler(ctx, req)
		}

		// Bước 2: Extract Token
		jwtToken, err := extractTokenFromMetadata(ctx)
		if err != nil {
			log.Warn("missing or invalid metadata for auth", zap.Error(err))
			return nil, err
		}

		// Bước 3: Verify and Parse
		idToken, err := token.VerifyAndParseJWT(jwtToken, keyProvider, expectedIssuer)
		if err != nil {
			log.Warn("JWT verification failed", zap.Error(err))
			// Ngăn rò rỉ lỗi nội bộ ra ngoài
			return nil, status.Error(codes.Unauthenticated, "unauthorized or invalid token")
		}

		// Bước 3.5: Kiểm tra Revocation (Token Validity)
		if options.TokenValidityChecker != nil {
			valid, err := options.TokenValidityChecker.IsValid(ctx, idToken.Subject, idToken.IssuedAt)
			if err != nil {
				log.Error("token validity check failed", zap.Error(err), zap.String("subject", idToken.Subject))
				// Fail-Closed: Chặn nếu không kiểm tra được token validity
				return nil, status.Error(codes.Internal, "security check failed")
			}
			if !valid {
				log.Warn("token has been revoked", zap.String("subject", idToken.Subject))
				return nil, status.Error(codes.Unauthenticated, "token has been revoked")
			}
		}

		// Bước 4: Lưu Identity vào Context
		ctx = token.SetIdentityInContext(ctx, *idToken)

		// Bước 5: Gắn User ID vào OpenTelemetry Span
		token.RecordTracingData(ctx, idToken.Subject)

		return handler(ctx, req)
	}
}
