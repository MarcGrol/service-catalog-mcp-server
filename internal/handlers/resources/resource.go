package resources

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Resouce struct {
	Resource mcp.Resource
	Handler  func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
}

func (h Resouce) Attrs() (mcp.Resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)) {
	return h.Resource, h.Handler
}
