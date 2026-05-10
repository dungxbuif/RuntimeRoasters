package casbinhttp_tests

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/casbin/casbin/v3"
	casbinhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin/transport/http"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Missing identity", func(t *testing.T) {
		engine := new(MockEngine)
		middleware := casbinhttp.GinMiddleware(engine)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		middleware(c)

		assert.Len(t, c.Errors, 1)
		assert.True(t, c.IsAborted())
	})

	t.Run("Allowed access", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "admin", "/test", "read").Return(true, nil)
		middleware := casbinhttp.GinMiddleware(engine)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request = c.Request.WithContext(identity.InjectContext(c.Request.Context(), identity.Claims{Role: "admin"}))

		middleware(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.False(t, c.IsAborted())
	})

	t.Run("Denied access", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "user", "/test", "read").Return(false, nil)
		middleware := casbinhttp.GinMiddleware(engine)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request = c.Request.WithContext(identity.InjectContext(c.Request.Context(), identity.Claims{Role: "user"}))

		middleware(c)

		assert.Len(t, c.Errors, 1)
		assert.True(t, c.IsAborted())
	})

	t.Run("Engine error", func(t *testing.T) {
		engine := new(MockEngine)
		engine.On("Enforce", "admin", "/test", "read").Return(false, errors.New("engine error"))
		middleware := casbinhttp.GinMiddleware(engine)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request = c.Request.WithContext(identity.InjectContext(c.Request.Context(), identity.Claims{Role: "admin"}))

		middleware(c)

		assert.True(t, c.IsAborted())
	})
}
