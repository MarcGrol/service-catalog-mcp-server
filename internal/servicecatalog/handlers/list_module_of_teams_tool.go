package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/resp"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/search_index"
)

// NewListModulesOfTeamsTool returns the MCP tool definition and its handler for listing interfaces.
func NewListModulesOfTeamsTool(repo catalogrepo.Cataloger, idx search_index.Index) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules_of_teams",
			mcp.WithDescription("List all modules owned by a team"),
			mcp.WithString("team_id", mcp.Required(), mcp.Description("The ID of the team to list modules for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			teamID, err := request.RequireString("team_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput("Missing team_id",
					"team_id",
					"Use a valid team identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := repo.ListModulesOfTeam(ctx, teamID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(
					fmt.Sprintf("error listing modules of team %s: %s", teamID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(
						fmt.Sprintf("Team with ID %s not found", teamID),
						"team_id",
						idx.Search(ctx, teamID).Teams,
					)), nil

			}

			return mcp.NewToolResultText(resp.Success(moduleNames)), nil
		},
	}
}
