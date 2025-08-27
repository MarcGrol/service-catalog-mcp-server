package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/core"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogconstants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/sloconstants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v", err)
	}
}

func run() error {
	zerolog.TimeFieldFormat = time.RFC3339

	ctx := context.Background()

	// Override if embedded files exist
	serviceCatalogDatabaseFilename, slosDatabaseFilename, databaseCleanup, err := data.UnpackDatabases(ctx)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to unpack databases: %s", err)
	} else {
		defer databaseCleanup()
	}
	cfg := loadConfig(serviceCatalogDatabaseFilename, slosDatabaseFilename)

	var serviceCatalogHandler core.MCPService = nil
	{
		// Initialize catalog repository
		catalogRepo := catalogrepo.New(cfg.PluginConfigs[catalogconstants.CatalogDatabaseFilenameKey])
		err := catalogRepo.Open(ctx)
		if err != nil {
			log.Warn().Msgf("Error opening catalog-database: %v", err)
			return err
		}
		defer catalogRepo.Close(ctx)

		// Initialize catalog search index
		catalogSearchIndex := search.NewSearchIndex(ctx, catalogRepo)

		// Initialize MCP handler
		serviceCatalogHandler = servicecatalog.NewMCPHandler(catalogRepo, catalogSearchIndex)
	}

	var sloHandler core.MCPService = nil
	{
		// Initialize SLO repository
		sloRepo := repo.New(cfg.PluginConfigs[sloconstants.SLODatabaseFilenameKey])
		err := sloRepo.Open(ctx)
		if err != nil {
			return err
		}

		// Initialize slo search index
		sloSearchIndex := slosearch.NewSearchIndex(ctx, sloRepo)

		// Initialize MCP handler
		sloHandler = slo.NewMCPHandler(sloRepo, sloSearchIndex)
	}

	application := core.New(cfg, serviceCatalogHandler, sloHandler)

	cleanup, err := application.Initialize(ctx)
	if err != nil {
		log.Warn().Msgf("Error initializing application: %v", err)
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	err = application.Run()
	if err != nil {
		log.Warn().Msgf("Error running application: %v", err)
		return err
	}
	return nil
}
