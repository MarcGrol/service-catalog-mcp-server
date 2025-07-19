package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Tool struct {
	Tool    mcp.Tool
	Handler func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func (h Tool) Funcs() (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	return h.Tool, h.Handler
}
