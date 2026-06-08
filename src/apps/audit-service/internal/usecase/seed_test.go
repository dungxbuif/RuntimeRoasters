package usecase

import (
	"context"
	"testing"

	"RuntimeRoasters/apps/audit-service/internal/domain"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeedDemoDataCreatesIdempotentAuditHashChain(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&domain.AuditLog{}))
	service := NewService(db, nil)

	created, err := service.SeedDemoData(context.Background(), false, nil)
	require.NoError(t, err)
	require.Equal(t, int64(6), created)

	createdAgain, err := service.SeedDemoData(context.Background(), false, nil)
	require.NoError(t, err)
	require.Equal(t, int64(6), createdAgain)

	count, err := service.DemoAuditCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(6), count)

	logs, err := service.ListByPartition(context.Background(), stableUUID("order", "succeeded"))
	require.NoError(t, err)
	require.NotEmpty(t, logs)
	for _, log := range logs {
		require.NotEmpty(t, log.CurrentHash)
	}
}
