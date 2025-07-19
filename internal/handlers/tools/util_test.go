package tools

import (
	"context"
	"testing"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

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

func expectSuccess(t *testing.T, result *mcp.CallToolResult, successText string) {
	assert.False(t, result.IsError)
	content, ok := result.Content[0].(mcp.TextContent)
	assert.True(t, ok)
	assert.Contains(t, content.Text, successText)
}

func expectEmpty(t *testing.T, ctx context.Context, store mystore.Store[model.Project]) {
	stored, err := store.List(ctx)
	assert.NoError(t, err)
	assert.Empty(t, stored)
}

func expectProject(t *testing.T, ctx context.Context, store mystore.Store[model.Project], expected model.Project) {
	stored, err := store.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, stored, 1)
	assert.Equal(t, expected.Name, stored[0].Name)
	assert.Equal(t, expected.Description, stored[0].Description)
	assert.Equal(t, expected.Authors, stored[0].Authors)
}
