package catalogrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/MarcGrol/learnmcp/internal/constants"
	"github.com/stretchr/testify/assert"
)

const doLog = false

func TestListModules(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	modules, err := repo.ListModules(ctx, "")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(modules), 1000)
	assert.LessOrEqual(t, len(modules), 5000)
	if doLog {
		for _, m := range modules {
			t.Logf("%+v", m)
		}
	}
}

func TestListModulesFilered(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	keyword := "kyc"
	if doLog {
		t.Logf("List modules with keyword %s:\n", keyword)
	}

	modules, err := repo.ListModules(ctx, keyword)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(modules), 40)
	assert.LessOrEqual(t, len(modules), 100)
	assert.Equal(t, "onboarding-and-compliance/kyc/webapp/kyc", modules[0].ModuleID)

	if doLog {
		for _, m := range modules {
			t.Logf("%+v", m)
		}
	}
}

func TestListListModulesOfTeam(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	teamID := "CustomerArea"

	if doLog {
		t.Logf("Modules owned by team %s:\n", teamID)
	}

	modules, exists, err := repo.ListModulesOfTeam(ctx, teamID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 100)
	assert.Equal(t, "adyen", modules[0])

	if doLog {
		asJson, _ := json.MarshalIndent(modules, "", "  ")
		t.Logf("%s", asJson)
	}
}

func TestListListModulesOfTeamNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	teamID := "partner"

	if doLog {
		t.Logf("Modules owned by team %s:\n", teamID)
	}

	_, exists, err := repo.ListModulesOfTeam(ctx, teamID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestModuleDetails(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "psp"
	if doLog {
		t.Logf("Details of module %s:\n", interfaceID)
	}

	module, exists, err := repo.GetModuleOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "psp", module.ModuleID)
	assert.Equal(t, "Internal Accounting System", module.Name)
	assert.Equal(t, "Keep track of all transactions from Captured state to Settled state", module.Description)
	assert.Equal(t, "psp/module-metadata.json", module.Spec)

	if doLog {
		asJson, _ := json.MarshalIndent(module, "", "  ")
		t.Logf("%s", asJson)
	}
}

func TestModuleDetailsNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Backoffice"
	if doLog {
		t.Logf("Details of module %s:\n", interfaceID)
	}

	_, exists, err := repo.GetModuleOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListInterfaces(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaces, err := repo.ListInterfaces(ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(interfaces), 10)
	assert.LessOrEqual(t, len(interfaces), 2000)
	assert.Equal(t, "AboutResourceV2", interfaces[0].InterfaceID)

	if doLog {
		for _, m := range interfaces {
			t.Logf("%+v", m)
		}
	}
}

func TestGetInterfaceOnID(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "com.adyen.services.acm.AcmService"

	if doLog {
		fmt.Printf("Details of interface %s:\n", interfaceID)
	}

	module, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.Equal(t, "com.adyen.services.acm.AcmService", module.InterfaceID)
	assert.Equal(t, "ACM", module.Description)
	assert.Equal(t, "RPL", module.Kind)
	assert.Equal(t, "paymentengine/acm/webapp/acm/src/main/resources/rpl-acm.xml", module.Spec)

	if doLog {
		asJson, _ := json.MarshalIndent(module, "", "  ")
		t.Logf("%s", asJson)
	}
}

func TestGetInterfaceOnIDNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Acm"

	if doLog {
		fmt.Printf("Details of interface %s:\n", interfaceID)
	}

	_, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListInterfaceConsumers(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "com.adyen.services.acm.AcmService"
	if doLog {
		fmt.Printf("Modules consuming interface %s:\n", interfaceID)
	}

	modules, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 50)
	assert.Equal(t, "adyen", modules[0])

	if doLog {
		asJson, _ := json.MarshalIndent(modules, "", "  ")
		t.Logf("%s", asJson)
	}
}

func TestListInterfaceConsumersNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	interfaceID := "Acm"
	if doLog {
		fmt.Printf("Modules consuming interface %s:\n", interfaceID)
	}

	_, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestListDatabaseConsumers(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	databaseID := "billing"
	if doLog {
		fmt.Printf("Modules consuming database %s:\n", databaseID)
	}

	modules, exists, err := repo.ListDatabaseConsumers(ctx, databaseID)
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.GreaterOrEqual(t, len(modules), 10)
	assert.LessOrEqual(t, len(modules), 20)
	assert.Equal(t, "airflowjob", modules[0])

	if doLog {
		asJson, _ := json.MarshalIndent(modules, "", "  ")
		t.Logf("%s", asJson)
	}
}

func TestListDatabaseConsumersNotFound(t *testing.T) {
	repo, ctx, cleanup := setup(t)
	defer cleanup()

	databaseID := "bill"
	if doLog {
		fmt.Printf("Modules consuming database %s:\n", databaseID)
	}

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
