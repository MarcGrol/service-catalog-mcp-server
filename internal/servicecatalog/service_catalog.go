package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/handlers"
	search "github.com/MarcGrol/learnmcp/internal/servicecatalog/search"
)

type ServiceCatalog struct {
	server      *server.MCPServer
	repo        catalogrepo.Cataloger
	searchIndex search.Index
}

func New(s *server.MCPServer, repo catalogrepo.Cataloger, searchIndex search.Index) *ServiceCatalog {
	return &ServiceCatalog{
		server:      s,
		repo:        repo,
		searchIndex: searchIndex,
	}
}

func (p *ServiceCatalog) RegisterHandlers(ctx context.Context) {

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
	)

	p.server.AddResources(
		handlers.NewModulesResource(p.repo),
	)

	p.server.AddPrompts(
		handlers.NewServiceCatalogPrompt(),
	)
}
