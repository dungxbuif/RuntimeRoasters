package repository_tests

import (
	"context"
	"testing"
	"time"

	"RuntimeRoasters/apps/auth-service/internal/infrastructure/repository"
	"RuntimeRoasters/pkg/database"
	"github.com/stretchr/testify/assert"
)

func TestSecurityStateRepository_Integration(t *testing.T) {
	dsn := "postgres://dev:password@localhost:5432/auth_db?sslmode=disable"
	db, err := database.NewPostgres(database.PostgresConfig{URL: dsn})
	if err != nil {
		t.Skip("Skipping integration test: database not available")
		return
	}

	repo := repository.NewSecurityStateRepository(db)
	ctx := context.Background()
	subject := "test-subject-repo"

	// 1. Initial get should return epoch
	state, err := repo.GetBySubject(ctx, subject)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), state.Version)
	assert.True(t, state.TokensValidAfter.Unix() <= 0) // Should be epoch

	// 2. MonotonicUpdate
	state, err = repo.MonotonicUpdate(ctx, subject, "test_revoke", "admin")
	assert.NoError(t, err)
	assert.Equal(t, subject, state.Subject)
	assert.Equal(t, int64(1), state.Version)
	assert.Equal(t, "test_revoke", state.Reason)
	assert.True(t, state.TokensValidAfter.After(time.Now().Add(-1*time.Minute)))

	// 3. Second MonotonicUpdate should increment version
	time.Sleep(1 * time.Second)
	state2, err := repo.MonotonicUpdate(ctx, subject, "test_revoke_2", "admin")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), state2.Version)
	assert.True(t, state2.TokensValidAfter.After(state.TokensValidAfter) || state2.TokensValidAfter.Equal(state.TokensValidAfter))

	// 4. Get again
	state3, err := repo.GetBySubject(ctx, subject)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), state3.Version)
}
