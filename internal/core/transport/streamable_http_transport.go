package transport

import (
	"net/http"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// StreamableHTTPTransport implements the ServerTransport interface for streamable HTTP.
type StreamableHTTPTransport struct {
	mcpServer *server.MCPServer
	port      string
	apiKey    string
}

// NewStreamableHTTPTransport creates a new StreamableHTTPTransport instance.
func NewStreamableHTTPTransport(s *server.MCPServer, port, apiKey string) *StreamableHTTPTransport {
	return &StreamableHTTPTransport{
		mcpServer: s,
		port:      port,
		apiKey:    apiKey,
	}
}

// Start starts the streamable HTTP transport server.
func (t *StreamableHTTPTransport) Start() error {
	streamableServer := server.NewStreamableHTTPServer(t.mcpServer,
		server.WithStateLess(true),
	)

	// Create a new HTTP server
	httpServer := &http.Server{
		Addr: ":" + t.port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// API Key authentication
			if t.apiKey != "" {
				providedAPIKey := r.Header.Get("X-API-Key")
				if providedAPIKey == "" || providedAPIKey != t.apiKey {
					log.Error().
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Str("key", providedAPIKey).
						Interface("status", http.StatusUnauthorized).
						Msg("Unauthorized request")
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				log.Info().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Interface("status", http.StatusOK).
					Msg("Authorized request")
			}

			// Serve the MCP requests
			streamableServer.ServeHTTP(w, r)
		}),
	}

	log.Info().Msgf("Starting MCP server with Streamable HTTP transport on :%s", t.port)
	log.Info().Msgf("HTTP endpoint: http://localhost:%s/mcp (direct JSON-RPC calls)", t.port)
	log.Info().Msg("Test with: curl -X POST http://localhost:" + t.port + "/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'")

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
