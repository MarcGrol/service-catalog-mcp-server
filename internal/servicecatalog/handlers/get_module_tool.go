package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/resp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

// NewGetSingleModuleTool returns the MCP tool definition and its handler for listing interfaces.
func NewGetSingleModuleTool(repo catalogrepo.Cataloger, idx search.Index) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_module",
			mcp.WithDescription("Gives details about a single module in the catalog"),
			mcp.WithString("module_id", mcp.Required(), mcp.Description("The ID of the module to get details for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			moduleID, err := request.RequireString("module_id")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing module_id",
						"module_id",
						"Use a valid module identifier")), nil
			}

			// call business logic
			module, exists, err := repo.GetModuleOnID(ctx, moduleID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error getting module %s: %s", moduleID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Module with ID %s not found", moduleID),
						"interface_id",
						idx.Search(ctx, moduleID, 10).Modules,
					)), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, module)), nil
		},
	}
}
