package servicecatalog

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
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
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, `{"moduleID":"paymentengine/acm/webapp/acm",`)
}

func TestGetInterfaceTool_Success2(t *testing.T) {
	store, idx, ctx, cleanup := setup(t)
	defer cleanup()

	// when
	result, err := NewMCPHandler(store, idx).getSingleInterfaceTool().Handler(ctx, createRequest("interface_id", map[string]interface{}{
		"interface_id": "com.adyen.services.configurationapi.MeService",
	}))

	// then
	assert.NoError(t, err)
	textResult := result.Content[0].(mcp.TextContent)
	t.Logf("result: %+v", textResult.Text)
	//assert.Contains(t, textResult.Text, `{"moduleID":"paymentengine/acm/webapp/acm",`)
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
