package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/resp"
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
			// call business logic
			interfaces, err := repo.ListInterfaces(ctx)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error listing interfaces with keyword: %s", err))), nil
			}

			results := []interfaceDescriptor{}
			for _, i := range interfaces {
				results = append(results, interfaceDescriptor{
					InterfaceID: i.InterfaceID,
					Description: i.Description,
					Kind:        i.Kind,
				})
			}
			return mcp.NewToolResultText(resp.Success(ctx, results)), nil
		},
	}
}

type interfaceDescriptor struct {
	InterfaceID string `json:"interface_id"`
	Description string `json:"description"`
	Kind        string `json:"kind"`
}
