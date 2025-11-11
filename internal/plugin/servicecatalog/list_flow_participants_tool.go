package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

// NewListFlowParticipantsTool returns the MCP tool definition and its handler for listing flow participants.
func (h *mcpHandler) listFlowParticipantsTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_flow_participants",
			mcp.WithDescription("List all modules that that are participants of this flow"),
			mcp.WithString("flow_id", mcp.Required(), mcp.Description("The ID of the flow")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[[]string](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// extract params
			flowID, err := request.RequireString("flow_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing flow_id",
					"flow_id",
					"Use a valid flow identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.ListParticpantsOfFlow(ctx, flowID)
			if err != nil {
				return mcp.NewToolResultError(
					resp.InternalError(ctx, // Corrected error message
						fmt.Sprintf("error listing participants of flow %s: %s", flowID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Flow with ID %s not found", flowID),
						"flow_id", // Corrected parameter name
						h.idx.Search(ctx, flowID, 10).Flows,
					)), nil
			}

			// return result
			return mcp.NewToolResultText(resp.Success(ctx, moduleNames)), nil
		},
	}
}
