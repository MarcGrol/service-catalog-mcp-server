package transport

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

type StdioTransport struct {
	mcpServer *server.MCPServer
}

func NewStdioTransport(s *server.MCPServer) *StdioTransport {
	return &StdioTransport{
		mcpServer: s,
	}
}

func (t *StdioTransport) Start(addr string) error {
	log.Println("Starting MCP server with stdio transport...")
	return server.ServeStdio(t.mcpServer)
}
