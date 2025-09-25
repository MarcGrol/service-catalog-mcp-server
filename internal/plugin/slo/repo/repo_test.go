package repo

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
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
		assert.Equal(t, 0.0, slo.BusinessCriticality)

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
		assert.Equal(t, "a2apayments_a2aissuer-api_test-ideal-certificates_availability", slos[0].UID)
		assert.Equal(t, "open-banking", slos[0].Team)
		assert.Equal(t, "a2apayments", slos[0].Application)
		assert.Equal(t, "a2aissuer-api", slos[0].Service)
		assert.Equal(t, "test-ideal-certificates", slos[0].Component)
		assert.Equal(t, "Availability", slos[0].Category)
	})

	// Test ListSLOsByTeam
	t.Run("ListSLOsByTeam", func(t *testing.T) {
		slos, exists, err := repo.listSLOsByTeam(ctx, "be-internal_services")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "accessportal_accessportal_main_availability", slos[0].UID)
		assert.Equal(t, 1.05, slos[0].OperationalReadiness)
		assert.Equal(t, 0.0, slos[0].BusinessCriticality)

		slos, exists, err = repo.listSLOsByTeam(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Len(t, slos, 0.0)
	})

	// Test ListSLOsByApplication
	t.Run("ListSLOsByApplication", func(t *testing.T) {
		slos, exists, err := repo.listSLOsByApplication(ctx, "accessportal")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "accessportal_accessportal_main_availability", slos[0].UID)
		assert.Equal(t, 1.05, slos[0].OperationalReadiness)
		assert.Equal(t, 0.0, slos[0].BusinessCriticality)

		slos, exists, err = repo.listSLOsByApplication(ctx, "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Len(t, slos, 0)
	})

	// Test SearchSLOs
	t.Run("SearchSLOs", func(t *testing.T) {
		slos, exists, err := repo.SearchSLOs(ctx, "methods", "generateDataInsightsLink")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.GreaterOrEqual(t, len(slos), 1)
		assert.Equal(t, "obgateway_obgateway_results_availability", slos[0].UID)
		assert.InDelta(t, 1.15, slos[0].OperationalReadiness, 0.001)
		assert.Equal(t, 0.0, slos[0].BusinessCriticality)

		slos, exists, err = repo.SearchSLOs(ctx, "methods", "nonexistent")
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.Len(t, slos, 0)
	})
}

func createRealDatabase(t *testing.T) (*sloRepo, context.Context, func()) {
	ctx := context.Background()

	sloDatabaseFilename, fileCleanup, err := data.UnpackSLODatabase(ctx)
	assert.NoError(t, err)
	defer fileCleanup()

	repo := New(sloDatabaseFilename)
	err = repo.Open(ctx)
	assert.NoError(t, err)

	log.Printf("Opened %s", sloDatabaseFilename)

	return repo, ctx, func() {
		log.Printf("Cleanup %s", sloDatabaseFilename)
		repo.Close(ctx)
	}
}
