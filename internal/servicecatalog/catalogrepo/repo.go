package catalogrepo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
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
	var err error

	if repo.db != nil {
		// already opened
		return nil
	}
	repo.db, err = sqlx.Connect("sqlite", repo.filename)
	if err != nil {
		return fmt.Errorf("connect error: %s", err)
	}
	return nil
}

func (repo *CatalogRepo) Close(ctx context.Context) error {
	if repo.db == nil {
		// already closed
		return nil
	}
	return repo.db.Close()
}

func (repo *CatalogRepo) ListModules(ctx context.Context) ([]Module, error) {
	if repo.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	modules := []Module{}
	err := repo.db.Select(&modules, "SELECT * FROM enriched_module ORDER BY module_id ASC")
	if err != nil {
		if err == sql.ErrNoRows {
			return modules, nil
		}
		return nil, fmt.Errorf("select error: %s", err)
	}
	return modules, nil
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
		return Module{}, false, fmt.Errorf("select module error: %s", err)
	}

	// What kinds?
	err = repo.db.Select(&module.ApplicationKinds, "SELECT kind_id FROM mod_kind WHERE module_id = $1 ORDER BY kind_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select kind error: %s", err)
	}

	//What flows?
	err = repo.db.Select(&module.Flows, "SELECT flow_id FROM mod_flow WHERE module_id = $1 ORDER BY flow_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select flow error: %s", err)
	}

	//What teams?
	err = repo.db.Select(&module.Teams, "SELECT team_id FROM mod_team WHERE module_id = $1 ORDER BY team_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select team error: %s", err)
	}

	// What exposed interfaces?
	err = repo.db.Select(&module.ExposedInterfaces, "SELECT interface_id FROM mod_exposed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select exposed-interfaces error: %s", err)
	}

	// What consumed interfaces?
	err = repo.db.Select(&module.ConsumedInterfaces, "SELECT interface_id FROM mod_consumed_interface WHERE module_id = $1 ORDER BY interface_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select consumed-interfaces error: %s", err)
	}

	// What databases?
	err = repo.db.Select(&module.Databases, "SELECT database_id FROM mod_database WHERE module_id = $1 ORDER BY database_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select database error: %s", err)
	}

	// What jobs?
	err = repo.db.Select(&module.Jobs, "SELECT job_id FROM mod_job WHERE module_id = $1 ORDER BY job_id", id)
	if err != nil {
		return Module{}, false, fmt.Errorf("select jobs error: %s", err)
	}

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
		return Interface{}, false, fmt.Errorf("select interface error: %s", err)
	}

	// What methods?
	err = repo.db.Select(&api.Methods, "SELECT method_id FROM interface_method WHERE interface_id = $1 ORDER BY method_id", id)
	if err != nil {
		return Interface{}, false, fmt.Errorf("select meth error: %s", err)
	}

	return api, true, nil
}

func (repo *CatalogRepo) ListInterfaces(ctx context.Context) ([]Interface, error) {
	if repo.db == nil {
		// already opened
		return nil, fmt.Errorf("database not yet opened")
	}

	interfaces := []Interface{}
	err := repo.db.Select(&interfaces, "SELECT * FROM interface ORDER BY interface_id")
	if err != nil {
		if err == sql.ErrNoRows {
			return interfaces, nil
		}
		return interfaces, fmt.Errorf("lit interface error: %s", err)
	}

	return interfaces, nil
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
		return []string{}, false, fmt.Errorf("select interface error: %s", err)
	}

	// Who consume this interface
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_consumed_interface WHERE interface_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %s", err)
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
			return []string{}, false, nil
		}
		return []string{}, false, fmt.Errorf("select database error: %s", err)
	}

	// Who consume this database
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_database WHERE database_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %s", err)
	}

	return interfaces, true, nil
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
		return []string{}, false, fmt.Errorf("select team error: %s", err)
	}

	// Who consume this interface
	interfaces := []string{}
	err = repo.db.Select(&interfaces, "SELECT module_id FROM mod_team WHERE team_id = $1 ORDER BY module_id", id)
	if err != nil {
		return []string{}, false, fmt.Errorf("select consumers error: %s", err)
	}

	return interfaces, true, nil
}
