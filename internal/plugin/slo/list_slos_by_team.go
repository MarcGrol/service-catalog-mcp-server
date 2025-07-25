package slo

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

func (h *mcpHandler) listSLOByTeamTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_slos_by_team",
			mcp.WithDescription("List all SLO's owned by a team"),
			mcp.WithString("team_id", mcp.Required(), mcp.Description("The ID of the team to list modules for")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Info().Any("request", request).Msg("list_slos_by_team")
			// extract params
			teamID, err := request.RequireString("team_id")
			if err != nil {
				log.Error().Err(err).Str("team_id", teamID).Msg("error extracting team_id")
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing team_id",
					"team_id",
					"Use a valid team identifier")), nil
			}
			log.Info().Str("team_id", teamID).Msg("success extracting team_id")

			// call business logic
			slos, exists, err := h.repo.ListSLOsByTeam(ctx, teamID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error listing slos of team %s: %s", teamID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Team with ID %s not found", teamID),
						"team_id",
						h.idx.Search(ctx, teamID, 10).Teams)), nil

			}

			return mcp.NewToolResultText(resp.Success(ctx, slos)), nil
		},
	}
}
