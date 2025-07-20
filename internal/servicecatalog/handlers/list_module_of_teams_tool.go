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
func NewListModulesOfTeamsTool(repo catalogrepo.Cataloger) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules_of_teams",
			mcp.WithDescription("List all modules owned by a team"),
			mcp.WithString("team_id", mcp.Required(), mcp.Description("The ID of the team to list modules for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			teamID, err := request.RequireString("team_id")
			if err != nil {
				return mcp.NewToolResultError("Missing team_id"), nil
			}
			moduleNames, exists, err := repo.ListModulesOfTeam(ctx, teamID)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing modules of team", err), nil
			}
			if !exists {
				return mcp.NewToolResultError("Team with given ID not found"), nil
			}

			result := fmt.Sprintf("Found %d modules for team %s:\n\n%s", len(moduleNames), teamID, strings.Join(moduleNames, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
