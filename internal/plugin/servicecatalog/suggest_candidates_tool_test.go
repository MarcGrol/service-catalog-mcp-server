package servicecatalog

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestSuggestCandidatesSuccess(t *testing.T) {
	_, idx, ctx, cleanup := setup(t)
	defer cleanup()

	// when
	result, err := NewMCPHandler(nil, idx).suggestCandidatesTool().Handler(ctx,
		createRequest("suggest_candidates", map[string]interface{}{
			"keyword": "partner",
		}))

	// then
	assert.NoError(t, err)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, `{"Modules":["partner","partner-jobs"`)
	assert.Contains(t, textResult.Text, `"Teams":["partner-experience","partner-experience-fe"`)
}
