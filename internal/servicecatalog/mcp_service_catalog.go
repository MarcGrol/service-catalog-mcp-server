package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/handlers"
	search "github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

// MCPServiceCatalog represents the service catalog for MCP.
type MCPServiceCatalog struct {
	server      *server.MCPServer
	repo        catalogrepo.Cataloger
	searchIndex search.Index
}

// New creates a new MCPServiceCatalog instance.
func New(s *server.MCPServer, repo catalogrepo.Cataloger, searchIndex search.Index) *MCPServiceCatalog {
	return &MCPServiceCatalog{
		server:      s,
		repo:        repo,
		searchIndex: searchIndex,
	}
}

// RegisterHandlers registers the service catalog handlers with the MCP server.
func (p *MCPServiceCatalog) RegisterHandlers(ctx context.Context) {

	p.server.AddTools(
		handlers.NewSuggestCandidatesTool(p.searchIndex),
		handlers.NewListModulesTool(p.repo),
		handlers.NewListModulesByComplexityTool(p.repo),
		handlers.NewGetSingleModuleTool(p.repo, p.searchIndex),
		handlers.NewListInterfacesTool(p.repo),
		handlers.NewListInterfacesByComplexityTool(p.repo),
		handlers.NewLGetSingleInterfaceTool(p.repo, p.searchIndex),
		handlers.NewListModulesOfTeamsTool(p.repo, p.searchIndex),
		handlers.NewListMDatabaseConsumersTool(p.repo, p.searchIndex),
		handlers.NewListInterfaceConsumersTool(p.repo, p.searchIndex),
		handlers.NewListFlowsTool(p.repo),
		handlers.NewListFlowParticipantsTool(p.repo, p.searchIndex),
		handlers.NewListKindsTool(p.repo),
		handlers.NewListModulesWithKindTool(p.repo, p.searchIndex),
	)

	p.server.AddResources(
		handlers.NewModulesResource(p.repo),
	)

	p.server.AddPrompts(
		handlers.NewServiceCatalogPrompt(),
	)
}
