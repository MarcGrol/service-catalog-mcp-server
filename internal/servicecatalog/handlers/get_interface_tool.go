package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/resp"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/search"
)

// NewLGetSingleInterfaceTool returns the MCP tool definition and its handler for listing interfaces.
func NewLGetSingleInterfaceTool(repo catalogrepo.Cataloger, idx search.Index) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"get_interface",
			mcp.WithDescription("Gives details about a single interface (=web-api)"),
			mcp.WithString("interface_id", mcp.Required(), mcp.Description("The ID of the interface to get details for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			interfaceID, err := request.RequireString("interface_id")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing interface_id",
						"interface_id",
						"Use a valid interface identifier")), nil
			}

			// call business logic
			iface, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error getting interface %s: %s", interfaceID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Interface with ID %s not found", interfaceID),
						"interface_id",
						idx.Search(ctx, interfaceID, 10).Interfaces,
					)), nil
			}

			return mcp.NewToolResultText(resp.Success(ctx, iface)), nil
		},
	}
}
