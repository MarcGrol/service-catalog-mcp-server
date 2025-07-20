package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"gopkg.in/yaml.v3"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

// NewListInterfacesTool returns the MCP tool definition and its handler for listing interfaces.
func NewLGetSingleInterfaceTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_interface",
			mcp.WithDescription("Gives details about a single interface (=web-api)"),
			mcp.WithString("interface_id", mcp.Required(), mcp.Description("The ID of the interface to get details for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			interfaceID, err := request.RequireString("interface_id")
			if err != nil {
				return mcp.NewToolResultError("Missing interface_id"), nil
			}
			iface, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error getting interface", err), nil
			}
			if !exists {
				return mcp.NewToolResultError("Interface with given ID not found"), nil
			}

			readableResult, err := yaml.Marshal(iface)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error serializing interface to yaml", err), nil
			}
			result := fmt.Sprintf("Found interface %s\n\n: %s",
				iface.InterfaceID, readableResult)
			return mcp.NewToolResultText(result), nil
		},
	}
}
