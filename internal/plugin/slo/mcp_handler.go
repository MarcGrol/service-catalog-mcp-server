package slo

import (
	"context"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

type mcpHandler struct {
	repo repo.SLORepo
	idx  slosearch.Index
}

// NewMCPHandler creates a new instance of MCPHandler.
func NewMCPHandler(repo repo.SLORepo, idx slosearch.Index) *mcpHandler {
	return &mcpHandler{
		repo: repo,
		idx:  idx,
	}
}

// RegisterAllHandlers registers all tools, resources, and prompts with the MCP server.
func (h *mcpHandler) RegisterAllHandlers(ctx context.Context, s *server.MCPServer) {
	s.AddTools(
		//h.listSLOTool(), // Response is too big
		h.listSLOByTeamTool(),
		h.listSLOByApplicationTool(),
		h.listSLOByTeamTool(),
		h.getSLOByIDTool(),
		h.suggestCandidatesTool(),
	)
	s.AddResources(
		h.sloResource(),
	)

	s.AddPrompts(
		h.sloPrompt(),
	)
}
