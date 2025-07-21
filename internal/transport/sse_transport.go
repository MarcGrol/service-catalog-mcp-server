package transport

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mark3labs/mcp-go/server"
)

type SSETransport struct {
	mcpServer *server.MCPServer
	baseURL   string
	port      string
}

func NewSSETransport(s *server.MCPServer, baseURL, port string) *SSETransport {
	return &SSETransport{
		mcpServer: s,
		baseURL:   baseURL,
		port:      port,
	}
}

func (t *SSETransport) Start() error {
	fullBaseURL := fmt.Sprintf("%s:%s", t.baseURL, t.port)

	sseServer := server.NewSSEServer(t.mcpServer,
		server.WithBaseURL(fullBaseURL),
	)

	log.Info().Msgf("Starting MCP server with SSE transport on %s", fullBaseURL)
	log.Info().Msgf("SSE endpoint: %s/sse", fullBaseURL)
	log.Info().Msgf("Message endpoint: %s/message", fullBaseURL)

	if err := sseServer.Start(":" + t.port); err != nil {
		return err
	}
	return nil
}
