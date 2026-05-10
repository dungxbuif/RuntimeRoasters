package auth_tests

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/token"
	"github.com/golang-jwt/jwt/v5"
)

type mockKeyProvider struct {
	keys map[string]*rsa.PublicKey
	err  error
}

func (m *mockKeyProvider) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	if m.err != nil {
		return nil, m.err
	}
	key, ok := m.keys[kid]
	if !ok {
		return nil, errors.New("key not found")
	}
	return key, nil
}

func generateRSAKeys(t *testing.T) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA keys: %v", err)
	}
	return privateKey, &privateKey.PublicKey
}

func createTestToken(t *testing.T, privateKey *rsa.PrivateKey, issuer, sub, role string) string {
	claims := jwt.MapClaims{
		"sub":    sub,
		"iss":    issuer,
		"role":   role,
		"exp":    time.Now().Add(time.Hour).Unix(),
		"jti":    "test-jti",
		"org_id": "test-org",
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header[token.JWTHeaderKeyID] = "test-key-id"
	ss, err := tok.SignedString(privateKey)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return ss
}
