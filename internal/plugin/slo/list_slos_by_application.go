package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

func (h *mcpHandler) listSLOByApplicationTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_slos_by_application",
			mcp.WithDescription("List all SLO's owned by an application"),
			mcp.WithString("application_id", mcp.Required(), mcp.Description("The ID of the application to list SLOs for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			applicationID, err := request.RequireString("application_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing application_id",
					"application_id",
					"Use a valid application identifier")), nil
			}

			// call business logic
			slos, exists, err := h.repo.ListSLOsByApplication(ctx, applicationID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error listing slos of application %s: %s", applicationID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Application with ID %s not found", applicationID),
						"application_id",
						h.idx.Search(ctx, applicationID, 10).Applications)), nil

			}

			return mcp.NewToolResultText(resp.Success(ctx, slos)), nil
		},
	}
}
