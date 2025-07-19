package transport

import (
	"github.com/mark3labs/mcp-go/server"
)

type ServerTransport interface {
	Start() error
}

func NewServerTransport(s *server.MCPServer, useSSE, useStreamable bool, port, baseURL string) ServerTransport {
	if useStreamable {
		return NewStreamableHTTPTransport(s, port)
	} else if useSSE {
		return NewSSETransport(s, baseURL, port)
	} else {
		return NewStdioTransport(s)
	}
}
