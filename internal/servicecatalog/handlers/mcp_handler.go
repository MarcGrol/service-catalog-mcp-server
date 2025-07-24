package handlers

import (
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

// MCPHandler handles the MCP (Machine Configuration Pool) related operations.
// It uses a catalog repository to manage the catalog and a search index to search for items.
// The catalog repository is used to add, update, and delete items in the catalog.
// The search index is used to search for items in the catalog.
// The MCPHandler is responsible for handling the MCP related operations.
type MCPHandler struct {
	repo catalogrepo.Cataloger
	idx  search.Index
}

// NewMCPHandler creates a new instance of MCPHandler.
func NewMCPHandler(repo catalogrepo.Cataloger, idx search.Index) *MCPHandler {
	return &MCPHandler{
		repo: repo,
		idx:  idx,
	}
}
