package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewListProjectToolAndHandler returns the MCP tool definition and its handler for listing projects.
func NewListProjectToolAndHandler(store mystore.Store[model.Project]) (mcp.Tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)) {
	tool := mcp.NewTool(
		"list_projects",
		mcp.WithDescription("Lists all projects"),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
	}
	return tool, handler
}
