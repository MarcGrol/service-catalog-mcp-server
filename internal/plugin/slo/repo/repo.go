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

// New creates a new CatalogRepo.
func New(filename string) SLORepo {
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

// ListSLOsByTeam retrieves all SLOs for a given team.
func (r *sloRepo) ListSLOsByTeam(ctx context.Context, teamID string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT *	FROM slo WHERE team = ? ORDER BY uid ASC`, teamID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by team: %w", err)
	}

	return addMetricsToSLOs(slos), len(slos) > 0, nil
}

// ListSLOsByApplication retrieves all SLOs for a given application.
func (r *sloRepo) ListSLOsByApplication(ctx context.Context, id string) ([]SLO, bool, error) {
	slos := []SLO{}
	err := r.db.Select(&slos, `SELECT *	FROM slo WHERE application = ? ORDER BY uid`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []SLO{}, false, nil // Not found
		}
		return nil, false, fmt.Errorf("failed to select SLOs by application: %w", err)
	}
	return addMetricsToSLOs(slos), len(slos) > 0, nil
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
