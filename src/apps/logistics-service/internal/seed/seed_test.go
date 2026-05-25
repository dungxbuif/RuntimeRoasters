package seed

import "testing"

func TestLoad(t *testing.T) {
	data, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(data.Vehicles) != 3 {
		t.Fatalf("Load() vehicles len = %d, want 3", len(data.Vehicles))
	}
	if len(data.Drivers) != 5 {
		t.Fatalf("Load() drivers len = %d, want 5", len(data.Drivers))
	}
	if len(data.Locations) != 4 {
		t.Fatalf("Load() locations len = %d, want 4", len(data.Locations))
	}
}
