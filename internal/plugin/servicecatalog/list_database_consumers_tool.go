package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListMDatabaseConsumersTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) listMDatabaseConsumersTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_database_consumers",
			mcp.WithDescription("List all modules that consume a given database"),
			mcp.WithString("database_id", mcp.Required(), mcp.Description("The ID of the database to list modules for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[resp.List](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			databaseID, err := request.RequireString("database_id")
			if err != nil {
				return mcp.NewToolResultError(
					resp.InvalidInput(ctx, "Missing database_id",
						"database_id",
						"Use a valid database identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.ListDatabaseConsumers(ctx, databaseID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error getting database %s: %s", databaseID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Module with ID %s not found", databaseID),
						"database_id",
						h.idx.Search(ctx, databaseID, 10).Databases,
					)), nil
			}

			return mcp.NewToolResultJSON[resp.List](resp.SliceToList(moduleNames))
		},
	}
}
