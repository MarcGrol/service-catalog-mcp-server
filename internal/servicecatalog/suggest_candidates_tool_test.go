package servicecatalog

import (
	"testing"

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
	expectSuccess(t, result, `"status": "success"`)
}
