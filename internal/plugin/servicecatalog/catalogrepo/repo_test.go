package catalogrepo

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/constants"
)

func TestListModules(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	modules, err := repo.ListModules(ctx, "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(modules), 1000)
	assert.LessOrEqual(t, len(modules), 5000)
	assert.Equal(t, "psp", modules[0].ModuleID)
}

func TestListModulesByComplexity(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	modules, err := repo.ListModulesByCompexity(ctx, 5)
	assert.NoError(t, err)
	assert.Len(t, modules, 5)
	top5 := lo.Map(modules, func(m Module, _ int) string {
		return m.ModuleID
	})

	assert.Equal(t, []string{"psp",
		"onboarding-and-compliance/kyc/webapp/kyc",
		"paymentengine/acm/webapp/acm",
		"bcm-container/protocol/step2/sct", "vias"}, top5)
}

func TestListModulesFilered(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	keyword := "kyc"

	modules, err := repo.ListModules(ctx, keyword)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(modules), 40)
	assert.LessOrEqual(t, len(modules), 100)
	assert.Equal(t, "onboarding-and-compliance/kyc/webapp/kyc", modules[0].ModuleID)

}

func TestListListModulesOfTeam(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	teamID := "CustomerArea"

	modules, exists, err := repo.ListModulesOfTeam(ctx, teamID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 100)
	assert.Equal(t, "adyen", modules[0])

}

func TestListListModulesOfTeamNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	teamID := "partner"

	_, exists, err := repo.ListModulesOfTeam(ctx, teamID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestModuleDetails(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "psp"

	module, exists, err := repo.GetModuleOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "psp", module.ModuleID)
	assert.Equal(t, "Internal Accounting System", module.Name)
	assert.Equal(t, "Keep track of all transactions from Captured state to Settled state", module.Description)
	assert.Equal(t, "psp/module-metadata.json", module.Spec)
}

func TestModuleDetailsNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Backoffice"

	_, exists, err := repo.GetModuleOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListInterfaces(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaces, err := repo.ListInterfaces(ctx, "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(interfaces), 10)
	assert.LessOrEqual(t, len(interfaces), 2000)
	assert.Equal(t, "threesixty/compass/webapp/compass", interfaces[0].ModuleID)
	assert.Equal(t, "AboutResourceV2", interfaces[0].InterfaceID)

}

func TestListInterfacesFiltered(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaces, err := repo.ListInterfaces(ctx, "Partner")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(interfaces), 1)
	assert.LessOrEqual(t, len(interfaces), 20)
	assert.Equal(t, "partner", interfaces[0].ModuleID)
	assert.Equal(t, "PartnerDocumentsResourceV1", interfaces[0].InterfaceID)

}

func TestListInterfacesByComplexity(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaces, err := repo.ListInterfacesByComplexity(ctx, 5)
	assert.NoError(t, err)
	assert.Len(t, interfaces, 5)
	names := lo.Map(interfaces, func(i Interface, _ int) string { return i.InterfaceID })
	assert.Equal(t, []string{
		"ManagementServiceV3",
		"ManagementServiceV1",
		"com.adyen.services.postfm.PosTFMService",
		"com.adyen.services.acm.AcmService",
		"com.adyen.services.configurationapi.MerchantConfigurationService"}, names)
}

