package main

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/app"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/config"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339

	ctx := context.Background()

	cfg := config.LoadConfig()

	// Initialize catalog repository
	catalogRepo := catalogrepo.New(cfg.DatabaseFile)
	err := catalogRepo.Open(ctx)
	if err != nil {
		log.Fatal().Msgf("Error opening database: %v", err)
	}
	defer catalogRepo.Close(ctx)

	// Initialize search index
	searchIndex := search.NewSearchIndex(ctx, catalogRepo)

	// Initialize MCP handler
	mcpHandler := servicecatalog.NewMCPHandler(catalogRepo, searchIndex)

	application := app.New(cfg, mcpHandler)

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
