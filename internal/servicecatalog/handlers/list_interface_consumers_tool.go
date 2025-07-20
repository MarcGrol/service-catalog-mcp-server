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
func NewListInterfaceConsumersTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_interface_consumers",
			mcp.WithDescription("List all modules that consume a given interface (=web-api)"),
			mcp.WithString("interface_id", mcp.Required(), mcp.Description("The ID of the interface (=web-api) to list modules for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			interfaceID, err := request.RequireString("interface_id")
			if err != nil {
				return mcp.NewToolResultError("Missing interface_id"), nil
			}
			moduleNames, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing modules of interface", err), nil
			}
			if !exists {
				return mcp.NewToolResultError("Interface with given ID not found"), nil
			}

			result := fmt.Sprintf("Found %d modules for interface %s:\n\n%s", len(moduleNames), interfaceID, strings.Join(moduleNames, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
