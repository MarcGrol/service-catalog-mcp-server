package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListModulesWithKindTool returns the MCP tool definition and its handler for listing modules with kind.
func (h *mcpHandler) listModulesWithKindTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules_with_kind",
			mcp.WithDescription("List all modules that are of this kind"),
			mcp.WithString("kind_id", mcp.Required(), mcp.Description("The ID of the kind")),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]string](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			kindID, err := request.RequireString("kind_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing kind_id",
					"kind_id",
					"Use a valid kind identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.ListModulesWithKind(ctx, kindID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx, // Corrected error message
						fmt.Sprintf("error listing modules with kind %s: %s", kindID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("No modules found for kind with ID %s not found", kindID),
						"kind_id", // Corrected parameter name
						h.idx.Search(ctx, kindID, 10).Kinds,
					)), nil
			}

			// return result
			return mcp.NewToolResultText(resp.Success(ctx, moduleNames)), nil
		},
	}
}
