package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
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
		h.suggestCandidatesTool(),
		h.listModulesTool(),
		h.listModulesByComplexityTool(),
		h.getSingleModuleTool(),
		h.listInterfacesTool(),
		h.listInterfacesByComplexityTool(),
		h.getSingleInterfaceTool(),
		h.listModulesOfTeamsTool(),
		h.listMDatabaseConsumersTool(),
		h.listInterfaceConsumersTool(),
		h.listFlowsTool(),
		h.listFlowParticipantsTool(),
		h.listKindsTool(),
		h.listModulesWithKindTool(),
	)

	s.AddResources(
		h.modulesResource(),
	)

	s.AddPrompts(
		h.serviceCatalogPrompt(),
	)
}
