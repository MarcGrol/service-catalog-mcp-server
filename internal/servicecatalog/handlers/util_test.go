package handlers

import (
	"context"
	"testing"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/constants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (catalogrepo.Cataloger, search.Index, context.Context, func()) {
	ctx := context.Background()

	repo := catalogrepo.New(constants.DatabaseFilename)
	err := repo.Open(ctx)
	assert.NoError(t, err)
	cleanup := func() {
		repo.Close(ctx)
	}

	idx := search.NewSearchIndex(ctx, repo)

	return repo, idx, ctx, cleanup
}

func createRequest(name string, args map[string]interface{}) mcp.CallToolRequest {
	req := mcp.CallToolRequest{Params: mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	}}
	return req
}

func expectSuccess(t *testing.T, result *mcp.CallToolResult, successText string) {
	assert.False(t, result.IsError)
	content, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, content.Text, successText)
}

func expectError(t *testing.T, result *mcp.CallToolResult, errorText string) {
	assert.True(t, result.IsError)
	content, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, content.Text, errorText)
}
