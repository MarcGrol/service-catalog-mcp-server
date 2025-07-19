package main

import (
	"context"
	"log"

	"github.com/MarcGrol/learnmcp/internal/app"
	"github.com/MarcGrol/learnmcp/internal/config"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	application := app.New(cfg)

	cleanup, err := application.Initialize(ctx)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	defer cleanup()

	if err := application.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
