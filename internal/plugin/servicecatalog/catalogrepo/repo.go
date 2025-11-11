package catalogrepo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"

	_ "github.com/glebarez/go-sqlite" // sqlite driver
	"github.com/jmoiron/sqlx"
)

// New creates a new Cataloger instance.
func New(filename string) Cataloger {
	return newCatalogRepo(filename)
}

// CatalogRepo is an implementation of Cataloger using a SQLite database.
type CatalogRepo struct {
	filename string
	db       *sqlx.DB
}

func newCatalogRepo(filename string) *CatalogRepo {
	return &CatalogRepo{
		filename: filename,
	}
}

// Open opens the database connection.
func (r *CatalogRepo) Open(ctx context.Context) error {
	log.Printf("Opening database: %s", r.filename)

	_, err := os.Stat(r.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s must exist", r.filename)
		}
		return fmt.Errorf("Error opening file %s: %s", r.filename, err)
	}

	if r.db != nil {
		// already opened
		return nil
	}

	r.db, err = sqlx.Connect("sqlite", r.filename)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (r *CatalogRepo) Close(ctx context.Context) error {
	//log.Printf("Closing database: %s", repo.filename)
	if r.db == nil {
		// already closed
		return nil
	}
	return r.db.Close()
}

// ListModules lists modules based on a keyword.
func (r *CatalogRepo) ListModules(ctx context.Context, keyword string) ([]Module, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not yet opened")
	}

	if keyword == "" {
		modules := []Module{}
		// This must use module and fails with enriched_module. Don't know why.
		// Currently returns about 2500 entries. Acceptable for now.
		err := r.db.Select(&modules, "SELECT * FROM module ORDER BY line_count DESC")
		if err != nil {
			if err == sql.ErrNoRows {
				return modules, nil
			}
			return nil, fmt.Errorf("select error: %w", err)
		}

		return enrichWithComplexityScore(modules), nil
	}

	modules := []Module{}
	err := r.db.Select(&modules, "SELECT * FROM module WHERE module_id LIKE $1 ORDER BY line_count DESC",
		wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return modules, nil
		}
	}

	return enrichWithComplexityScore(modules), nil
}

// ListModulesByCompexity lists modules ordered by complexity.
func (r *CatalogRepo) ListModulesByCompexity(ctx context.Context, limit int) ([]Module, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database not yet opened")
	}

	modules := []Module{}
	// This must use module and fails with enriched_module. Don't know why.
	// Currently returns about 2500 entries. Acceptable for now.
	err := r.db.Select(&modules, "SELECT * FROM enriched_module")
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

// ListModulesOfTeam lists modules belonging to a specific team.
func (r *CatalogRepo) ListModulesOfTeam(ctx context.Context, id string) ([]string, bool, error) {
	if r.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	team := ""
	err := r.db.Get(&team, "SELECT team_id FROM team WHERE team_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select team error: %w", err)
	}

	// Who consume this interface
	modules := []string{}
	err = r.db.Select(&modules, "SELECT module_id FROM mod_team WHERE team_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %w", err)
	}

	return modules, true, nil
}

// GetModuleOnID retrieves a module by its ID.
func (r *CatalogRepo) GetModuleOnID(ctx context.Context, id string) (Module, bool, error) {
	if r.db == nil {
		// already opened
		return Module{}, false, fmt.Errorf("database not yet opened")
	}

	module := Module{}
	err := r.db.Get(&module, "SELECT * FROM module WHERE module_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return module, false, nil
		}
		return Module{}, false, fmt.Errorf("get module error: %w", err)
	}

	// What kinds?
	err = r.db.Select(&module.ApplicationKinds, "SELECT kind_id FROM mod_kind WHERE module_id = $1 ORDER BY kind_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select kind error: %w", err)
	}
	module.KindCount = intPointer(len(module.ApplicationKinds))

	//What flows?
	err = r.db.Select(&module.Flows, "SELECT flow_id FROM mod_flow WHERE module_id = $1 ORDER BY flow_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select flow error: %w", err)
	}
	module.FlowCount = intPointer(len(module.Flows))

	//What teams?
	err = r.db.Select(&module.Teams, "SELECT team_id FROM mod_team WHERE module_id = $1 ORDER BY team_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select team error: %w", err)
	}
	module.TeamCount = intPointer(len(module.Teams))

	// What exposed interfaces?
	err = r.db.Select(&module.ExposedInterfaces, "SELECT interface_id FROM mod_exposed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select exposed-interfaces error: %w", err)
	}
	module.ExposedAPICount = intPointer(len(module.ExposedInterfaces))

	// What consumed interfaces?
	err = r.db.Select(&module.ConsumedInterfaces, "SELECT interface_id FROM mod_consumed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select consumed-interfaces error: %w", err)
	}
	module.ConsumedAPICount = intPointer(len(module.ConsumedInterfaces))

	// What databases?
	err = r.db.Select(&module.Databases, "SELECT database_id FROM mod_database WHERE module_id = $1 ORDER BY database_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select database error: %w", err)
	}
	module.DatabaseCount = intPointer(len(module.Databases))

	// What jobs?
	err = r.db.Select(&module.Jobs, "SELECT job_id FROM mod_job WHERE module_id = $1 ORDER BY job_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select jobs error: %w", err)
	}
	module.JobCount = intPointer(len(module.Jobs))

	module.ComplexityScore = module.CalculateComplexityScore()

	return module, true, nil
}

