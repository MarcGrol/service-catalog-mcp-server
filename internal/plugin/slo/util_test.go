package slo

import (
	"testing"

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
