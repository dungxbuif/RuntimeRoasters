package tests

import (
	"context"
	"testing"

	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/delivery/grpc"
	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/demo-service/internal/usecase"
	demov1 "github.com/dungxbuif/RuntimeRoasters/runtime/demo/v1"
)

// mockDemoUsecase satisfies the usecase.DemoUsecase interface
type mockDemoUsecase struct {
	demo *domain.Demo
	err  error
}

func (m *mockDemoUsecase) GetDemo(ctx context.Context) (*domain.Demo, error) {
	return m.demo, m.err
}

func TestDemoUseCase(t *testing.T) {
	u := usecase.NewDemoUsecase()
	
	t.Run("returns static demo data", func(t *testing.T) {
		got, err := u.GetDemo(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ID != "demo-1" {
			t.Errorf("got id %s, want demo-1", got.ID)
		}
	})
}

func TestDemoHandler(t *testing.T) {
	t.Run("GetDemo returns protobuf response", func(t *testing.T) {
		mockUC := &mockDemoUsecase{
			demo: &domain.Demo{ID: "test-id", Message: "test-msg"},
		}
		handler := grpc.NewDemoHandler(mockUC)

		resp, err := handler.GetDemo(context.Background(), &demov1.GetDemoRequest{})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if resp.Id != "test-id" {
			t.Errorf("got id %s, want test-id", resp.Id)
		}
		if resp.Message != "test-msg" {
			t.Errorf("got message %s, want test-msg", resp.Message)
		}
	})
}
