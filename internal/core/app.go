package core

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/config"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/transport"
)

// MCPService represents the interface for the MCP service.
type MCPService interface {
	RegisterAllHandlers(ctx context.Context, s *server.MCPServer)
}

// Application represents the main application structure.
type Application struct {
	config          config.Config
	mcpServer       *server.MCPServer
	mcpServices     []MCPService
	serverTransport transport.Transport
}

// New creates a new Application instance.
func New(cfg config.Config, mcpServices ...MCPService) *Application {
	return &Application{
		config:      cfg,
		mcpServices: mcpServices,
	}
}

// Initialize initializes the application.
func (a *Application) Initialize(ctx context.Context) (func(), error) {
	// Create a new MCP server
	a.mcpServer = server.NewMCPServer(
		"Marc's MCP Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(loggingHooks()))

	for _, service := range a.mcpServices {
		service.RegisterAllHandlers(ctx, a.mcpServer)
	}

	a.serverTransport = transport.NewTransport(a.mcpServer, a.config.UseSSE, a.config.UseStreamable, a.config.Port, a.config.BaseURL, a.config.APIKey)

	return func() {}, nil
}

// Run starts the application's server transport.
func (a *Application) Run() error {
	if err := a.serverTransport.Start(); err != nil {
		return err
	}
	return nil
}

func loggingHooks() *server.Hooks {
	return &server.Hooks{
		OnBeforeCallTool: []server.OnBeforeCallToolFunc{
			func(ctx context.Context, id any, req *mcp.CallToolRequest) {
				log.Info().
					Str("method", "tool").
					Any("request_id", id).
					Str("name", req.Params.Name).
					Any("args", req.Params.Arguments).
					Send()
			},
		},
		OnAfterCallTool: []server.OnAfterCallToolFunc{
			func(ctx context.Context, id any, req *mcp.CallToolRequest, resp *mcp.CallToolResult) {
				log.Info().
					Str("method", "tool").
					Any("request_id", id).
					Str("name", req.Params.Name).
					Any("args", req.Params.Arguments).
					Bool("success", !resp.IsError).Send()
			},
		},
		OnBeforeReadResource: []server.OnBeforeReadResourceFunc{
			func(ctx context.Context, id any, req *mcp.ReadResourceRequest) {
				log.Info().Str("method", "resource").
					Any("request_id", id).
					Str("resource_method", req.Request.Method).
					Any("args", req.Params.Arguments).
					Send()
			},
		},
		OnAfterReadResource: []server.OnAfterReadResourceFunc{
			func(ctx context.Context, id any, req *mcp.ReadResourceRequest, resp *mcp.ReadResourceResult) {
				log.Info().
					Str("method", "resource").
					Any("request_id", id).
					Str("resource_method", req.Request.Method).
					Any("args", req.Params.Arguments).
					Send()
			},
		},
	}
}
