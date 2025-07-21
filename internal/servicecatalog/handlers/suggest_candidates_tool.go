package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/resp"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/search"
)

// NewSuggestCandidatesTool returns the MCP tool definition and its handler for listing interfaces.
func NewSuggestCandidatesTool(index search.Index) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"suggest_candidates",
			mcp.WithDescription("Suggest matching modules, interfaces, databases, or teams based on user input."),
			mcp.WithString("keyword", mcp.Description("The keyword to search modules, interfaces, databases, or teams for.")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			keyword, err := request.RequireString("keyword")
			if err != nil {
				resp.InvalidInput(ctx, "Missing keyword",
					"keyword",
					"Use a valid keyword")
			}

			// call business logic
			searchResult := index.Search(ctx, keyword)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error searching for candidates like %s: %s", keyword, err))), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, searchResult)), nil
		},
	}
}
