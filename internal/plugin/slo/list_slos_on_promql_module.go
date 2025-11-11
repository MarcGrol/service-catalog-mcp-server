package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
)

func (h *mcpHandler) listSLOsOnPromQLModule() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_slos_on_module",
			mcp.WithDescription("Search all SLO's based on their module"),
			mcp.WithString("module_id", mcp.Required(), mcp.Description("Name of the module to list SLOs for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]repo.SLO](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			moduleID, err := request.RequireString("module_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing module_id",
					"module_id",
					"Use a valid module_id")), nil
			}

			// call business logic
			slos, exists, err := h.repo.ListSLOsByPromQLModule(ctx, moduleID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error searching slos on module_id %s: %s", moduleID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("No SLOs with module_id %s found", moduleID),
						"module_id",
						h.idx.Search(ctx, moduleID, 10).Webapps)), nil

			}

			return mcp.NewToolResultJSON[[]repo.SLO](slos)
		},
	}
}
