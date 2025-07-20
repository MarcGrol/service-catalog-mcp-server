package servicecatalog

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/handlers"
)

type ServiceCatalog struct {
	server *server.MCPServer
	repo   catalogrepo.Cataloger
}

func New(s *server.MCPServer, repo catalogrepo.Cataloger) *ServiceCatalog {
	return &ServiceCatalog{
		server: s,
		repo:   repo,
	}
}

func (p *ServiceCatalog) Initialize(ctx context.Context) error {

	p.register()

	return nil
}

func (p *ServiceCatalog) register() {
	p.server.AddTools(
		handlers.NewListModulesTool(p.repo),
		handlers.NewLGetSingleModuleTool(p.repo),
		handlers.NewListInterfacesTool(p.repo),
		handlers.NewLGetSingleInterfaceTool(p.repo),
	)

	p.server.AddResources(
		handlers.NewModulesResource(p.repo),
	)

	p.server.AddPrompts(
		handlers.NewServiceCatalogTeamPrompt(),
	)
}
