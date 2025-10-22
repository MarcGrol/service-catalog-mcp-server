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

// UnpackServiceCatalogDatabase copies the ServiceCatalog databases from embedding to the filesystem
func UnpackServiceCatalogDatabase(c context.Context) (string, func(), error) {
	serviceCatalogatabaseFilename, err := copyDatabase("service-catalog.sqlite", serviceCatalogDatabase)
	if err != nil {
		return "", nil, err
	}
	log.Info().Msgf("Service-catalog-database filename: %s", serviceCatalogatabaseFilename)

	cleanup := func() {
		log.Info().Msgf("Removing temporary databases %s",
			serviceCatalogatabaseFilename)
		err := os.Remove(serviceCatalogatabaseFilename)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to remove service catalog database file %s: %s", serviceCatalogatabaseFilename, err)
		}
	}
	return serviceCatalogatabaseFilename, cleanup, nil
}

// UnpackSLODatabase copies the SLO databases from embedding to the filesystem
func UnpackSLODatabase(c context.Context) (string, func(), error) {
	sloDatabaseFilename, err := copyDatabase("slos.sqlite", sloDatabase)
	if err != nil {
		return "", nil, err
	}

	log.Info().Msgf("SLO-database filename: %s", sloDatabaseFilename)

	cleanup := func() {
		log.Info().Msgf("Removing temporary databases %s",
			sloDatabaseFilename)
		err = os.Remove(sloDatabaseFilename)
		if err != nil {
			log.Warn().Err(err).Msgf("Failed to remove service catalog database file %s: %s", sloDatabaseFilename, err)
		}
	}
	return sloDatabaseFilename, cleanup, nil
}

func copyDatabase(name string, databaseBlob []byte) (string, error) {
	// Create temporary file for database.
	fp, err := os.CreateTemp("", name)
	if err != nil {
		return "", fmt.Errorf("error creating file for database %s: %v", name, err)
	}

	// Write database to file.
	_, err = fp.Write(databaseBlob)
	if err != nil {
		return "", fmt.Errorf("error copying database %s: %v", name, err)
	}

	return fp.Name(), nil
}
