package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/config"
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
	serviceCatalogDatabaseFilename, serviceCatalogDatabaseCleanup, err := data.UnpackServiceCatalogDatabase(ctx)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to unpack service-catalog database: %s", err)
		return err
	}
	defer serviceCatalogDatabaseCleanup()

	slosDatabaseFilename, sloDatabaseCleanup, err := data.UnpackSLODatabase(ctx)
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to unpack slo database: %s", err)
		return err
	}
	defer sloDatabaseCleanup()

	cfg := loadConfig(serviceCatalogDatabaseFilename, slosDatabaseFilename)

	mcpHandlers := []core.MCPService{}
	if cfg.Mode == config.Both || cfg.Mode == config.ServiceCatalog {
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
		mcpHandlers = append(mcpHandlers, servicecatalog.NewMCPHandler(catalogRepo, catalogSearchIndex))
	}

	if cfg.Mode == config.Both || cfg.Mode == config.SLO {
		// Initialize SLO repository
		sloRepo := repo.New(cfg.PluginConfigs[sloconstants.SLODatabaseFilenameKey])
		err := sloRepo.Open(ctx)
		if err != nil {
			return err
		}

		// Initialize slo search index
		sloSearchIndex := slosearch.NewSearchIndex(ctx, sloRepo)

		// Initialize MCP handler
		mcpHandlers = append(mcpHandlers, slo.NewMCPHandler(sloRepo, sloSearchIndex))
	}

	application := core.New(cfg, mcpHandlers)

	applicationCleanup, err := application.Initialize(ctx)
	if err != nil {
		log.Warn().Msgf("Error initializing application: %v", err)
		return err
	}
	if applicationCleanup != nil {
		defer applicationCleanup()
	}

	err = application.Run()
	if err != nil {
		log.Warn().Msgf("Error running application: %v", err)
		return err
	}
	return nil
}
