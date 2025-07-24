package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

type mcpHandler struct {
	repo catalogrepo.Cataloger
	idx  search.Index
}

// NewMCPHandler creates a new instance of MCPHandler.
func NewMCPHandler(repo catalogrepo.Cataloger, idx search.Index) *mcpHandler {
	return &mcpHandler{
		repo: repo,
		idx:  idx,
	}
}

// RegisterAllHandlers registers all tools, resources, and prompts with the MCP server.
func (h *mcpHandler) RegisterAllHandlers(ctx context.Context, s *server.MCPServer) {
	s.AddTools(
		h.NewSuggestCandidatesTool(),
		h.NewListModulesTool(),
		h.NewListModulesByComplexityTool(),
		h.NewGetSingleModuleTool(),
		h.NewListInterfacesTool(),
		h.NewListInterfacesByComplexityTool(),
		h.NewGetSingleInterfaceTool(),
		h.NewListModulesOfTeamsTool(),
		h.NewListMDatabaseConsumersTool(),
		h.NewListInterfaceConsumersTool(),
		h.NewListFlowsTool(),
		h.NewListFlowParticipantsTool(),
		h.NewListKindsTool(),
		h.NewListModulesWithKindTool(),
	)

	s.AddResources(
		h.NewModulesResource(),
	)

	s.AddPrompts(
		h.NewServiceCatalogPrompt(),
	)
}
