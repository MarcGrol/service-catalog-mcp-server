package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
)

// NewSearchContentTool returns the MCP tool definition and its handler for searching content.
func NewSearchContentTool(store mystore.Store[model.Project]) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"search_content",
			mcp.WithDescription("Search for content in projects and tasks"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
			mcp.WithString("type", mcp.Description("Content type to search: project, task, all")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, err := request.RequireString("query")
			if err != nil {
				return mcp.NewToolResultError("Missing search query"), nil
			}
			searchType := request.GetString("type", "all")
			results := []string{
				fmt.Sprintf("Found in project config: %s", strings.ToLower(query)),
				fmt.Sprintf("Found in task #123: %s related item", query),
				fmt.Sprintf("Found in documentation: %s reference", query),
			}
			result := fmt.Sprintf("Search Results for '%s' (type: %s):\n\n%s", query, searchType, strings.Join(results, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
