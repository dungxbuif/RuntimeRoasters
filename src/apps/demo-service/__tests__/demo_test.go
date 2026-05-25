package tests

import (
	"context"
	"testing"

	"RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/demo-service/internal/domain"
	demov1 "RuntimeRoasters/runtime/demo/v1"
	"github.com/stretchr/testify/assert"
)

// mockDemoUsecase implements DemoUsecase for testing
type mockDemoUsecase struct {
	demo *domain.Demo
	err  error
}

func (m *mockDemoUsecase) GetDemo(ctx context.Context) (*domain.Demo, error) {
	return m.demo, m.err
}

func TestDemoUseCase(t *testing.T) {
	// Simple test since logic is static for now
	t.Run("returns static demo data", func(t *testing.T) {
		// Since we can't easily mock the unexported struct without moving test
		// we'll just test the current implementation if it was exported or via New.
		// For now, this is a placeholder to establish the __tests__ pattern.
	})
}

func TestDemoHandler(t *testing.T) {
	t.Run("GetDemo returns protobuf response", func(t *testing.T) {
		mockUC := &mockDemoUsecase{
			demo: &domain.Demo{ID: "test-id", Message: "test-msg"},
		}
		handler := grpc.NewDemoHandler(mockUC)

		resp, err := handler.GetDemo(context.Background(), &demov1.GetDemoRequest{})

		assert.NoError(t, err)
		assert.Equal(t, "test-id", resp.Id)
		assert.Equal(t, "test-msg", resp.Message)
	})
}
