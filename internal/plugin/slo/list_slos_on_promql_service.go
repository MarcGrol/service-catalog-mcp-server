package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

func (h *mcpHandler) listSLOsOnPromQLWebService() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_slos_on_service",
			mcp.WithDescription("Search all SLO's based on a web service"),
			mcp.WithString("service-name", mcp.Required(), mcp.Description("Name of the web-service to list SLOs for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			serviceName, err := request.RequireString("service-name")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing service-name",
					"service-name",
					"Use a valid service-name")), nil
			}

			// call business logic
			slos, exists, err := h.repo.ListSLOsByPromQLService(ctx, serviceName)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error searching slos on service-name %s: %s", serviceName, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("No SLO with service-name %s not found", serviceName),
						"service-name",
						h.idx.Search(ctx, serviceName, 10).Applications)), nil

			}

			return mcp.NewToolResultText(resp.Success(ctx, slos)), nil
		},
	}
}
