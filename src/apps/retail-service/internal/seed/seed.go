package seed

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"RuntimeRoasters/apps/retail-service/internal/domain"
)

//go:embed stores.json
var storesSeed []byte

func LoadStores() ([]domain.Store, error) {
	var stores []domain.Store
	if err := json.Unmarshal(storesSeed, &stores); err != nil {
		return nil, fmt.Errorf("load retail stores seed: %w", err)
	}
	if len(stores) == 0 {
		return nil, errors.New("retail stores seed is empty")
	}
	for _, store := range stores {
		if store.ID == "" || store.Name == "" || store.Status == "" {
			return nil, fmt.Errorf("invalid retail store seed: id=%q name=%q status=%q", store.ID, store.Name, store.Status)
		}
	}
	return stores, nil
}
