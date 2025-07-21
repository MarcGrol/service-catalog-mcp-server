package transport

import (
	"github.com/rs/zerolog/log"

	"github.com/mark3labs/mcp-go/server"
)

type StreamableHTTPTransport struct {
	mcpServer *server.MCPServer
	port      string
}

func NewStreamableHTTPTransport(s *server.MCPServer, port string) *StreamableHTTPTransport {
	return &StreamableHTTPTransport{
		mcpServer: s,
		port:      port,
	}
}

func (t *StreamableHTTPTransport) Start() error {
	streamableServer := server.NewStreamableHTTPServer(t.mcpServer,
		server.WithStateLess(true),
	)
	log.Info().Msgf("Starting MCP server with Streamable HTTP transport on :%s", t.port)
	log.Info().Msgf("HTTP endpoint: http://localhost:%s/mcp (direct JSON-RPC calls)", t.port)
	log.Info().Msg("Test with: curl -X POST http://localhost:" + t.port + "/mcp -H 'Content-Type: application/json' -d '{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"tools/list\"}'")

	if err := streamableServer.Start(":" + t.port); err != nil {
		return err
	}
	return nil
}
