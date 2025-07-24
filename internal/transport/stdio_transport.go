package transport

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// StdioTransport implements the ServerTransport interface for stdio.
type StdioTransport struct {
	mcpServer *server.MCPServer
}

// NewStdioTransport creates a new StdioTransport instance.
func NewStdioTransport(s *server.MCPServer) *StdioTransport {
	return &StdioTransport{
		mcpServer: s,
	}
}

// Start starts the stdio transport server.
func (t *StdioTransport) Start() error {
	log.Info().Msgf("Starting MCP server with stdio transport...")
	return server.ServeStdio(t.mcpServer)
}
