package slo

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/search"
)

type mcpHandler struct {
	repo repo.SLORepo
	idx  search.Index
}

// NewMCPHandler creates a new instance of mcpHandler.
func NewMCPHandler(repo repo.SLORepo, idx search.Index) *mcpHandler {
	return &mcpHandler{
		repo: repo,
		idx:  idx,
	}
}

// RegisterAllHandlers registers all tools, resources, and prompts with the MCP server.
func (h *mcpHandler) RegisterAllHandlers(ctx context.Context, s *server.MCPServer) {
	s.AddTools(
		h.suggestCandidatesTool(),
		h.searchSLOs(),
		h.listSLOsOnPromQLWebService(),
		h.listSLOsOnPromQLModule(),
		h.getSLOByIDTool(),
	)
	s.AddResources(
		h.sloResource(),
	)

	s.AddPrompts(
		h.sloPrompt(),
	)
}
