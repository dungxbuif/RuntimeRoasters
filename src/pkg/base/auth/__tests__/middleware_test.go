package auth_tests

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/token"
	authhttp "github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/transport/http"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	kid := "test-key"
	issuer := "http://test-issuer"

	keyProvider := &mockKeyProvider{
		keys: map[string]*rsa.PublicKey{
			kid: publicKey,
		},
	}

	setupRouter := func() (*httptest.ResponseRecorder, *gin.Engine) {
		w := httptest.NewRecorder()
		r := gin.New()
		r.Use(errs.GinErrorHandler())
		r.Use(authhttp.GinMiddleware(keyProvider, issuer))
		return w, r
	}

	createToken := func(claims jwt.MapClaims) string {
		t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		t.Header[token.JWTHeaderKeyID] = kid
		s, _ := t.SignedString(privateKey)
		return s
	}

	t.Run("valid token propagates identity", func(t *testing.T) {
		w, r := setupRouter()

		r.GET("/test", func(c *gin.Context) {
			id, ok := identity.FromContext(c.Request.Context())
			if !ok {
				t.Error("expected identity in context")
				return
			}
			if id.Subject != "user-1" {
				t.Errorf("got subject %s, want user-1", id.Subject)
			}
			c.Status(http.StatusOK)
		})

		claims := jwt.MapClaims{
			"sub":  "user-1",
			"iss":  issuer,
			"exp":  time.Now().Add(time.Hour).Unix(),
			"role": "user",
		}
		rawToken := createToken(claims)

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+rawToken)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
	})

	t.Run("missing authorization header", func(t *testing.T) {
		w, r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("invalid token format", func(t *testing.T) {
		w, r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "NotBearer some-token")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("invalid token returns unauthorized", func(t *testing.T) {
		w, r := setupRouter()
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})
		
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
