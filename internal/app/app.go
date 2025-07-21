package app

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/config"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/search"
	"github.com/MarcGrol/learnmcp/internal/transport"
)

type Application struct {
	config          config.Config
	mcpServer       *server.MCPServer
	serviceCatalog  *servicecatalog.ServiceCatalog
	serverTransport transport.ServerTransport
}

func New(cfg config.Config) *Application {
	return &Application{
		config: cfg,
	}
}

func (a *Application) Initialize(ctx context.Context) (func(), error) {
	// Create a new MCP server
	a.mcpServer = server.NewMCPServer(
		"Marc's MCP Server", // Server name
		"1.0.0",             // Version
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging())

	{
		catalogRepo := catalogrepo.New(a.config.DatabaseFile)
		err := catalogRepo.Open(ctx)
		if err != nil {
			return nil, fmt.Errorf("error opening database: %s", err)
		}

		searchIndex := search.NewSearchIndex(ctx, catalogRepo)

		a.serviceCatalog = servicecatalog.New(a.mcpServer, catalogRepo, searchIndex)
		a.serviceCatalog.RegisterHandlers(ctx)
	}

	a.serverTransport = transport.NewServerTransport(a.mcpServer, a.config.UseSSE, a.config.UseStreamable, a.config.Port, a.config.BaseURL)

	return func() {}, nil
}

func (a *Application) Run() error {
	if err := a.serverTransport.Start(); err != nil {
		return err
	}
	return nil
}
