package e2e

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

// SecurityTestSuite verifies the Two-Gate Security Model (KrakenD Scope + Service Casbin)
type SecurityTestSuite struct {
	suite.Suite
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	jwksServer *http.Server
	issuer     string
	gatewayURL string
}

func (s *SecurityTestSuite) SetupSuite() {
	s.issuer = "http://localhost:4444/"
	s.gatewayURL = "http://localhost:8081"

	// 1. Generate RSA Key Pair
	var err error
	s.privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	s.NoError(err)
	s.publicKey = &s.privateKey.PublicKey

	// 2. Start Mock JWKS Server on a fixed port that KrakenD can reach
	// Note: In deployments/krakend/krakend.json, jwk_url should point to this if testing locally.
	// For this test to work with the real KrakenD container, we'd need to bridge network.
	// For now, we'll implement the logic and provide the verification steps.
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/jwks.json", s.handleJWKS)
	
	s.jwksServer = &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	go func() {
		if err := s.jwksServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Mock JWKS server failed: %v\n", err)
		}
	}()
}

func (s *SecurityTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.jwksServer.Shutdown(ctx)
}

func (s *SecurityTestSuite) handleJWKS(w http.ResponseWriter, r *http.Request) {
	// Standard JWKS format
	n := base64.RawURLEncoding.EncodeToString(s.publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.publicKey.E)).Bytes())

	jwks := map[string]interface{}{
		"keys": []map[string]interface{}{
			{
				"kty": "RSA",
				"use": "sig",
				"kid": "test-key-id",
				"alg": "RS256",
				"n":   n,
				"e":   e,
			},
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}

func (s *SecurityTestSuite) createToken(role string, scopes []string) string {
	claims := jwt.MapClaims{
		"sub":   "user-123",
		"iss":   s.issuer,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
		"jti":   "test-jti-456",
		"role":  role,
		"scope": scopes,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "test-key-id"
	ss, err := token.SignedString(s.privateKey)
	s.NoError(err)
	return ss
}

// --- Test Cases ---

func (s *SecurityTestSuite) Test_S1_Unauthorized_NoToken() {
	resp, err := http.Get(s.gatewayURL + "/v1/farm/ping")
	s.NoError(err)
	defer resp.Body.Close()

	// Gate 1 (KrakenD) should block this
	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *SecurityTestSuite) Test_S2_Forbidden_Gate1_InvalidScope() {
	// Token has correct role but WRONG scope for KrakenD
	token := s.createToken("admin", []string{"wrong:scope"})
	
	req, _ := http.NewRequest("GET", s.gatewayURL+"/v1/farm/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	// Gate 1 (KrakenD) should block this because it expects 'farm:read'
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *SecurityTestSuite) Test_S3_Forbidden_Gate2_InvalidRole() {
	// Token has correct scope (Passes Gate 1) but WRONG role (Blocked by Casbin at Gate 2)
	token := s.createToken("guest", []string{"farm:read"})
	
	req, _ := http.NewRequest("GET", s.gatewayURL+"/v1/farm/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	// Gate 2 (Service/Casbin) should block this. 
	// Note: Middleware usually returns 403 for RBAC failures.
	s.Equal(http.StatusForbidden, resp.StatusCode)
}

func (s *SecurityTestSuite) Test_S4_Success_FullAccess() {
	// Token has correct scope AND correct role
	token := s.createToken("admin", []string{"farm:read"})
	
	req, _ := http.NewRequest("GET", s.gatewayURL+"/v1/farm/ping", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	
	resp, err := http.DefaultClient.Do(req)
	s.NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)
}

func TestSecuritySuite(t *testing.T) {
	// Skip by default if KrakenD is not reachable to avoid CI failure
	conn, err := net.DialTimeout("tcp", "localhost:8081", 1*time.Second)
	if err != nil {
		t.Skip("KrakenD not running at localhost:8081, skipping E2E security tests")
		return
	}
	conn.Close()
	
	suite.Run(t, new(SecurityTestSuite))
}
