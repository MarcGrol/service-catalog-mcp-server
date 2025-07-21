package main

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/learnmcp/internal/app"
	"github.com/MarcGrol/learnmcp/internal/config"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339

	ctx := context.Background()

	cfg := config.LoadConfig()

	application := app.New(cfg)

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
