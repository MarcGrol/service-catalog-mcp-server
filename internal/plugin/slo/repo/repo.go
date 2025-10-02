package repo

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

// New creates a new sloRepo.
func New(filename string) *sloRepo {
	return newSLORepo(filename)
}

// CatalogRepo is an implementation of Cataloger using a SQLite database.
type sloRepo struct {
	filename string
	db       *sqlx.DB
}

func newSLORepo(filename string) *sloRepo {
	return &sloRepo{
		filename: filename,
	}
}

// Open opens the database connection.
func (r *sloRepo) Open(ctx context.Context) error {
	log.Printf("Opening database: %s", r.filename)

	// Check if the file exists
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

	// Create the SLO table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS SLO (
		UID TEXT NOT NULL PRIMARY KEY,
		CreatedAt TEXT,
		LastModified TEXT, 
		ModificationCount INTEGER NOT NULL DEFAULT 0,
		Filename TEXT NOT NULL,
		DisplayName TEXT,
		Team TEXT NOT NULL,
		Application TEXT NOT NULL,
		Service TEXT NOT NULL,
		Component TEXT NOT NULL,
		Category TEXT NOT NULL,
		RelativeThroughput REAL NOT NULL,
		PromQLQuery TEXT,
		PromQLMetrics TEXT,
		PromQLService TEXT,
		Methods TEXT,
		TargetSLO REAL NOT NULL,
		Duration TEXT,
		SLI REAL NOT NULL,
		DashboardLinkCount INTEGER NOT NULL DEFAULT 0,
		AlertLinkCount INTEGER NOT NULL DEFAULT 0,
		EmailChannelCount INTEGER NOT NULL DEFAULT 0,
		ChatChannelCount INTEGER NOT NULL DEFAULT 0,
		IsEnriched BOOLEAN NOT NULL DEFAULT FALSE,
		IsCritical BOOLEAN NOT NULL DEFAULT FALSE,
		IsFrontdoor BOOLEAN NOT NULL DEFAULT FALSE,
		IsOnlinePaymentsFlow BOOLEAN NOT NULL DEFAULT FALSE,
		IsIPPPaymentsFlow BOOLEAN NOT NULL DEFAULT FALSE,
		IsPayoutFlow BOOLEAN NOT NULL DEFAULT FALSE,
		IsReportingFlow BOOLEAN NOT NULL DEFAULT FALSE,
		IsOnboardingFlow BOOLEAN NOT NULL DEFAULT FALSE,
		IsCustomerPortalFlow BOOLEAN NOT NULL DEFAULT FALSE,
		CriticalFlows TEXT
	);`
	_, err = r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create SLO table: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (r *sloRepo) Close(ctx context.Context) error {
	//log.Printf("Closing database: %s", repo.filename)
	if r.db == nil {
		// already closed
		return nil
	}
	return r.db.Close()
}

// GetSLOByID retrieves a single SLO by its UID.
func (r *sloRepo) GetSLOByID(ctx context.Context, id string) (SLO, bool, error) {
	slo := SLO{}
	err := r.db.Get(&slo, `SELECT * FROM slo WHERE uid = ?`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return SLO{}, false, nil // Not found
		}
		return SLO{}, false, fmt.Errorf("failed to get SLO by ID: %w", err)
	}

	return addMetricsToSLO(slo), true, nil
}

func (r *sloRepo) ListSLOs(ctx context.Context) ([]SLO, error) {
	slos := []SLO{}
	// The row header is also part of the result set, but these have the wrong type, so they break parsing
	// So this is why we filter out the row with 'UID'
	err := r.db.Select(&slos, `SELECT * FROM slo WHERE uid IS NOT 'UID' ORDER BY uid`)
	if err != nil {
		if err == sql.ErrNoRows {
			return slos, nil // Not found
		}
		return nil, fmt.Errorf("failed to select all SLOs: %w", err)
	}
	return slos, nil
}

// listSLOsByTeam retrieves all SLOs for a given team.
func (r *sloRepo) listSLOsByTeam(ctx context.Context, keyword string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT *FROM slo WHERE team LIKE ? ORDER BY uid ASC`, wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by team '%s': %w", keyword, err)
	}

	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// listSLOsByApplication retrieves all SLOs for a given application.
func (r *sloRepo) listSLOsByApplication(ctx context.Context, keyword string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT *	FROM slo WHERE application LIKE ? ORDER BY application`, wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by application '%s': %w", keyword, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// listSLOsByComponent retrieves all SLOs for a given application.
func (r *sloRepo) listSLOsByComponent(ctx context.Context, keyword string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT * FROM slo WHERE component like ? ORDER BY component`, wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by component '%s': %w", keyword, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// listSLOsByService retrieves all SLOs for a given service.
func (r *sloRepo) listSLOsByService(ctx context.Context, keyword string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT *FROM slo WHERE service LIKE ? OR PromQLService LIKE ? ORDER BY service,PromQLService`,
		wildcard(keyword), wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by service '%s': %w", keyword, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// ListSLOsByPromQLService retrieves all SLOs for a given promql-servoce.
func (r *sloRepo) ListSLOsByPromQLService(ctx context.Context, serviceName string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT * FROM slo WHERE PromQLService LIKE ? ORDER BY PromQLService`,
		wildcard(serviceName))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by promQL-service '%s': %w", serviceName, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

func (r *sloRepo) ListSLOsByPromQLModule(ctx context.Context, webappName string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT * FROM slo WHERE PromQLWebapp LIKE ? ORDER BY PromQLWebapp`,
		wildcard(webappName))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by promQL-webapp '%s': %w", webappName, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// listSLOsByMethods retrieves all SLOs for a given service.
func (r *sloRepo) listSLOsByMethods(ctx context.Context, keyword string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT * FROM slo WHERE PromQLMethods LIKE ? ORDER BY PromQLMethods`,
		wildcard(keyword))
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by method '%s': %w", keyword, err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// SearchSLOs searches all SLOs based on category and keyword
func (r *sloRepo) SearchSLOs(ctx context.Context, category, keyword string) ([]SLO, bool, error) {
	switch category {
	case "team":
		return r.listSLOsByTeam(ctx, keyword)
	case "application":
		return r.listSLOsByApplication(ctx, keyword)
	case "webapp", "module":
		return r.ListSLOsByPromQLModule(ctx, keyword)
	case "service", "webservice":
		return r.listSLOsByService(ctx, keyword)
	case "component":
		return r.listSLOsByComponent(ctx, keyword)
	case "method", "methods":
		return r.listSLOsByMethods(ctx, keyword)
	default:
		return nil, false, fmt.Errorf("unknown category: %s", category)
	}

}

func addMetricsToSLOs(slos []SLO) []SLO {
	for i, slo := range slos {
		slos[i] = addMetricsToSLO(slo)
	}

	sort.Slice(slos, func(i, j int) bool {
		return slos[i].BusinessCriticality > slos[j].BusinessCriticality
	})

	return slos
}

func addMetricsToSLO(slo SLO) SLO {
	slo.OperationalReadiness = slo.calculateOperationalReadinessMultiplier()
	slo.BusinessCriticality = slo.calculateBusinessCriticalityMultiplier()
	return slo
}

func wildcard(in string) string {
	if in == "" {
		return in
	}
	return "%" + in + "%"
}
