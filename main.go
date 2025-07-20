package main

import (
	"context"
	"log"

	"github.com/MarcGrol/learnmcp/internal/app"
	"github.com/MarcGrol/learnmcp/internal/config"
)

func main() {
	// Set log flags for more detailed output
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	ctx := context.Background()

	cfg := config.LoadConfig()

	application := app.New(cfg)

	cleanup, err := application.Initialize(ctx)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
