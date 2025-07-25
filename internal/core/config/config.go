package config

import (
	"flag"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/constants"
)

// Config holds the application's configuration.
type Config struct {
	UseSSE              bool
	UseStreamable       bool
	Port                string
	BaseURL             string
	CatalogDatabaseFile string
	SLODatabaseFile     string
	APIKey              string
}

// LoadConfig loads the application configuration from command-line flags.
func LoadConfig() Config {
	useSSE := flag.Bool("sse", false, "Use SSE transport instead of stdio")
	useStreamable := flag.Bool("http", false, "Use Streamable HTTP transport (easier for testing)")
	port := flag.String("port", "8080", "Port for SSE server")
	baseURL := flag.String("baseurl", "http://localhost", "Base URL for SSE server")
	catalogDatabaseFile := flag.String("catalog-databasefile", constants.CatalogDatabaseFilename, "Full path to the catalog SQLite database file")
	sloDatabaseFile := flag.String("slo-databasefile", constants.SLODatabaseFilename, "Full path to the SLO SQLite database file")
	apiKey := flag.String("api-key", "", "API key for authentication")
	flag.Parse()

	return Config{
		UseSSE:              *useSSE,
		UseStreamable:       *useStreamable,
		Port:                *port,
		BaseURL:             *baseURL,
		CatalogDatabaseFile: *catalogDatabaseFile,
		SLODatabaseFile:     *sloDatabaseFile,
		APIKey:              *apiKey,
	}
}
