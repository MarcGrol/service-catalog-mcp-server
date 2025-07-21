package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewModulesResource returns the MCP resource contract and handler for modules configuration.
func NewModulesResource(repo catalogrepo.Cataloger) server.ServerResource {
	return server.ServerResource{
		Resource: mcp.NewResource(
			"catalog://modules",
			"List of modules in the catalog",
			mcp.WithMIMEType("application/json"),
		),
		Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			// call business logic
			modules, err := repo.ListModules(ctx, "")
			if err != nil {
				return nil, fmt.Errorf("error listing modules: %s", err)
			}

			// return result
			modulesJson, err := json.MarshalIndent(
				map[string]interface{}{
					"total_modules": len(modules),
					"modules":       modules,
				}, "", "  ")
			if err != nil {
				return nil, err
			}

			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "application/json",
					Text:     string(modulesJson),
				},
			}, nil
		},
	}
}
