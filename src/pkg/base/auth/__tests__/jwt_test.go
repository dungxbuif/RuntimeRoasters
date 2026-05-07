package auth_tests

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/token"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/golang-jwt/jwt/v5"
)

func TestVerifyAndParseJWT(t *testing.T) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	kid := "test-key-1"
	issuer := "http://test-issuer"

	keyProvider := &mockKeyProvider{
		keys: map[string]*rsa.PublicKey{
			kid: publicKey,
		},
	}

	createToken := func(claims jwt.MapClaims, signingKey *rsa.PrivateKey, kidHeader string) string {
		t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		if kidHeader != "" {
			t.Header[token.JWTHeaderKeyID] = kidHeader
		}
		s, _ := t.SignedString(signingKey)
		return s
	}

	t.Run("successful verification", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub":      "user-123",
			"iss":      issuer,
			"exp":      time.Now().Add(time.Hour).Unix(),
			"role":     "admin",
			"org_id":   "org-1",
			"jti":      "jti-1",
		}
		rawToken := createToken(claims, privateKey, kid)

		got, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		want := &identity.Claims{
			Subject: "user-123",
			Role:    "admin",
			OrgID:   "org-1",
			JTI:     "jti-1",
		}

		if *got != *want {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": issuer,
			"exp": time.Now().Add(-time.Hour).Unix(),
		}
		rawToken := createToken(claims, privateKey, kid)

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for expired token, got nil")
		}
	})

	t.Run("wrong issuer", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": "http://wrong-issuer",
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		rawToken := createToken(claims, privateKey, kid)

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for wrong issuer, got nil")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		otherPrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		// Sign with a different key
		rawToken := createToken(claims, otherPrivateKey, kid)

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for invalid signature, got nil")
		}
	})

	t.Run("missing kid", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		rawToken := createToken(claims, privateKey, "") // empty kid

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for missing kid, got nil")
		}
	})

	t.Run("unexpected signing method", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": issuer,
		}
		// Create HS256 token instead of RS256
		tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		rawToken, _ := tok.SignedString([]byte("some-secret"))

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for HS256 token, got nil")
		}
	})

	t.Run("key provider error", func(t *testing.T) {
		claims := jwt.MapClaims{
			"sub": "user-123",
			"iss": issuer,
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		rawToken := createToken(claims, privateKey, "unknown-kid")

		_, err := token.VerifyAndParseJWT(rawToken, keyProvider, issuer)
		if err == nil {
			t.Error("expected error for unknown kid, got nil")
		}
	})
}
