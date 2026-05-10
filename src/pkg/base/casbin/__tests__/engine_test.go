package casbin_tests

import (
	"testing"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/stretchr/testify/assert"
)

func TestMapHTTPMethodToAction(t *testing.T) {
	tests := []struct {
		method   string
		expected string
	}{
		{"GET", "read"},
		{"POST", "write"},
		{"PUT", "write"},
		{"PATCH", "write"},
		{"DELETE", "delete"},
		{"HEAD", "read"}, // Default case
		{"OPTIONS", "read"}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			assert.Equal(t, tt.expected, casbin.MapHTTPMethodToAction(tt.method))
		})
	}
}
