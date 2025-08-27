package data

import (
	"context"
	// Needed for embed
	_ "embed"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

//go:embed service-catalog.sqlite
var serviceCatalogDatabase []byte

//go:embed slos.sqlite
var sloDatabase []byte

// UnpackDatabases copies databases from embedding to the filesystem
func UnpackDatabases(c context.Context) (string, string, func(), error) {
	serviceCatalogatabaseFilename, err := copyDatabase("service-catalog.sqlite", serviceCatalogDatabase)
	if err != nil {
		return "", "", nil, err
	}
	sloDatabaseFilename, err := copyDatabase("slos.sqlite", sloDatabase)
	if err != nil {
		return "", "", nil, err
	}

	log.Info().Msgf("Service-catalog-database filename: %s", serviceCatalogatabaseFilename)
	log.Info().Msgf("SLO-database filename: %s", sloDatabaseFilename)

	cleanup := func() {
		log.Info().Msgf("Removing temporary databases %s and %s",
			serviceCatalogatabaseFilename, sloDatabaseFilename)
		err := os.Remove(serviceCatalogatabaseFilename)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to remove service catalog database file %s: %s", serviceCatalogatabaseFilename, err)
		}
		err = os.Remove(sloDatabaseFilename)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to remove service catalog database file %s: %s", sloDatabaseFilename, err)
		}
	}
	return serviceCatalogatabaseFilename, sloDatabaseFilename, cleanup, nil
}

func copyDatabase(name string, databaseBlob []byte) (string, error) {
	// Create temporary file for database.
	tmpDB, err := os.CreateTemp("", name)
	if err != nil {
		return "", fmt.Errorf("error creating database %s: %v", name, err)
	}

	// Write database to file.
	_, err = tmpDB.Write(databaseBlob)
	if err != nil {
		return "", fmt.Errorf("error copying database %s: %v", name, err)
	}

	return tmpDB.Name(), nil
}
