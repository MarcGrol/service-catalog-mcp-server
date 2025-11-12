package servicecatalog

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/data"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
)

func setup(t *testing.T) (repo.Cataloger, search.Index, context.Context, func()) {
	ctx := context.Background()
	serviceCatalogDatabaseFilename, fileCleanup, err := data.UnpackServiceCatalogDatabase(ctx)
	assert.NoError(t, err)
	defer fileCleanup()

	repo := repo.New(serviceCatalogDatabaseFilename)
	err = repo.Open(ctx)
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

func expectError(t *testing.T, result *mcp.CallToolResult, errorText string) {
	assert.True(t, result.IsError)
	content, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, content.Text, errorText)
}
