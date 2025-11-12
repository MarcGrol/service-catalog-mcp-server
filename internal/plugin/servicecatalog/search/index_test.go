package search

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
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
		Methods: []string{
			"OnboardPartner",
			"GetPartnerTerms",
			"InvitePartnerUser",
			"CreatePartnerTerms",
			"AcceptPartnerTerms",
		},
		Flows: []string{},
		Kinds: []string{},
	}, result)

}

func setup(t *testing.T) (repo.Cataloger, context.Context, func()) {
	ctx := context.TODO()

	serviceCatalogDatabaseFilename, fileCleanup, err := data.UnpackServiceCatalogDatabase(ctx)
	assert.NoError(t, err)
	defer fileCleanup()

	repo := repo.New(serviceCatalogDatabaseFilename)

	err = repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}
	return repo, ctx, cleanup
}
