package slo

import (
	"context"
	"fmt"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (h *mcpHandler) getSLOByIDTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_slo",
			mcp.WithDescription("Gives details about a single slo"),
			mcp.WithString("slo_id", mcp.Required(), mcp.Description("The ID of the slo to get details for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			sloID, err := request.RequireString("slo_id")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing slo_id",
						"slo_id",
						"Use a valid slo identifier")), nil
			}

			// call business logic
			module, exists, err := h.repo.GetSLOByID(ctx, sloID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error getting slo %s: %s", sloID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("SLO with ID %s not found", sloID),
						"slo_id",
						h.idx.Search(ctx, sloID, 10).SLOs,
					)), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, module)), nil
		},
	}
}
