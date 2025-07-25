package slosearch

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/sloconstants"
)

func TestSearchIndex_Search(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	idx := NewSearchIndex(ctx, repo)

	result := idx.Search(ctx, "partner", 5)

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	assert.NoError(t, err)
	t.Logf("Search result:\n %v", string(jsonResult))

	assert.Equal(t, Result{
		SLOs: []string{
			"partner_referral_general_latency",
			"partner_onboarding_general_latency",
			"partner_oauth_authorization_latency",
			"partner_commissions_general_latency",
			"partner_referral_general_availability",
		},
		Teams: []string{
			"partner-experience",
			"platform-integration-experience-test-latency",
			"platform-integration-experience-live-latency",
			"payments-engine-alternative-singapore-critical",
			"platform-integration-experience-test-availability",
		},
		Applications: []string{
			"partner",
		},
	}, result)

}

func setup(t *testing.T) (repo.SLORepo, context.Context, func()) {
	ctx := context.TODO()

	repo := repo.New(sloconstants.SLODatabaseFilename())
	err := repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}
	return repo, ctx, cleanup
}
