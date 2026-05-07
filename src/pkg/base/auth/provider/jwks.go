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
)

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKSCache struct {
	jwksURL        string
	internalSecret string
	cacheTTL       time.Duration

	mu         sync.RWMutex
	keys       map[string]*rsa.PublicKey
	fetchedAt  time.Time
	refreshing bool
}

// NewJWKSCache initializes a new JWKS cache and performs an initial fetch.
// Nó thực hiện cơ chế Retry với Exponential Backoff để đảm bảo service có thể đợi
// hạ tầng Identity sẵn sàng mà không bị sập ngay lập tức (Panic).
func NewJWKSCache(jwksURL, internalSecret string, ttl time.Duration) (*JWKSCache, error) {
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
	maxRetries := 5
	backoff := 2 * time.Second

	var err error
	for i := 0; i < maxRetries; i++ {
		if err = cache.refresh(); err == nil {
			return cache, nil
		}
		
		fmt.Printf("Attempt %d: Failed to fetch JWKS from %s, retrying in %v... (error: %v)\n", i+1, jwksURL, backoff, err)
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	return nil, fmt.Errorf("initial jwks fetch failed after %d retries: %w", maxRetries, err)
}

func (c *JWKSCache) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	c.mu.RLock()
	key, ok := c.keys[kid]
	fetchedAt := c.fetchedAt
	refreshing := c.refreshing
	c.mu.RUnlock()

	if ok && time.Since(fetchedAt) < c.cacheTTL {
		return key, nil
	}

	if !refreshing {
		go func() {
			_ = c.refresh()
		}()
	}

	if ok {
		return key, nil
	}

	return nil, fmt.Errorf("key id %s not found in jwks cache", kid)
}

func (c *JWKSCache) refresh() error {
	c.mu.Lock()
	if c.refreshing && time.Since(c.fetchedAt) < 10*time.Second {
		c.mu.Unlock()
		return nil
	}
	c.refreshing = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.refreshing = false
		c.mu.Unlock()
	}()

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, c.jwksURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-Internal-Secret", c.internalSecret)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks endpoint returned status %d", resp.StatusCode)
	}

	var jwks jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
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
	c.fetchedAt = time.Now()
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