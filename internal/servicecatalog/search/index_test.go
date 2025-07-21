package search

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

func TestSearchIndex_Search(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	idx := NewSearchIndex(ctx, repo)

	result := idx.Search(ctx, "partner")

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	assert.NoError(t, err)
	t.Logf("Search result:\n %v", string(jsonResult))

	assert.Equal(t, SearchResult{
		Modules: []string{
			"partner",
			"partner-jobs",
			"common/partner",
			"ui/resources/partner",
			"communication/services/partner",
		},
		Teams: []string{
			"PartnerExperience",
			"PartnerExperience_FE",
			"PlatformIntegrationExperience",
		},
		Interfaces: []string{
			"PartnerTermsResourceV1",
			"PartnerReferralResourceV1",
			"PartnerMarketingResourceV1",
			"PartnerDocumentsResourceV1",
			"PartnerOnboardingResourceV1",
		},
		Databases: []string{
			"partner",
		},
	}, result)

}

func setup(t *testing.T) (catalogrepo.Cataloger, context.Context, func()) {
	ctx := context.TODO()

	repo := catalogrepo.New("/Users/marcgrol/src/learnmcp/internal/servicecatalog/service-catalog.sqlite")
	err := repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}
	return repo, ctx, cleanup
}
