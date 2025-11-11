package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListInterfacesTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) listInterfacesTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_interfaces",
			mcp.WithDescription("Lists all interfaces (=web-api's) in the catalog"),
			mcp.WithString("filter_keyword", mcp.Required(), mcp.Description("The keyword to filter interfaces by.")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]interfaceDescriptor](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			keyword, err := request.RequireString("filter_keyword")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing filter_keyword",
						"filter_keyword",
						"Use a non-empty string as keyword")), nil
			}

			// call business logic
			interfaces, err := h.repo.ListInterfaces(ctx, keyword)
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
	InterfaceID     string
	Description     string
	Kind            string
	ComplexityScore int `yaml:",omitempty"`
}
