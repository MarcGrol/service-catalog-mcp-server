package catalogrepo

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

func New(filename string) Cataloger {
	return newCatalogRepo(filename)
}

type CatalogRepo struct {
	filename string
	db       *sqlx.DB
}

func newCatalogRepo(filename string) *CatalogRepo {
	return &CatalogRepo{
		filename: filename,
	}
}

func (repo *CatalogRepo) Open(ctx context.Context) error {
	//log.Printf("Opening database: %s", repo.filename)

	var err error
	if repo.db != nil {
		// already opened
		return nil
	}

	repo.db, err = sqlx.Connect("sqlite", repo.filename)
	if err != nil {
				return fmt.Errorf("connect error: %w", err)
	}

	return nil
}

func (repo *CatalogRepo) Close(ctx context.Context) error {
	//log.Printf("Closing database: %s", repo.filename)
	if repo.db == nil {
		// already closed
		return nil
	}
	return repo.db.Close()
}

func (repo *CatalogRepo) ListModules(ctx context.Context, keyword string) ([]Module, error) {
	if repo.db == nil {
		return nil, fmt.Errorf("database not yet opened")
	}

	if keyword == "" {
		modules := []Module{}
		// This must use module and fails with enriched_module. Don't know why.
		// Currently returns about 2500 entries. Acceptable for now.
		err := repo.db.Select(&modules, "SELECT * FROM module ORDER BY line_count DESC")
		if err != nil {
			if err == sql.ErrNoRows {
				return modules, nil
			}
			return nil, fmt.Errorf("select error: %w", err)
		}

		return enrichWithComplexityScore(modules), nil
	}

	modules := []Module{}
	err := repo.db.Select(&modules, "SELECT * FROM module WHERE module_id LIKE $1 ORDER BY line_count DESC", "%"+keyword+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			return modules, nil
		}
	}

	return enrichWithComplexityScore(modules), nil
}

func (repo *CatalogRepo) ListModulesByCompexity(ctx context.Context, limit int) ([]Module, error) {
	if repo.db == nil {
		return nil, fmt.Errorf("database not yet opened")
	}

	modules := []Module{}
	// This must use module and fails with enriched_module. Don't know why.
	// Currently returns about 2500 entries. Acceptable for now.
	err := repo.db.Select(&modules, "SELECT * FROM enriched_module")
	if err != nil {
		if err == sql.ErrNoRows {
			return modules, nil
		}
		return nil, fmt.Errorf("select error: %w", err)
	}

	for i, module := range modules {
		modules[i].ComplexityScore = module.CalculateComplexityScore()
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].ComplexityScore > modules[j].ComplexityScore
	})

	return modules[0:min(limit, len(modules))], nil
}

func enrichWithComplexityScore(modules []Module) []Module {
	for _, module := range modules {
		module.ComplexityScore = module.CalculateComplexityScore()
	}
	return modules
}

func (repo *CatalogRepo) ListModulesOfTeam(ctx context.Context, id string) ([]string, bool, error) {
	if repo.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	team := ""
	err := repo.db.Get(&team, "SELECT team_id FROM team WHERE team_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select team error: %w", err)
	}

	// Who consume this interface
	modules := []string{}
	err = repo.db.Select(&modules, "SELECT module_id FROM mod_team WHERE team_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %w", err)
	}

	return modules, true, nil
}

func (repo *CatalogRepo) GetModuleOnID(ctx context.Context, id string) (Module, bool, error) {
	if repo.db == nil {
		// already opened
		return Module{}, false, fmt.Errorf("database not yet opened")
	}

	module := Module{}
	err := repo.db.Get(&module, "SELECT * FROM module WHERE module_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return module, false, nil
		}
		return Module{}, false, fmt.Errorf("select module error: %w", err)
	}

	// What kinds?
	err = repo.db.Select(&module.ApplicationKinds, "SELECT kind_id FROM mod_kind WHERE module_id = $1 ORDER BY kind_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select kind error: %w", err)
	}

	//What flows?
	err = repo.db.Select(&module.Flows, "SELECT flow_id FROM mod_flow WHERE module_id = $1 ORDER BY flow_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select flow error: %w", err)
	}

	//What teams?
	err = repo.db.Select(&module.Teams, "SELECT team_id FROM mod_team WHERE module_id = $1 ORDER BY team_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select team error: %w", err)
	}

	// What exposed interfaces?
	err = repo.db.Select(&module.ExposedInterfaces, "SELECT interface_id FROM mod_exposed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select exposed-interfaces error: %w", err)
	}

	// What consumed interfaces?
	err = repo.db.Select(&module.ConsumedInterfaces, "SELECT interface_id FROM mod_consumed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select consumed-interfaces error: %w", err)
	}

	// What databases?
	err = repo.db.Select(&module.Databases, "SELECT database_id FROM mod_database WHERE module_id = $1 ORDER BY database_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select database error: %w", err)
	}

	// What jobs?
	err = repo.db.Select(&module.Jobs, "SELECT job_id FROM mod_job WHERE module_id = $1 ORDER BY job_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select jobs error: %w", err)
	}

	module.ComplexityScore = module.CalculateComplexityScore()

	return module, true, nil
}

