package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewListInterfacesTool returns the MCP tool definition and its handler for listing interfaces.
func NewLGetSingleModuleTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_module",
			mcp.WithDescription("Gives details about a single module in the catalog"),
			mcp.WithString("module_id", mcp.Required(), mcp.Description("The ID of the module to get details for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			moduleID, err := request.RequireString("module_id")
			if err != nil {
				return mcp.NewToolResultError("Missing module_id"), nil
			}
			module, exists, err := repo.GetModuleOnID(ctx, moduleID)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error getting module", err), nil
			}
			if !exists {
				return mcp.NewToolResultError("Module with given ID not found"), nil
			}

			readableResult, err := yaml.Marshal(module)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error serializing module to yaml", err), nil
			}
			result := fmt.Sprintf("Found interface %s\n\n: %s",
				module.ModuleID, readableResult)
			return mcp.NewToolResultText(result), nil
		},
	}
}
