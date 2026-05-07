package auth_tests

import (
	"crypto/rsa"
	"errors"
)

type mockKeyProvider struct {
	keys map[string]*rsa.PublicKey
}

func (m *mockKeyProvider) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	key, ok := m.keys[kid]
	if !ok {
		return nil, errors.New("key not found")
	}
	return key, nil
}
