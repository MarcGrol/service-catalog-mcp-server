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
func NewListMDatabaseConsumersTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_database_consumers",
			mcp.WithDescription("List all modules that consume a given database"),
			mcp.WithString("database_id", mcp.Required(), mcp.Description("The ID of the database to list modules for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			databaseID, err := request.RequireString("database_id")
			if err != nil {
				return mcp.NewToolResultError("Missing database_id"), nil
			}
			moduleNames, exists, err := repo.ListDatabaseConsumers(ctx, databaseID)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing modules of database", err), nil
			}
			if !exists {
				return mcp.NewToolResultError("Database with given ID not found"), nil
			}

			result := fmt.Sprintf("Found %d modules for database %s:\n\n%s", len(moduleNames), databaseID, strings.Join(moduleNames, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
