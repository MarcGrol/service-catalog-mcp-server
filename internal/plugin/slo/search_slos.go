package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
)

func (h *mcpHandler) searchSLOs() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"search_slos",
			mcp.WithDescription("Search all SLO's based on application, webapp,service, component or methods"),
			mcp.WithString("category", mcp.Required(), mcp.Description("Category to search on: Must be one of 'team', 'application', 'webapp', 'service', 'component' or 'methods'")),
			mcp.WithString("keyword", mcp.Required(), mcp.Description("The keyword to list SLOs for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]repo.SLO](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			category, err := request.RequireString("category")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing category",
					"keyword",
					"Use a keyword")), nil
			}
			keyword, err := request.RequireString("keyword")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing keyword",
					"keyword",
					"Use a keyword")), nil
			}

			// call business logic
			slos, exists, err := h.repo.SearchSLOs(ctx, category, keyword)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error searching slos on keyword %s: %s", keyword, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("No SLO with keyword %s.%s not found", category, keyword),
						"keyword",
						h.idx.Search(ctx, keyword, 10).Applications)), nil

			}

			return mcp.NewToolResultJSON[[]repo.SLO](slos)
		},
	}
}
