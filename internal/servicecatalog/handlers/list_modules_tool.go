package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewListModulesTool returns the MCP tool definition and its handler for listing modules.
func NewListModulesTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules",
			mcp.WithDescription("Lists all modules in the catalog."),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			modules, err := repo.ListModules(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing modules", err), nil
			}
			results := []string{}
			for _, p := range modules {
				results = append(results, fmt.Sprintf("%s: %s", p.Name, p.Description))
			}
			result := fmt.Sprintf("Found modules:\n\n%s", strings.Join(results, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
