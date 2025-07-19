package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Contract mcp.Tool
	Handler  func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func (h Tool) Attrs() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	return h.Contract, h.Handler
}
