package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewListInterfacesTool returns the MCP tool definition and its handler for listing interfaces.
func NewListInterfacesTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_interfaces",
			mcp.WithDescription("Lists all interfaces (=web-api's) in the catalog"),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			interfaces, err := repo.ListInterfaces(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing modules", err), nil
			}
			results := []string{}
			for _, p := range interfaces {
				results = append(results, fmt.Sprintf("%s: %s", p.InterfaceID, p.Description))
			}
			result := fmt.Sprintf("Found %d interfaces:\n\n%s", len(results), strings.Join(results, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
