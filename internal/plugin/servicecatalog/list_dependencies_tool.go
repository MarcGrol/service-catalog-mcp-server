package servicecatalog

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/core/resp"
)

func (h *mcpHandler) listDependenciesTool() server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool(
			"list_dependencies_of_module",
			mcp.WithDescription("List all gradle dependencies of a module"),
			mcp.WithString("module_id", mcp.Required(), mcp.Description("The ID of the module to list gradle dependencies for")),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithOutputSchema[resp.List](),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Info()
			// extract params
			moduleID, err := request.RequireString("module_id")
			if err != nil {
				return mcp.NewToolResultError(resp.InvalidInput(ctx, "Missing module_id",
					"module_id",
					"Use a valid module identifier")), nil
			}

			// call business logic
			moduleNames, exists, err := h.repo.GetGradleDependenciesOfModule(ctx, moduleID)
			if err != nil {
				return mcp.NewToolResultError(resp.InternalError(ctx,
					fmt.Sprintf("error listing gradle dependencies of module %s: %s", moduleID, err))), nil
			}
			if !exists {
				return mcp.NewToolResultError(
					resp.NotFound(ctx,
						fmt.Sprintf("Module with ID %s not found", moduleID),
						"module_id",
						h.idx.Search(ctx, moduleID, 10).Modules)), nil

			}

			return mcp.NewToolResultJSON[resp.List](resp.SliceToList(moduleNames))
		},
	}
}
