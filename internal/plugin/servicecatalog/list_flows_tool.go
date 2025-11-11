package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListFlowsTool returns the MCP tool definition and its handler for listing flows.
func (h *mcpHandler) listFlowsTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_flows",
			mcp.WithDescription("Lists all critical flows in the catalog."),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithIdempotentHintAnnotation(true),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]string](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

			// call business logic
			flows, err := h.repo.ListFlows(ctx)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error listing flows: %s", err))), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, flows)), nil
		},
	}
}
