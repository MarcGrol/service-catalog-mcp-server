package main

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	zerolog.TimeFieldFormat = time.RFC3339

	ctx := context.Background()

	cfg := loadConfig()

	var serviceCatalogHandler core.MCPService = nil
	{
		// Initialize catalog repository
		catalogRepo := catalogrepo.New(cfg.PluginConfigs[catalogconstants.CatalogDatabaseFilenameKey])
		err := catalogRepo.Open(ctx)
		if err != nil {
			log.Fatal().Msgf("Error opening catalog-database: %v", err)
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
			log.Fatal().Msgf("Error opening SLO database: %v", err)
		}

		// Initialize slo search index
		sloSearchIndex := slosearch.NewSearchIndex(ctx, sloRepo)

		// Initialize MCP handler
		sloHandler = slo.NewMCPHandler(sloRepo, sloSearchIndex)
	}

	application := core.New(cfg, serviceCatalogHandler, sloHandler)

	cleanup, err := application.Initialize(ctx)
	if err != nil {
		log.Fatal().Msgf("Error initializing application: %v", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	if err := application.Run(); err != nil {
		log.Fatal().Msgf("Error running application: %v", err)
	}
}