func intPointer(val int) *int {
	return &val
}

// GetInterfaceOnID retrieves an interface by its ID.
func (r *CatalogRepo) GetInterfaceOnID(ctx context.Context, id string) (Interface, bool, error) {
	if r.db == nil {
		// already opened
		return Interface{}, false, fmt.Errorf("database not yet opened")
	}

	api := Interface{}
	err := r.db.Get(&api, `
		SELECT
			m.module_id, i.interface_id, i.description, i.kind, i.openapi_specification, i.rpl_specification, i.method_count
		FROM
			enriched_interface i
			LEFT JOIN mod_exposed_interface m ON i.interface_id = m.interface_id
		WHERE i.interface_id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return api, false, nil
		}
		return Interface{}, false, fmt.Errorf("select interface error: %w", err)
	}

	// What methods?
	err = r.db.Select(&api.Methods, "SELECT method_id FROM interface_method WHERE interface_id LIKE $1 ORDER BY method_id", id)
	if err != nil {
		return Interface{}, false, fmt.Errorf("select meth error: %w", err)
	}

	return api, true, nil
}

// ListInterfaces lists interfaces based on a keyword.
func (r *CatalogRepo) ListInterfaces(ctx context.Context, keyword string) ([]Interface, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	if keyword == "" {

		interfaces := []Interface{}
		err := r.db.Select(&interfaces, `
	SELECT 
		m.module_id, i.interface_id, i.description, i.kind, i.openapi_specification, i.rpl_specification, i.method_count
	FROM 
		enriched_interface i
		LEFT JOIN mod_exposed_interface m ON i.interface_id = m.interface_id 
	ORDER BY 
		i.interface_id`)
		if err != nil {
			if err == sql.ErrNoRows {
				return interfaces, nil
			}
			return interfaces, fmt.Errorf("list interface error: %w", err)
		}
		return interfaces, nil
	}

	interfaces := []Interface{}
	err := r.db.Select(&interfaces, `
	SELECT 
		m.module_id, i.interface_id, i.description, i.kind, i.openapi_specification, i.rpl_specification, i.method_count
	FROM 
		enriched_interface i
		LEFT JOIN mod_exposed_interface m ON i.interface_id = m.interface_id 
	WHERE 
		i.interface_id LIKE $1
	ORDER BY 
		i.interface_id`, wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return interfaces, nil
		}
		return interfaces, fmt.Errorf("list interface error: %w", err)
	}
	return interfaces, nil

}

// ListInterfacesByComplexity lists interfaces ordered by complexity.
func (r *CatalogRepo) ListInterfacesByComplexity(ctx context.Context, limit int) ([]Interface, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	interfaces := []Interface{}
	err := r.db.Select(&interfaces, `
	SELECT 
		m.module_id, i.interface_id, i.description, i.kind, i.openapi_specification, i.rpl_specification, i.method_count
	FROM 
		enriched_interface i
		LEFT JOIN mod_exposed_interface m ON i.interface_id = m.interface_id 
	ORDER BY 
		i.method_count DESC LIMIT $1`, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return interfaces, nil
		}
		return interfaces, fmt.Errorf("list interface error: %w", err)
	}

	return interfaces[0:min(limit, len(interfaces))], nil
}

// ListInterfaceConsumers lists modules that consume a given interface.
func (r *CatalogRepo) ListInterfaceConsumers(ctx context.Context, id string) ([]string, bool, error) {
	if r.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	api := ""
	err := r.db.Get(&api, "SELECT interface_id FROM interface WHERE interface_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select interface error: %w", err)
	}

	// Who consume this interface
	interfaces := []string{}
	err = r.db.Select(&interfaces, "SELECT module_id FROM mod_consumed_interface WHERE interface_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %w", err)
	}

	return interfaces, true, nil
}

// ListMethods lists all web-methods.
func (r *CatalogRepo) ListMethods(ctx context.Context) ([]string, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	methods := []string{}
	err := r.db.Select(&methods, "SELECT DISTINCT method_id FROM interface_method ORDER BY method_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return methods, nil
		}
		return []string{}, fmt.Errorf("select interfaces error: %w", err)
	}
	return methods, nil
}

// ListDatabaseConsumers lists modules that consume a given database.
func (r *CatalogRepo) ListDatabaseConsumers(ctx context.Context, id string) ([]string, bool, error) {
	if r.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	database := ""
	err := r.db.Get(&database, "SELECT database_id FROM database WHERE database_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select database error: %w", err)
	}

	// Who consume this database
	interfaces := []string{}
	err = r.db.Select(&interfaces, "SELECT module_id FROM mod_database WHERE database_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select databases error: %w", err)
	}

	return interfaces, true, nil
}

// ListDatabases lists all databases.
func (r *CatalogRepo) ListDatabases(ctx context.Context) ([]string, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	databases := []string{}
	err := r.db.Select(&databases, "SELECT DISTINCT database_id FROM database ORDER BY database_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return databases, nil
		}
		return []string{}, fmt.Errorf("select database error: %w", err)
	}
	return databases, nil
}

// ListTeams lists all teams.
func (r *CatalogRepo) ListTeams(ctx context.Context) ([]string, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	teams := []string{}
	err := r.db.Select(&teams, "SELECT DISTINCT team_id FROM team ORDER BY team_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return teams, nil
		}
		return []string{}, fmt.Errorf("select team error: %w", err)
	}
	return teams, nil
}

// ListFlows lists all flows.
func (r *CatalogRepo) ListFlows(ctx context.Context) ([]string, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	flows := []string{}
	err := r.db.Select(&flows, "SELECT DISTINCT flow_id FROM flow ORDER BY flow_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return flows, nil
		}
		return []string{}, fmt.Errorf("select flow error: %w", err)
	}
	return flows, nil
}

// ListParticpantsOfFlow lists modules participating in a given flow.
func (r *CatalogRepo) ListParticpantsOfFlow(ctx context.Context, id string) ([]string, bool, error) {
	if r.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	flow := ""
	err := r.db.Get(&flow, "SELECT flow_id FROM flow WHERE flow_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select flow error: %w", err)
	}

	// Who is part of this flow?
	interfaces := []string{}
	err = r.db.Select(&interfaces, "SELECT module_id FROM mod_flow WHERE flow_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select modules of flow error: %w", err)
	}

	return interfaces, true, nil
}

// ListKinds lists all module kinds.
func (r *CatalogRepo) ListKinds(ctx context.Context) ([]string, error) {
	if r.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	flows := []string{}
	err := r.db.Select(&flows, "SELECT DISTINCT kind_id FROM kind ORDER BY kind_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return flows, nil
		}
		return []string{}, fmt.Errorf("select kind error: %w", err)
	}
	return flows, nil
}

// ListModulesWithKind lists modules of a specific kind.
func (r *CatalogRepo) ListModulesWithKind(ctx context.Context, id string) ([]string, bool, error) {
	if r.db == nil {
		// already opened
		return nil, false, fmt.Errorf("database not yet opened")
	}

	kind := ""
	err := r.db.Get(&kind, "SELECT kind_id FROM kind WHERE kind_id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			// not found, do return others with similar names
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select kind error: %w", err)
	}

	// Which application are of this kind
	interfaces := []string{}
	err = r.db.Select(&interfaces, "SELECT module_id FROM mod_kind WHERE kind_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select modules with kind error: %w", err)
	}

	return interfaces, true, nil
}

func wildcard(in string) string {
	if in == "" {
		return in
	}
	return "%" + in + "%"
}
