package casbingrpc_tests

import (
	"context"
	"errors"
	"testing"

	"github.com/casbin/casbin/v3"
	casbingrpc "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/grpc"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	)

	type MockEngine struct {
	mock.Mock
	}

	func (m *MockEngine) Enforce(rvals ...interface{}) (bool, error) {
	args := m.Called(rvals...)
	return args.Bool(0), args.Error(1)
	}

	func (m *MockEngine) GetEnforcer() *casbin.Enforcer {
	return nil
	}

	func (m *MockEngine) Sync(ctx context.Context) error {
	return nil
	}

	func TestGRPCUnaryInterceptor(t *testing.T) {
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	}

	t.Run("Missing identity", func(t *testing.T) {
		engine := new(MockEngine)
		interceptor := casbingrpc.GRPCUnaryInterceptor(engine)


		_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
	})

	t.Run("Allowed access", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "admin", "/test", "write").Return(true, nil)
		interceptor := casbingrpc.GRPCUnaryInterceptor(engine)

		ctx := identity.InjectContext(context.Background(), identity.Claims{Role: "admin"})
		resp, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)

		assert.NoError(t, err)
		assert.Equal(t, "ok", resp)
	})

	t.Run("Denied access", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "user", "/test", "write").Return(false, nil)
		interceptor := casbingrpc.GRPCUnaryInterceptor(engine)

		ctx := identity.InjectContext(context.Background(), identity.Claims{Role: "user"})
		_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("Engine error", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "admin", "/test", "write").Return(false, errors.New("engine error"))
		interceptor := casbingrpc.GRPCUnaryInterceptor(engine)

		ctx := identity.InjectContext(context.Background(), identity.Claims{Role: "admin"})
		_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test"}, handler)

		assert.Error(t, err)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
	})
}
