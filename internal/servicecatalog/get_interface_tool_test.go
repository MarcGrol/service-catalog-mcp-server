package servicecatalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInterfaceTool_Success(t *testing.T) {
	store, idx, ctx, cleanup := setup(t)
	defer cleanup()

	// when
	result, err := NewMCPHandler(store, idx).getSingleInterfaceTool().Handler(ctx, createRequest("interface_id", map[string]interface{}{
		"interface_id": "com.adyen.services.acm.AcmService",
	}))

	// then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	t.Logf("result: %+v", result)
}

func TestGetInterfaceTool_NotFound(t *testing.T) {
	store, idx, ctx, cleanup := setup(t)
	defer cleanup()

	// when
	result, err := NewMCPHandler(store, idx).getSingleInterfaceTool().Handler(ctx, createRequest("interface_id", map[string]interface{}{
		"interface_id": "lalala",
	}))

	// then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	t.Logf("result: %+v", result)

}
