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
	mcpHandler  *handlers.MCPHandler
}

// New creates a new MCPServiceCatalog instance.
func New(s *server.MCPServer, repo catalogrepo.Cataloger, searchIndex search.Index) *MCPServiceCatalog {
	return &MCPServiceCatalog{
		server:      s,
		repo:        repo,
		searchIndex: searchIndex,
		mcpHandler:  handlers.NewMCPHandler(repo, searchIndex),
	}
}

// RegisterHandlers registers the service catalog handlers with the MCP server.
func (p *MCPServiceCatalog) RegisterHandlers(ctx context.Context) {

	p.server.AddTools(
		p.mcpHandler.NewSuggestCandidatesTool(),
		p.mcpHandler.NewListModulesTool(),
		p.mcpHandler.NewListModulesByComplexityTool(),
		p.mcpHandler.NewGetSingleModuleTool(),
		p.mcpHandler.NewListInterfacesTool(),
		p.mcpHandler.NewListInterfacesByComplexityTool(),
		p.mcpHandler.NewLGetSingleInterfaceTool(),
		p.mcpHandler.NewListModulesOfTeamsTool(),
		p.mcpHandler.NewListMDatabaseConsumersTool(),
		p.mcpHandler.NewListInterfaceConsumersTool(),
		p.mcpHandler.NewListFlowsTool(),
		p.mcpHandler.NewListFlowParticipantsTool(),
		p.mcpHandler.NewListKindsTool(),
		p.mcpHandler.NewListModulesWithKindTool(),
	)

	p.server.AddResources(
		p.mcpHandler.NewModulesResource(),
	)

	p.server.AddPrompts(
		p.mcpHandler.NewServiceCatalogPrompt(),
	)
}
