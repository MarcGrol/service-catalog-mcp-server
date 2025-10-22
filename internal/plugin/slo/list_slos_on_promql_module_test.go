package slo

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

func TestListSLOsOnPromQLModuleTool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockSLORepo(ctrl)
	idxMock := slosearch.NewMockIndex(ctrl)

	h := NewMCPHandler(repoMock, idxMock)
	tool := h.listSLOsOnPromQLModule()
	ctx := context.Background()

	t.Run("Successful list", func(t *testing.T) {
		moduleID := "test-app"
		expectedSLOs := []repo.SLO{
			{UID: "slo1", Application: moduleID},
		}
		repoMock.EXPECT().ListSLOsByPromQLModule(ctx, moduleID).Return(expectedSLOs, true, nil).Times(1)

		req := createRequest("list_slos_on_module", map[string]interface{}{"module_id": moduleID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		expectSuccess(t, result, `"status": "success"`)
		textResult := result.Content[0].(mcp.TextContent).Text
		assert.Contains(t, textResult, `"uid": "slo1",`)
	})

	t.Run("SLOs for module not found", func(t *testing.T) {
		moduleID := "nonexistent-module"
		repoMock.EXPECT().ListSLOsByPromQLModule(ctx, moduleID).Return([]repo.SLO{}, false, nil).Times(1)
		idxMock.EXPECT().Search(ctx, moduleID, gomock.Any()).Return(slosearch.Result{Webapps: []string{"suggested-module"}}).Times(1)

		req := createRequest("list_slos_on_module", map[string]interface{}{"module_id": moduleID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "not_found"`)
		textResult := result.Content[0].(mcp.TextContent).Text
		assert.Contains(t, textResult, "No SLOs with module_id nonexistent-module found")
		assert.Contains(t, textResult, "suggested-module")
	})

	t.Run("Repository error", func(t *testing.T) {
		moduleID := "error-module"
		expectedErr := fmt.Errorf("database error")
		repoMock.EXPECT().ListSLOsByPromQLModule(ctx, moduleID).Return(nil, false, expectedErr).Times(1)

		req := createRequest("list_slos_on_module", map[string]interface{}{"module_id": moduleID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "error"`)
		textResult := result.Content[0].(mcp.TextContent).Text
		assert.Contains(t, textResult, "error searching slos on module_id error-module: database error")
	})

	t.Run("Missing module_id", func(t *testing.T) {
		req := createRequest("list_slos_on_module", map[string]interface{}{})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "invalid_input"`)
		textResult := result.Content[0].(mcp.TextContent).Text
		assert.Contains(t, textResult, "Missing module_id")
	})
}
