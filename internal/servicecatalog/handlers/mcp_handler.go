package handlers

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

// MCPHandler handles the MCP (Machine Configuration Pool) related operations.
// It uses a catalog repository to manage the catalog and a search index to search for items.
// The catalog repository is used to add, update, and delete items in the catalog.
// The search index is used to search for items in the catalog.
// The MCPHandler is responsible for handling the MCP related operations.
type MCPHandler struct {
	repo catalogrepo.Cataloger
	idx  search.Index
}

// NewMCPHandler creates a new instance of MCPHandler.
func NewMCPHandler(repo catalogrepo.Cataloger, idx search.Index) *MCPHandler {
	return &MCPHandler{
		repo: repo,
		idx:  idx,
	}
}

// RegisterAllHandlers registers all tools, resources, and prompts with the MCP server.
func (h *MCPHandler) RegisterAllHandlers(s *server.MCPServer, ctx context.Context) {
	s.AddTools(
		h.NewSuggestCandidatesTool(),
		h.NewListModulesTool(),
		h.NewListModulesByComplexityTool(),
		h.NewGetSingleModuleTool(),
		h.NewListInterfacesTool(),
		h.NewListInterfacesByComplexityTool(),
		h.NewLGetSingleInterfaceTool(),
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
