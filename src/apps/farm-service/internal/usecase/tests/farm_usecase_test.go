package __tests__

import (
	"context"
	"testing"

	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/apps/farm-service/internal/usecase"
	"RuntimeRoasters/pkg/base/identity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the farm repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, farm *domain.Farm) error {
	args := m.Called(ctx, farm)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uint64) (*domain.Farm, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Farm), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context) ([]*domain.Farm, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Farm), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, farm *domain.Farm) error {
	args := m.Called(ctx, farm)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateFarm(t *testing.T) {
	repo := new(MockRepository)
	u := usecase.NewFarmUsecase(repo)

	ctx := context.Background()
	// Mock identity in context
	ctx = identity.InjectContext(ctx, identity.Claims{Subject: "owner-123", Role: "FARM_MANAGER"})

	farmReq := &domain.Farm{
		Name:       "Test Farm",
		Location:   "CAU_DAT",
		Area:       10.5,
		CoffeeType: "ARABICA",
	}

	repo.On("Create", ctx, mock.MatchedBy(func(f *domain.Farm) bool {
		return f.Name == "Test Farm" && f.OwnerID == "owner-123"
	})).Return(nil)

	created, err := u.CreateFarm(ctx, farmReq)

	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.Equal(t, "owner-123", created.OwnerID)
	repo.AssertExpectations(t)
}

func TestListFarms(t *testing.T) {
	repo := new(MockRepository)
	u := usecase.NewFarmUsecase(repo)

	ctx := context.Background()
	ctx = identity.InjectContext(ctx, identity.Claims{Subject: "owner-123", Role: "FARM_MANAGER"})

	expectedFarms := []*domain.Farm{
		{ID: 1, Name: "Farm 1", OwnerID: "owner-123"},
	}

	repo.On("List", ctx).Return(expectedFarms, nil)

	farms, err := u.ListFarms(ctx)

	assert.NoError(t, err)
	assert.Len(t, farms, 1)
	assert.Equal(t, "Farm 1", farms[0].Name)
	repo.AssertExpectations(t)
}
