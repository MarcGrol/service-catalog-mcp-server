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
	p.mcpHandler.RegisterAllHandlers(p.server, ctx)
}
