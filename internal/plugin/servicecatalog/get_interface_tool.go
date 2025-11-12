package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
)

// NewLGetSingleInterfaceTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) getSingleInterfaceTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_interface",
			mcp.WithDescription("Gives details about a single interface (=web-api)"),
			mcp.WithString("interface_id", mcp.Required(), mcp.Description("The ID of the interface to get details for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[repo.Interface](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			interfaceID, err := request.RequireString("interface_id")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing interface_id",
						"interface_id",
						"Use a valid interface identifier")), nil
			}

			// call business logic
			iface, exists, err := h.repo.GetInterfaceOnID(ctx, interfaceID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error getting interface %s: %s", interfaceID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Interface with ID %s not found", interfaceID),
						"interface_id",
						h.idx.Search(ctx, interfaceID, 10).Interfaces,
					)), nil
			}
			log.Printf("inline: %+v", iface)

			return mcp.NewToolResultJSON[repo.Interface](iface)
		},
	}
}
