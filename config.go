package main

import (
	"flag"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/config"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogconstants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/sloconstants"
)

// loadConfig loads the application configuration from command-line flags.
func loadConfig() config.Config {
	useSSE := flag.Bool("sse", false, "Use SSE transport instead of stdio")
	useStreamable := flag.Bool("http", false, "Use Streamable HTTP transport instead of stdio")
	port := flag.String("port", "8080", "Port for SSE server")
	baseURL := flag.String("baseurl", "http://localhost", "Base URL for SSE server")
	catalogDatabaseFile := flag.String("catalog-databasefile", catalogconstants.CatalogDatabaseFilename(), "Full path to the catalog SQLite database file")
	sloDatabaseFile := flag.String("slo-databasefile", sloconstants.SLODatabaseFilename(), "Full path to the SLO SQLite database file")
	apiKey := flag.String("api-key", "", "API key for authentication (default empty)")
	flag.Parse()

	return config.Config{
		UseSSE:        *useSSE,
		UseStreamable: *useStreamable,
		Port:          *port,
		BaseURL:       *baseURL,
		APIKey:        *apiKey,
		PluginConfigs: map[string]string{
			catalogconstants.CatalogDatabaseFilenameKey: *catalogDatabaseFile,
			sloconstants.SLODatabaseFilenameKey:         *sloDatabaseFile,
		},
	}
}
