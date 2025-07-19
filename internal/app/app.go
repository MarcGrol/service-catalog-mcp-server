package app

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/config"
	"github.com/MarcGrol/learnmcp/internal/model"

	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/MarcGrol/learnmcp/project"
	"github.com/MarcGrol/learnmcp/transport"
)

type Application struct {
	config          config.Config
	mcpServer       *server.MCPServer
	projectService  *project.ProjectService
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
		"Marc's project MCP Server", // Server name
		"1.0.0",                     // Version
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging())

	projectStore, projectStoreCleanup, err := mystore.New[model.Project](ctx)
	if err != nil {
		return nil, err
	}

	a.projectService = project.New(a.mcpServer, projectStore)
	a.projectService.Initialize(ctx)

	a.serverTransport = transport.NewServerTransport(a.mcpServer, a.config.UseSSE, a.config.UseStreamable, a.config.Port, a.config.BaseURL)

	return projectStoreCleanup, nil
}

func (a *Application) Run() error {
	if err := a.serverTransport.Start(); err != nil {
		return err
	}
	return nil
}
