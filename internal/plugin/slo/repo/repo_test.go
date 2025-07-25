package repo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	repo, ctx, cleanup := createRealDatabase(t)
	defer cleanup()

	// Test GetSLOByID
	t.Run("GetSLOByID", func(t *testing.T) {
		slo, found, err := repo.GetSLOByID(ctx, "accessportal_accessportal_main_availability")
		assert.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "accessportal_accessportal_main_availability", slo.UID)
		assert.Equal(t, "Access Portal - Availability", slo.DisplayName)
		assert.Equal(t, 1.05, slo.OperationalReadiness)
		assert.Equal(t, 1.0, slo.BusinessCriticality)

		// Test not found
		_, found, err = repo.GetSLOByID(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, found)
	})

	// Test ListSLOs
	t.Run("ListSLOs", func(t *testing.T) {
		slos, err := repo.ListSLOs(ctx)
		assert.NoError(t, err)

		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "accessportal_accessportal_main_availability", slos[0].UID)
		assert.Equal(t, "be-internal_services", slos[0].Team)
		assert.Equal(t, "accessportal", slos[0].Application)
		assert.Equal(t, "accessportal", slos[0].Service)
		assert.Equal(t, "main", slos[0].Component)
		assert.Equal(t, "Availability", slos[0].Category)
	})

	// Test ListSLOsByTeam
	t.Run("ListSLOsByTeam", func(t *testing.T) {
		slos, err := repo.ListSLOsByTeam(ctx, "be-internal_services")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "accessportal_accessportal_main_availability", slos[0].UID)
		assert.Equal(t, 1.05, slos[0].OperationalReadiness)
		assert.Equal(t, 1.0, slos[0].BusinessCriticality)

		slos, err = repo.ListSLOsByTeam(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.Len(t, slos, 0)
	})

	// Test ListSLOsByApplication
	t.Run("ListSLOsByApplication", func(t *testing.T) {
		slos, err := repo.ListSLOsByApplication(ctx, "accessportal")
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "accessportal_accessportal_main_availability", slos[0].UID)
		assert.Equal(t, 1.05, slos[0].OperationalReadiness)
		assert.Equal(t, 1.0, slos[0].BusinessCriticality)

		slos, err = repo.ListSLOsByApplication(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.Len(t, slos, 0)
	})
}

func createRealDatabase(t *testing.T) (SLORepo, context.Context, func()) {
	ctx := context.Background()
	repo := New("/Users/marcgrol/src/service-catalog-mcp-server/internal/plugin/slo/slos.sqlite")
	err := repo.Open(ctx)
	assert.NoError(t, err)
	return repo, ctx, func() {
		repo.Close(ctx)
	}
}
