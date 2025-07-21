package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/handlers"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/search_index"
)

type ServiceCatalog struct {
	server      *server.MCPServer
	repo        catalogrepo.Cataloger
	searchIndex search_index.SearchIndex
}

func New(s *server.MCPServer, repo catalogrepo.Cataloger) *ServiceCatalog {
	return &ServiceCatalog{
		server: s,
		repo:   repo,
	}
}

func (p *ServiceCatalog) Initialize(ctx context.Context) error {
	p.register()

	p.searchIndex = search_index.NewSearchIndex(ctx, p.repo)

	return nil
}

func (p *ServiceCatalog) register() {
	p.server.AddTools(
		handlers.NewSuggestCandidatesTool(p.searchIndex),
		handlers.NewListModulesTool(p.repo),
		handlers.NewGetSingleModuleTool(p.repo, p.searchIndex),
		handlers.NewListInterfacesTool(p.repo),
		handlers.NewLGetSingleInterfaceTool(p.repo, p.searchIndex),
		handlers.NewListModulesOfTeamsTool(p.repo, p.searchIndex),
		handlers.NewListMDatabaseConsumersTool(p.repo, p.searchIndex),
		handlers.NewListInterfaceConsumersTool(p.repo, p.searchIndex),
	)

	p.server.AddResources(
		handlers.NewModulesResource(p.repo),
	)

	p.server.AddPrompts(
		handlers.NewServiceCatalogTeamPrompt(),
	)
}
