/*
*Package provider cung cấp các cơ chế lấy khóa công khai (Public Key) để xác thực JWT.
* File jwks.go triển khai cơ chế lấy khóa từ endpoint JWKS của Identity Server (Hydra).
* Cơ chế hoạt động:
* 1. Fetcher: Truy cập endpoint JWKS qua Nginx Gatekeeper với mã bí mật nội bộ (Internal Secret).
* 2. Caching (Stale-While-Revalidate):
*   - Nếu khóa trong bộ nhớ đệm (cache) còn hiệu lực, trả về ngay lập tức.
*   - Nếu khóa đã hết hạn (TTL), hệ thống vẫn trả về khóa cũ (stale) cho request hiện tại
*     để đảm bảo không bị block (low latency), đồng thời khởi chạy một goroutine chạy ngầm
*     để cập nhật khóa mới (revalidate).
*   - Điều này giúp hệ thống luôn sẵn sàng ngay cả khi Identity Server gặp sự cố tạm thời.
*
* Luồng gọi hàm (Call Flow):
* GetPublicKey(kid)
*
*	-> Check Cache
*	-> [Nếu hết hạn] refresh() (goroutine)
*	   -> http.Get(JWKS_URL)
*	   -> decodeRSA(n, e)
*	   -> Update Cache
 */
package provider

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
)

type JWKSCache struct {
	jwksURL        string
	internalSecret string
	cacheTTL       time.Duration
	keys           map[string]*rsa.PublicKey
	mu             sync.RWMutex
	lastFetch      time.Time
	sf             singleflight.Group
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

// NewJWKSCache initializes a new JWKS cache and performs an initial fetch.
func NewJWKSCache(jwksURL, internalSecret string, ttl time.Duration) (*JWKSCache, error) {
	if jwksURL == "" {
		return nil, errors.New("JWKS URL is required")
	}
	if ttl == 0 {
		ttl = 5 * time.Minute
	}

	cache := &JWKSCache{
		jwksURL:        jwksURL,
		internalSecret: internalSecret,
		cacheTTL:       ttl,
		keys:           make(map[string]*rsa.PublicKey),
	}

	// Cơ chế Retry khi khởi động (Bootstrapping Resilience)
	log := logger.GetLogger()
	var err error
	for i := 0; i < 5; i++ {
		if err = cache.refresh(); err == nil {
			break
		}
		log.Warn("Failed to fetch JWKS", zap.String("url", jwksURL), zap.Int("attempt", i+1), zap.Error(err))
		time.Sleep(time.Duration(1<<i) * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS after retries: %w", err)
	}

	return cache, nil
}

func (c *JWKSCache) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	key, ok := c.keys[kid]
	lastFetch := c.lastFetch
	c.mu.RUnlock()

	// Case 1: Key not found - Could be a Key Rotation.
	// We MUST perform a synchronous refresh to check for new keys.
	if !ok {
		_, err, _ := c.sf.Do("refresh", func() (interface{}, error) {
			return nil, c.refresh()
		})
		if err != nil {
			return nil, fmt.Errorf("failed to refresh JWKS for missing kid %s: %w", kid, err)
		}

		// Re-check cache after sync refresh
		c.mu.RLock()
		defer c.mu.RUnlock()
		key, ok = c.keys[kid]
		if !ok {
			return nil, fmt.Errorf("kid %s not found in JWKS even after refresh", kid)
		}
		return key, nil
	}

	// Case 2: Key found but cache is stale (SWR).
	// Trigger background refresh using singleflight to prevent Thundering Herd.
	if time.Since(lastFetch) > c.cacheTTL {
		go func() {
			c.sf.Do("refresh", func() (interface{}, error) {
				return nil, c.refresh()
			})
		}()
	}

	return key, nil
}

func (c *JWKSCache) refresh() error {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.internalSecret != "" {
		req.Header.Set("X-Internal-Secret", c.internalSecret)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks endpoint returned status %d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS response: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Kty == "RSA" {
			pubKey, err := decodeRSA(key.N, key.E)
			if err != nil {
				return fmt.Errorf("failed to decode RSA key %s: %w", key.Kid, err)
			}
			newKeys[key.Kid] = pubKey
		}
	}

	if len(newKeys) == 0 {
		return errors.New("no valid RSA keys found in JWKS")
	}

	c.mu.Lock()
	c.keys = newKeys
	c.lastFetch = time.Now()
	c.mu.Unlock()

	return nil
}

func decodeRSA(nStr, eStr string) (*rsa.PublicKey, error) {
	if nStr == "" || eStr == "" {
		return nil, errors.New("empty N or E in JWK")
	}
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	var e int
	for _, b := range eBytes {
		e = (e << 8) | int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}
