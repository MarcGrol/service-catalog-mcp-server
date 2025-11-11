package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

// NewSuggestCandidatesTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) suggestCandidatesTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"suggest_slos",
			mcp.WithDescription("Suggest matching slos, applications or teams based on user input."),
			mcp.WithString("keyword", mcp.Required(), mcp.Description("The keyword to search modules, interfaces, databases, or teams for.")),
			mcp.WithNumber("limit_to", mcp.Description("Maximum number of results per category to return.")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[slosearch.Result](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyword, err := request.RequireString("keyword")
			if err != nil {
				resp.InvalidInput(ctx, "Missing keyword",
					"keyword",
					"Use a valid keyword")
			}
			limit := request.GetInt("limit_to", 10)

			// call business logic
			searchResult := h.idx.Search(ctx, keyword, limit)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error searching for slos like %s: %s", keyword, err))), nil
			}

			return mcp.NewToolResultJSON[slosearch.Result](searchResult)
		},
	}
}
