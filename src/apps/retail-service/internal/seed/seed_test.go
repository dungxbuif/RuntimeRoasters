package seed

import "testing"

func TestLoadStores(t *testing.T) {
	stores, err := LoadStores()
	if err != nil {
		t.Fatalf("LoadStores() error = %v", err)
	}
	if len(stores) != 5 {
		t.Fatalf("LoadStores() len = %d, want 5", len(stores))
	}
	if stores[0].ID == "" || stores[0].ManagerEmail == "" {
		t.Fatalf("LoadStores() returned incomplete first store: %+v", stores[0])
	}
}
