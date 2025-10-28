package main

import (
	"flag"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/config"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogconstants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/sloconstants"
)

// loadConfig loads the application configuration from command-line flags.
func loadConfig(catalogDatabaseFilename, sloDatabaseFilename string) config.Config {
	useSSE := flag.Bool("sse", false, "Use SSE transport instead of stdio")
	useStreamable := flag.Bool("http", false, "Use Streamable HTTP transport instead of stdio")
	port := flag.String("port", "8080", "Port for SSE server")
	baseURL := flag.String("baseurl", "http://localhost", "Base URL for SSE server")
	catalogDatabaseFile := flag.String("catalog-databasefile", catalogDatabaseFilename, "Full path to the catalog SQLite database file")
	sloDatabaseFile := flag.String("slo-databasefile", sloDatabaseFilename, "Full path to the SLO SQLite database file")
	apiKey := flag.String("api-key", "", "API key for authentication (default empty)")
	mode := flag.String("mode", "both", "slo, service-catalog or both")
	flag.Parse()

	return config.Config{
		UseSSE:        *useSSE,
		UseStreamable: *useStreamable,
		Port:          *port,
		BaseURL:       *baseURL,
		APIKey:        *apiKey,
		Mode:          config.Mode(*mode),
		PluginConfigs: map[string]string{
			catalogconstants.CatalogDatabaseFilenameKey: *catalogDatabaseFile,
			sloconstants.SLODatabaseFilenameKey:         *sloDatabaseFile,
		},
	}
}
