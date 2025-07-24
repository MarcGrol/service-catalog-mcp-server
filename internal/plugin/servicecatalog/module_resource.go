package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewModulesResource returns the MCP resource contract and handler for modules configuration.
func (h *mcpHandler) modulesResource() server.ServerResource {
	return server.ServerResource{
		Resource: mcp.NewResource(
			"catalog://modules",
			"List of modules in the catalog",
			mcp.WithMIMEType("application/json"),
		),
		Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			// call business logic
			modules, err := h.repo.ListModules(ctx, "")
			if err != nil {
				return nil, fmt.Errorf("error listing modules: %s", err)
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "application/json",
					Text:     resp.Success(ctx, modules),
				},
			}, nil
		},
	}
}
