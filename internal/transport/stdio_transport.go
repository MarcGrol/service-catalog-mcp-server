package transport

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type StdioTransport struct {
	mcpServer *server.MCPServer
}

func NewStdioTransport(s *server.MCPServer) *StdioTransport {
	return &StdioTransport{
		mcpServer: s,
	}
}

func (t *StdioTransport) Start() error {
	log.Info().Msgf("Starting MCP server with stdio transport...")
	return server.ServeStdio(t.mcpServer)
}
