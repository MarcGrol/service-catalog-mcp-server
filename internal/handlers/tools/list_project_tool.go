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

// NewListProjectTool returns the MCP tool definition and its handler for listing projects.
func NewListProjectTool(store mystore.Store[model.Project]) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_projects",
			mcp.WithDescription("Lists all projects"),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := store.List(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing projects", err), nil
			}
			results := []string{}
			for _, p := range projects {
				results = append(results, fmt.Sprintf("%s: %s", p.Name, p.Description))
			}
			result := fmt.Sprintf("Currently available project:\n\n%s", strings.Join(results, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