func (repo *CatalogRepo) GetInterfaceOnID(ctx context.Context, id string) (Interface, bool, error) {
	if repo.db == nil {
		// already opened
		return Interface{}, false, fmt.Errorf("database not yet opened")
	}

	api := Interface{}
	err := repo.db.Get(&api, "SELECT * FROM interface WHERE interface_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return api, false, nil
		}
		return Interface{}, false, fmt.Errorf("select interface error: %w", err)
	}

	// What methods?
	err = repo.db.Select(&api.Methods, "SELECT method_id FROM interface_method WHERE interface_id = $1 ORDER BY method_id", id)
	if err != nil {
		return Interface{}, false, fmt.Errorf("select meth error: %w", err)
	}

	return api, true, nil
}

func (repo *CatalogRepo) ListInterfaces(ctx context.Context, keyword string) ([]Interface, error) {
	if repo.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	if keyword == "" {

		interfaces := []Interface{}
		err := repo.db.Select(&interfaces, "SELECT * FROM interface ORDER BY interface_id")
		if err != nil {
			if err == sql.ErrNoRows {
				return interfaces, nil
			}
			return interfaces, fmt.Errorf("lit interface error: %w", err)
		}
		return interfaces, nil
	}

	interfaces := []Interface{}
	err := repo.db.Select(&interfaces, "SELECT * FROM interface WHERE interface_id LIKE $1 ORDER BY interface_id", "%"+keyword+"%")
	if err != nil {
		if err == sql.ErrNoRows {
			return interfaces, nil
		}
		return interfaces, fmt.Errorf("list interface error: %w", err)
	}
	return interfaces, nil

}

func (repo *CatalogRepo) ListInterfacesByComplexity(ctx context.Context, limit int) ([]Interface, error) {
	interfaces := []Interface{}
	err := repo.db.Select(&interfaces, "SELECT * FROM interface")
	if err != nil {
		if err == sql.ErrNoRows {
			return interfaces, nil
		}
		return interfaces, fmt.Errorf("list interface error: %w", err)
	}

	sort.Slice(interfaces, func(i, j int) bool {
		return interfaces[i].MethodCount > interfaces[j].MethodCount
	})

	return interfaces[0:min(limit, len(interfaces))], nil
}

// GroupInterfaces is experimentaal and very slow
func (repo *CatalogRepo) GroupInterfaces(ctx context.Context) (map[string][]Interface, error) {
	interfaces, err := repo.ListInterfaces(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error getting interface on ID: %w", err)
	}

	enrichedInterfaces := lo.Map(interfaces, func(item Interface, _ int) *Interface {
		enrichedItem, exists, err := repo.GetInterfaceOnID(ctx, item.InterfaceID)
		if err != nil || !exists {
			return nil
		}
		sort.Strings(enrichedItem.Methods)
		enrichedItem.MethodBasedID = strings.ToLower(strings.Join(enrichedItem.Methods, "-"))
		return &enrichedItem
	})
	enrichedInterfaces = lo.Filter(enrichedInterfaces, func(item *Interface, _ int) bool {
		return item != nil
	})
	interfaces = lo.Map(enrichedInterfaces, func(item *Interface, _ int) Interface {
		return *item
	})

	groupedInterfaces := lo.GroupBy(interfaces, func(item Interface) string {
		return item.MethodBasedID
	})

	return groupedInterfaces, nil
}

func (repo *CatalogRepo) ListInterfaceConsumers(ctx context.Context, id string) ([]string, bool, error) {
	if repo.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	api := ""
	err := repo.db.Get(&api, "SELECT interface_id FROM interface WHERE interface_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select interface error: %w", err)
	}

	// Who consume this interface
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_consumed_interface WHERE interface_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %w", err)
	}

	return interfaces, true, nil
}

func (repo *CatalogRepo) ListDatabaseConsumers(ctx context.Context, id string) ([]string, bool, error) {
	if repo.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	database := ""
	err := repo.db.Get(&database, "SELECT database_id FROM database WHERE database_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select database error: %w", err)
	}

	// Who consume this database
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_database WHERE database_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select databases error: %w", err)
	}

	return interfaces, true, nil
}

func (repo *CatalogRepo) ListDatabases(ctx context.Context) ([]string, error) {
	databases := []string{}
	err := repo.db.Select(&databases, "SELECT DISTINCT database_id FROM database ORDER BY database_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return databases, nil
		}
		return []string{}, fmt.Errorf("select database error: %w", err)
	}
	return databases, nil
}

func (repo *CatalogRepo) ListTeams(ctx context.Context) ([]string, error) {
	teams := []string{}
	err := repo.db.Select(&teams, "SELECT DISTINCT team_id FROM team ORDER BY team_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return teams, nil
		}
		return []string{}, fmt.Errorf("select team error: %w", err)
	}
	return teams, nil
}

func (repo *CatalogRepo) ListFlows(ctx context.Context) ([]string, error) {
	flows := []string{}
	err := repo.db.Select(&flows, "SELECT DISTINCT flow_id FROM flow ORDER BY flow_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return flows, nil
		}
		return []string{}, fmt.Errorf("select flow error: %w", err)
	}
	return flows, nil
}

func (repo *CatalogRepo) ListParticpantsOfFlow(ctx context.Context, id string) ([]string, bool, error) {
	if repo.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	flow := ""
	err := repo.db.Get(&flow, "SELECT flow_id FROM flow WHERE flow_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select flow error: %w", err)
	}

	// Who is part of this flow?
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_flow WHERE flow_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select flows error: %w", err)
	}

	return interfaces, true, nil
}
