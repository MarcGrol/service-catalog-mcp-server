package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/resp"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewListModulesTool returns the MCP tool definition and its handler for listing modules.
func NewListModulesTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules",
			mcp.WithDescription("Lists all modules in the catalog."),
			mcp.WithString("filter_keyword", mcp.Description("The keyword to filter modules by.")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			keyword := request.GetString("module_id", "")

			// call business logic
			modules, err := repo.ListModules(ctx, keyword)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(
						fmt.Sprintf("error listing modules with keyword %s: %s", keyword, err))), nil
			}

			results := []string{}
			for _, mod := range modules {
				results = append(results, fmt.Sprintf("%s: %s", mod.Name, mod.Description))
			}
			return mcp.NewToolResultText(resp.Success(results)), nil
		},
	}
}
