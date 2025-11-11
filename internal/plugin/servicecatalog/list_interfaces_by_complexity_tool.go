package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListInterfacesByComplexityTool returns the MCP tool definition and its handler for listing interfaces by complexity.
func (h *mcpHandler) listInterfacesByComplexityTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_interfaces_by_complexity",
			mcp.WithDescription("Lists all interfaces in the catalog ordered DESC on complexity limited up to limit_to interfaces."),
			mcp.WithNumber("limit_to", mcp.Description("Maximum number of interfaces to list.")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[InterfaceDescriptorList](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			limit := request.GetInt("limit_to", 20)

			// call business logic
			interfaces, err := h.repo.ListInterfacesByComplexity(ctx, limit)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx,
						fmt.Sprintf("error listing interfaces by complexity: %s", err))), nil
			}

			results := []InterfaceDescriptor{}
			for _, i := range interfaces {
				results = append(results, InterfaceDescriptor{
					InterfaceID:     i.InterfaceID,
					Description:     i.Description,
					Kind:            i.Kind,
					ComplexityScore: i.MethodCount,
				})
			}
			return mcp.NewToolResultJSON[InterfaceDescriptorList](InterfaceDescriptorList{
				Interfaces: results,
			})
		},
	}
}

type InterfaceDescriptorList struct {
	Interfaces []InterfaceDescriptor `json:"interfaces"`
}
