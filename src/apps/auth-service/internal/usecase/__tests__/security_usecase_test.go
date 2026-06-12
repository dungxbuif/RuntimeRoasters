package usecase_tests

import (
	"context"
	"testing"
	"time"

	"RuntimeRoasters/apps/auth-service/config"
	"RuntimeRoasters/apps/auth-service/internal/domain"
	"RuntimeRoasters/apps/auth-service/internal/usecase"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSecurityStateRepository struct {
	mock.Mock
}

func (m *mockSecurityStateRepository) GetBySubject(ctx context.Context, subject string) (*domain.UserSecurityState, error) {
	args := m.Called(ctx, subject)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserSecurityState), args.Error(1)
}

func (m *mockSecurityStateRepository) MonotonicUpdate(ctx context.Context, subject string, reason string, updatedBy string) (*domain.UserSecurityState, error) {
	args := m.Called(ctx, subject, reason, updatedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserSecurityState), args.Error(1)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) WithTx(ctx context.Context, fn func(txCtx context.Context) error) error {
	return fn(ctx)
}

func TestSecurityUsecase_GetTokensValidAfter(t *testing.T) {
	s := miniredis.RunT(t)
	rClient := redis.NewClient(&redis.Options{Addr: s.Addr()})
	
	repo := new(mockSecurityStateRepository)
	txm := new(mockTxManager)
	cfg := &config.Config{}
	
	u := usecase.NewSecurityUsecase(cfg, repo, txm, rClient, nil)
	ctx := context.Background()
	subject := "user123"

	// 1. Cache Miss - fallback to repo
	validAfterTime := time.Now().Truncate(time.Second)
	repo.On("GetBySubject", ctx, subject).Return(&domain.UserSecurityState{
		Subject:          subject,
		TokensValidAfter: validAfterTime,
		Version:          1,
	}, nil).Once()

	res, err := u.GetTokensValidAfter(ctx, subject)
	assert.NoError(t, err)
	assert.Equal(t, validAfterTime.Unix(), res.Unix())
	repo.AssertExpectations(t)

	// 2. Cache Hit - should not call repo
	// Wait a bit to ensure async cache write completes if it was async, but here it's sync
	res2, err := u.GetTokensValidAfter(ctx, subject)
	assert.NoError(t, err)
	assert.Equal(t, validAfterTime.Unix(), res2.Unix())
	repo.AssertNumberOfCalls(t, "GetBySubject", 1) // still 1 call
}
