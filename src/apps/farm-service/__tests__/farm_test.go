package tests

import (
	"context"
	"testing"

	"RuntimeRoasters/apps/farm-service/internal/delivery/grpc"
	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/apps/farm-service/internal/usecase"
	farmv1 "RuntimeRoasters/runtime/farm/v1"
)

// mockFarmUsecase satisfies the usecase.FarmUsecase interface
type mockFarmUsecase struct {
	farm *domain.Farm
	err  error
}

func (m *mockFarmUsecase) GetFarm(ctx context.Context) (*domain.Farm, error) {
	return m.farm, m.err
}

func TestFarmUseCase(t *testing.T) {
	u := usecase.NewFarmUsecase()
	
	t.Run("returns static farm data", func(t *testing.T) {
		got, err := u.GetFarm(context.Background())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.ID != "farm-1" {
			t.Errorf("got id %s, want farm-1", got.ID)
		}
	})
}

func TestFarmHandler(t *testing.T) {
	t.Run("GetFarm returns protobuf response", func(t *testing.T) {
		mockUC := &mockFarmUsecase{
			farm: &domain.Farm{ID: "test-id", Message: "test-msg"},
		}
		handler := grpc.NewFarmHandler(mockUC)

		resp, err := handler.GetFarm(context.Background(), &farmv1.GetFarmRequest{})
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
