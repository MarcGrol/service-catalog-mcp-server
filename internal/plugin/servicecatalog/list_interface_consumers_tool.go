package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListInterfaceConsumersTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) listInterfaceConsumersTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_interface_consumers",
			mcp.WithDescription("List all modules that consume a given interface (=web-api)"),
			mcp.WithString("interface_id", mcp.Required(), mcp.Description("The ID of the interface (=web-api) to list modules for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[resp.List](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			interfaceID, err := request.RequireString("interface_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing interface_id",
					"interface_id",
					"Use a valid interface identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.ListInterfaceConsumers(ctx, interfaceID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error listing consumers of interface %s: %s", interfaceID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Interface with ID %s not found", interfaceID),
						"interface_id",
						h.idx.Search(ctx, interfaceID, 10).Interfaces,
					)), nil
			}

			// return result
			return mcp.NewToolResultJSON[resp.List](resp.SliceToList(moduleNames))
		},
	}
}
