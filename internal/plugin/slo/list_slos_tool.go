package slo

import (
	"context"
	"fmt"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (h *mcpHandler) listSLOTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_slos",
			mcp.WithDescription("List all SLO's"),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

			// call business logic
			slos, err := h.repo.ListSLOs(ctx)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error listing slos: %s", err))), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, slos)), nil
		},
	}
}
