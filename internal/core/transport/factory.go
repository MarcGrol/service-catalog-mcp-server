package transport

import (
	"github.com/mark3labs/mcp-go/server"
)

// Transport defines the interface for a server transport.
type Transport interface {
	Start() error
}

// NewTransport creates a new server transport based on the provided options.
func NewTransport(s *server.MCPServer, useSSE, useStreamable bool, port, baseURL, apiKey string) Transport {
	if useStreamable {
		return NewStreamableHTTPTransport(s, port, apiKey)
	} else if useSSE {
		return NewSSETransport(s, baseURL, port)
	} else {
		return NewStdioTransport(s)
	}
}
