package auth_tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/http"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"crypto/rsa"
)

func TestGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	expectedIssuer := "http://localhost:4444/"

	t.Run("valid token propagates identity", func(t *testing.T) {
		privateKey, publicKey := generateRSAKeys(t)
		keyProvider := &mockKeyProvider{
			keys: map[string]*rsa.PublicKey{
				"test-key-id": publicKey,
			},
		}
		token := createTestToken(t, privateKey, expectedIssuer, "user-123", "admin")

		r := gin.New()
		r.Use(errs.GinErrorHandler())
		r.Use(authhttp.GinMiddleware(keyProvider, expectedIssuer))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		keyProvider := &mockKeyProvider{}
		r := gin.New()
		r.Use(errs.GinErrorHandler())
		r.Use(authhttp.GinMiddleware(keyProvider, expectedIssuer))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token format", func(t *testing.T) {
		keyProvider := &mockKeyProvider{}
		r := gin.New()
		r.Use(errs.GinErrorHandler())
		r.Use(authhttp.GinMiddleware(keyProvider, expectedIssuer))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token returns unauthorized", func(t *testing.T) {
		keyProvider := &mockKeyProvider{}
		r := gin.New()
		r.Use(errs.GinErrorHandler())
		r.Use(authhttp.GinMiddleware(keyProvider, expectedIssuer))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
