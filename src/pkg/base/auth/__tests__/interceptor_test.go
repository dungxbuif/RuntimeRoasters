package auth_tests

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"RuntimeRoasters/pkg/base/auth/token"
	"RuntimeRoasters/pkg/base/auth/transport/grpc"
	"RuntimeRoasters/pkg/base/identity"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestGRPCUnaryInterceptor(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	kid := "test-key"
	issuer := "http://test-issuer"

	keyProvider := &mockKeyProvider{
		keys: map[string]*rsa.PublicKey{
			kid: publicKey,
		},
	}

	createToken := func(claims jwt.MapClaims) string {
		t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		t.Header[token.JWTHeaderKeyID] = kid
		s, _ := t.SignedString(privateKey)
		return s
	}

	interceptor := authgrpc.GRPCUnaryInterceptor(keyProvider, issuer)

	t.Run("valid token propagates identity", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-grpc",
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		rawToken := createToken(claims)

		md := metadata.Pairs(token.GRPCAuthMetadataKey, "Bearer "+rawToken)
		ctx := metadata.NewIncomingContext(context.Background(), md)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			id, ok := identity.FromContext(ctx)
			if !ok {
				return nil, errors.New("no identity in context")
			}
			if id.Subject != "user-grpc" {
				return nil, errors.New("wrong subject")
			}
			return "ok", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/service.v1.Service/Method"}
		resp, err := interceptor(ctx, nil, info, handler)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp != "ok" {
			t.Errorf("got response %v, want ok", resp)
		}
	})

	t.Run("public route skips auth", func(t *testing.T) {
		ctx := context.Background() // No metadata
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		}

		// Health check is public by default
		info := &grpc.UnaryServerInfo{FullMethod: "/grpc.health.v1.Health/Check"}
		resp, err := interceptor(ctx, nil, info, handler)

		if err != nil {
			t.Fatalf("expected no error for public route, got %v", err)
		}
		if resp != "ok" {
			t.Errorf("got response %v, want ok", resp)
		}
	})

	t.Run("custom public route skips auth", func(t *testing.T) {
		// Create interceptor with custom public route
		interceptorCustom := authgrpc.GRPCUnaryInterceptor(keyProvider, issuer, authgrpc.WithPublicRoutes("/custom.v1.Service/PublicMethod"))
		
		ctx := context.Background()
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "custom-ok", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/custom.v1.Service/PublicMethod"}
		resp, err := interceptorCustom(ctx, nil, info, handler)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp != "custom-ok" {
			t.Errorf("got response %v, want custom-ok", resp)
		}
	})

	t.Run("malformed token returns unauthenticated", func(t *testing.T) {
		md := metadata.Pairs(token.GRPCAuthMetadataKey, "InvalidFormat token")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/service.v1.Service/Method"}
		_, err := interceptor(ctx, nil, info, handler)

		if status.Code(err) != codes.Unauthenticated {
			t.Errorf("expected Unauthenticated code, got %v", status.Code(err))
		}
	})

	t.Run("missing metadata returns unauthenticated", func(t *testing.T) {
		ctx := context.Background()
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/service.v1.Service/Method"}
		_, err := interceptor(ctx, nil, info, handler)

		if status.Code(err) != codes.Unauthenticated {
			t.Errorf("expected Unauthenticated code, got %v", status.Code(err))
		}
	})

	t.Run("metadata exists but missing authorization key", func(t *testing.T) {
		md := metadata.Pairs("some-other-key", "some-value")
		ctx := metadata.NewIncomingContext(context.Background(), md)

		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return "ok", nil
		}

		info := &grpc.UnaryServerInfo{FullMethod: "/service.v1.Service/Method"}
		_, err := interceptor(ctx, nil, info, handler)

		if status.Code(err) != codes.Unauthenticated {
			t.Errorf("expected Unauthenticated code, got %v", status.Code(err))
		}
	})
}
