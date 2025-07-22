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
			mcp.WithString("filter_keyword", mcp.Required(), mcp.Description("The keyword to filter interfaces by.")),
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
			interfaces, err := repo.ListInterfaces(ctx, keyword)
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
	InterfaceID     string `json:"interface_id"`
	Description     string `json:"description"`
	Kind            string `json:"kind"`
	ComplexityScore int    `json:"complexityScore"`
}

/*
complexity_score = (
  (line_count / max_line_count) * 0.25 +
  (database_count / max_database_count) * 0.20 +
  (team_count / max_team_count) * 0.15 +
  (interface_count / max_interface_count) * 0.15 +
  (job_count / max_job_count) * 0.10 +
  (file_count / max_file_count) * 0.10 +
  (flow_count / max_flow_count) * 0.05
) * 100
*/
