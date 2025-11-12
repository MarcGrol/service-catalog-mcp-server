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
	catalog_constants "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/constants"
	catalog_repo "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
	catalog_search "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo"
	slo_constants "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/constants"
	slo_repo "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	slo_search "github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/search"
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
		catalogRepo := catalog_repo.New(cfg.PluginConfigs[catalog_constants.CatalogDatabaseFilenameKey])
		err := catalogRepo.Open(ctx)
		if err != nil {
			log.Warn().Msgf("Error opening catalog-database: %v", err)
			return err
		}
		defer catalogRepo.Close(ctx)

		// Initialize catalog search index
		catalogSearchIndex := catalog_search.NewSearchIndex(ctx, catalogRepo)

		// Initialize MCP handler
		mcpHandlers = append(mcpHandlers, servicecatalog.NewMCPHandler(catalogRepo, catalogSearchIndex))
	}

	if cfg.Mode == config.Both || cfg.Mode == config.SLO {
		// Initialize SLO repository
		sloRepo := slo_repo.New(cfg.PluginConfigs[slo_constants.SLODatabaseFilenameKey])
		err := sloRepo.Open(ctx)
		if err != nil {
			return err
		}

		// Initialize slo search index
		sloSearchIndex := slo_search.NewSearchIndex(ctx, sloRepo)

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
