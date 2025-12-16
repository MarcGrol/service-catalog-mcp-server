package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListMDatabaseConsumersTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) listModuleConsumersTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_module_consumers",
			mcp.WithDescription("List all modules that consume a given gradle module"),
			mcp.WithString("module_id", mcp.Required(), mcp.Description("The ID of the gradle module to list consumers for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[resp.List](),
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
			moduleNames, exists, err := h.repo.ListConsumersOfGradleModule(ctx, moduleID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error getting consumers of module %s: %s", moduleID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Module with ID %s not found", moduleID),
						"module_id",
						h.idx.Search(ctx, moduleID, 10).Modules,
					)), nil
			}

			return mcp.NewToolResultJSON[resp.List](resp.SliceToList(moduleNames))
		},
	}
}
