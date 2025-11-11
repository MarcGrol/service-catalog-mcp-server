package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListModulesOfTeamsTool returns the MCP tool definition and its handler for listing interfaces.
func (h *mcpHandler) listModulesOfTeamsTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_modules_of_teams",
			mcp.WithDescription("List all modules owned by a team"),
			mcp.WithString("team_id", mcp.Required(), mcp.Description("The ID of the team to list modules for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[resp.List](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Info()
			// extract params
			teamID, err := request.RequireString("team_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing team_id",
					"team_id",
					"Use a valid team identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.ListModulesOfTeam(ctx, teamID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error listing modules of team %s: %s", teamID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Team with ID %s not found", teamID),
						"team_id",
						h.idx.Search(ctx, teamID, 10).Teams)), nil

			}

			return mcp.NewToolResultJSON[resp.List](resp.SliceToList(moduleNames))
		},
	}
}
