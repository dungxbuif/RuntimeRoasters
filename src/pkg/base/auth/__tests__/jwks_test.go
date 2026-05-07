package auth_tests

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/auth/provider"
)

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

func generateJWKS(kid string) jwksResponse {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	return jwksResponse{
		Keys: []jwk{
			{
				Kty: "RSA",
				Kid: kid,
				N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
				E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
			},
		},
	}
}

func TestJWKSCache_Resilience(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		json.NewEncoder(w).Encode(generateJWKS("test-kid"))
	}))
	defer server.Close()

	t.Run("NewJWKSCache retries on initial failure", func(t *testing.T) {
		cache, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err != nil {
			t.Fatalf("expected success after retries, got error: %v", err)
		}
		if atomic.LoadInt32(&callCount) < 3 {
			t.Errorf("expected at least 3 calls, got %d", callCount)
		}
		_, err = cache.GetPublicKey("test-kid")
		if err != nil {
			t.Errorf("failed to get key: %v", err)
		}
	})
}

func TestJWKSCache_PermanentFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	t.Run("fails after 5 retries", func(t *testing.T) {
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error after permanent failure, got nil")
		}
	})
}

func TestJWKSCache_SWR(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(generateJWKS("test-kid"))
	}))
	defer server.Close()

	t.Run("Stale-While-Revalidate logic", func(t *testing.T) {
		// TTL set to 100ms
		cache, _ := provider.NewJWKSCache(server.URL, "secret", 100*time.Millisecond)

		key1, _ := cache.GetPublicKey("test-kid")
		time.Sleep(150 * time.Millisecond) // wait for TTL expiry

		// SWR: Returns stale key1, triggers refresh
		key2, _ := cache.GetPublicKey("test-kid")
		if key1 != key2 {
			t.Error("SWR should return stale key immediately")
		}

		// Wait for refresh to complete
		var key3 *rsa.PublicKey
		var err error
		for i := 0; i < 10; i++ {
			time.Sleep(50 * time.Millisecond)
			key3, err = cache.GetPublicKey("test-kid")
			if err == nil && key3 != key1 {
				break
			}
		}

		if key3 == key1 {
			t.Error("expected new key after refresh")
		}
	})
}

func TestJWKSCache_Errors(t *testing.T) {
	t.Run("KID not found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(generateJWKS("test-kid"))
		}))
		defer server.Close()
		cache, _ := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		_, err := cache.GetPublicKey("non-existent")
		if err == nil {
			t.Error("expected error for missing KID, got nil")
		}
	})

	t.Run("malformed JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error for malformed JSON, got nil")
		}
	})

	t.Run("invalid base64 in JWK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := jwksResponse{
				Keys: []jwk{
					{Kty: "RSA", Kid: "bad-key", N: "!!!invalid-base64!!!", E: "AQAB"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error for invalid base64, got nil")
		}
	})

	t.Run("empty N or E in JWK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := jwksResponse{
				Keys: []jwk{
					{Kty: "RSA", Kid: "bad-key", N: "", E: ""},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error for empty fields, got nil")
		}
	})

	t.Run("non-RSA keys are skipped", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := jwksResponse{
				Keys: []jwk{
					{Kty: "EC", Kid: "ec-key"}, // Should be skipped
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error when no valid RSA keys found, got nil")
		}
	})

	t.Run("invalid E field in JWK", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp := jwksResponse{
				Keys: []jwk{
					{Kty: "RSA", Kid: "bad-key", N: "AQAB", E: "!!!invalid!!!"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()
		_, err := provider.NewJWKSCache(server.URL, "secret", 1*time.Minute)
		if err == nil {
			t.Error("expected error for invalid E field, got nil")
		}
	})
}