func TestGetInterfaceOnID(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	{
		interfaceID := "com.adyen.services.acm.AcmService"
		module, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, "paymentengine/acm/webapp/acm", module.ModuleID)
		assert.Equal(t, interfaceID, module.InterfaceID)
		assert.Equal(t, "ACM", module.Description)
		assert.Equal(t, "RPL", module.Kind)
		assert.Nil(t, module.OpenAPISpecs)
		assert.Equal(t, "paymentengine/acm/webapp/acm/src/main/resources/rpl-acm.xml", *module.RPLSpecs)
	}
	{
		interfaceID := "com.adyen.services.checkout.shopper.BinLookupService"
		iface, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, "checkoutshopper", iface.ModuleID)
		assert.Equal(t, interfaceID, iface.InterfaceID)
		assert.Equal(t, "BinLookupService", iface.Description)
		assert.Equal(t, "RPL", iface.Kind)
		assert.Nil(t, iface.OpenAPISpecs)
		assert.Equal(t, "checkoutshopper/src/main/resources/rpl-checkoutBinLookup.xml", *iface.RPLSpecs)
	}
	{
		interfaceID := "ManagementServiceV3"
		iface, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, "configurationapi", iface.ModuleID)
		assert.Equal(t, interfaceID, iface.InterfaceID)
		assert.Equal(t, "ManagementService", iface.Description)
		assert.Equal(t, "OpenAPI", iface.Kind)
		assert.Equal(t, "configurationapi/open-api-specs/prod/rpl/ManagementService-v3.json", *iface.OpenAPISpecs)
		assert.Nil(t, iface.RPLSpecs)
	}
}

func TestGetInterfaceOnIDNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Acm"

	_, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListInterfaceConsumers(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "com.adyen.services.acm.AcmService"

	modules, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 50)
	assert.Equal(t, "adyen", modules[0])
}

func TestListInterfaceConsumersNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Acm"

	_, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListDatabaseConsumers(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	databaseID := "billing"

	modules, exists, err := repo.ListDatabaseConsumers(ctx, databaseID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 20)
	assert.Equal(t, "airflowjob", modules[0])

}

func TestListDatabaseConsumersNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	databaseID := "bill"

	_, exists, err := repo.ListDatabaseConsumers(ctx, databaseID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListDatabases(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	databases, err := repo.ListDatabases(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(databases), 100)
	assert.LessOrEqual(t, len(databases), 500)
	assert.Equal(t, "a2aissuer", databases[0])
}

func TestListTeams(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	teams, err := repo.ListTeams(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(teams), 100)
	assert.LessOrEqual(t, len(teams), 500)
	assert.Equal(t, "AMLTech", teams[0])
}

func TestListFlows(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	flows, err := repo.ListFlows(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"CustomerPortals-TransactionSearch",
		"IPP_Payments-Authorization",
		"IPP_Payments-ModificationAndSettlement",
		"Onboarding-CreateAccountHolder",
		"Onboarding-CreateLegalEntity",
		"Onboarding-ModifyLegalEntity",
		"Onboarding-RequestAccountHolderCapabilities",
		"Online_Payments-Authorization",
		"Online_Payments-ModificationAndSettlement",
		"Payments-ModificationAndSettlement",
		"Payout-OnDemandPayout",
		"Payout-Sweep",
	}, flows)
}

func TestListFlowParticpants(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	modules, exists, err := repo.ListParticpantsOfFlow(ctx, "CustomerPortals-TransactionSearch")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, []string{"ca", "ca-core", "consumers", "pspdw"}, modules)
}

func TestListKinds(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	kinds, err := repo.ListKinds(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"?", "api", "cache-api", "cache-db", "common",
		"communicator", "constants", "consumer", "database", "domain", "entities",
		"flink-job", "framework", "frontend", "integration-test", "ipp-terminal",
		"job", "repositories", "services", "services_common", "terminal-component",
		"ui-beans-common", "ui-resources", "ui-resources-common", "unknown", "util",
		"webapp", "webapp_external", "webapp_internal"}, kinds)
}

func TestListAppsWithKind(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	kinds, exists, err := repo.ListModulesWithKind(ctx, "webapp_external")
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, []string{
		"acs/webapps/acs", "apipix/apipix", "attestpos"}, kinds[0:3])
}

func setup(t *testing.T) (Cataloger, context.Context, func()) {
	ctx := context.TODO()

	repo := New(constants.DatabaseFilename)
	err := repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}
	return repo, ctx, cleanup
}
