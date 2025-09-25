package slosearch

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
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
		Services: []string{
			"PartnerTermsResource",
			"PartnerUsersResource",
			"PartnerReferralResource",
			"all-partner-portal",
			"PartnerOnboardingResource",
		},
		Components: []string{
			"partner-flow",
			"v1_webhooks_alelo_partner_order_responses",
			"v1_webhooks_alelo_partner_order_enablements",
			"palauthorisation_internal",
			"capabilityprofilestatusconsumer",
		},
		Methods: []string{
			"/v1/webhooks/Alelo/partner-order/responses",
			"/v1/webhooks/Alelo/partner-order/enablements",
			"getAllCostContractPartnerPricingPlan,getAllPartnerPricingPlanAssignment,getBulkableSettings,updateSettingsForAccounts",
			"accountCanBeBilled,createCostContractPartnerPricingPlan,createPartnerPricingPlanAssignment,deleteCostContractPartnerPricingPlan,deletePartnerPricingPlanAssignment,getBillingInvoiceSettings,getPartnerPricingPlanAssignment,getSettingOptions,kycInfo,updateBillingInvoiceSettings",
			"accountCanBeBilled,createCostContractPartnerPricingPlan,createPartnerPricingPlanAssignment,deleteCostContractPartnerPricingPlan,deletePartnerPricingPlanAssignment,getBillingInvoiceSettings,getPartnerPricingPlanAssignment,getSettingOptions,kycInfo,setDefaultBilling,updateBillingInvoiceSettings",
		},
	}, result)

}

func setup(t *testing.T) (repo.SLORepo, context.Context, func()) {
	ctx := context.TODO()

	sloDatabaseFilename, fileCleanup, err := data.UnpackSLODatabase(ctx)
	assert.NoError(t, err)
	defer fileCleanup()

	repo := repo.New(sloDatabaseFilename)
	err = repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}
	return repo, ctx, cleanup
}
