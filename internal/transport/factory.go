package transport

import (
	"github.com/mark3labs/mcp-go/server"
)

type ServerTransport interface {
	Start() error
}

func NewServerTransport(s *server.MCPServer, useSSE, useStreamable bool, port, baseURL, apiKey string) ServerTransport {
	if useStreamable {
		return NewStreamableHTTPTransport(s, port, apiKey)
	} else if useSSE {
		return NewSSETransport(s, baseURL, port)
	} else {
		return NewStdioTransport(s)
	}
}
