package usecase

import (
	"context"
	"testing"

	"RuntimeRoasters/apps/trace-service/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestSeedDemoDataBuildsDeterministicTraceDocuments(t *testing.T) {
	db := setupTraceTestDB(t)
	service := NewService(db, nil)

	created, err := service.SeedDemoData(context.Background(), false, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, created, int64(12))

	createdAgain, err := service.SeedDemoData(context.Background(), false, nil)
	require.NoError(t, err)
	require.GreaterOrEqual(t, createdAgain, int64(12))

	count, err := service.DemoTraceDocumentCount(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(4), count)

	var docs []domain.TraceDocument
	require.NoError(t, db.Order("entity_id ASC").Find(&docs).Error)
	require.Len(t, docs, 4)
	for _, doc := range docs {
		require.NotEmpty(t, doc.Document)
		require.NotEmpty(t, doc.TraceIDs)
	}
}
