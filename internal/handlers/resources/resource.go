package resources

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Resource struct {
	Contract mcp.Resource
	Handler  func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
}

func (h Resource) Attrs() (mcp.Resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)) {
	return h.Contract, h.Handler
}
