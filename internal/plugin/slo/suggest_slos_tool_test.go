package slo

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/sloconstants"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

func TestSuggestCandidatesSuccess(t *testing.T) {
	ctx := context.Background()
	repo := repo.New(sloconstants.SLODatabaseFilename())
	repo.Open(ctx)
	idx := slosearch.NewSearchIndex(ctx, repo)

	// when
	result, err := NewMCPHandler(nil, idx).suggestCandidatesTool().Handler(ctx,
		createRequest("suggest_slos", map[string]interface{}{
			"keyword": "partner",
		}))

	// then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, `partner_referral_general_latency`)
}
