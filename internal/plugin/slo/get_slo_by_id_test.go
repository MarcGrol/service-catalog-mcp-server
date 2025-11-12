package slo

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/search"
)

func TestGetLOByIDTool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockSLORepo(ctrl)
	idxMock := search.NewMockIndex(ctrl)

	h := NewMCPHandler(repoMock, idxMock)
	tool := h.getSLOByIDTool()
	ctx := context.Background()

	t.Run("Successful get", func(t *testing.T) {
		sloID := "test-app"
		expectedSLO := repo.SLO{UID: "slo1", Application: sloID}
		repoMock.EXPECT().GetSLOByID(ctx, sloID).Return(expectedSLO, true, nil).Times(1)

		req := createRequest("get_slo", map[string]interface{}{"slo_id": sloID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, `"uid":"slo1",`)
	})

	t.Run("Application not found", func(t *testing.T) {
		sloID := "nonexistent-slo"
		repoMock.EXPECT().GetSLOByID(ctx, sloID).Return(repo.SLO{}, false, nil).Times(1)
		idxMock.EXPECT().Search(ctx, sloID, gomock.Any()).Return(search.Result{SLOs: []string{"suggested-slo"}}).Times(1)

		req := createRequest("get_slo", map[string]interface{}{"slo_id": sloID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		textResult := result.Content[0].(mcp.TextContent)
		expectError(t, result, `"status": "not_found"`)
		assert.Contains(t, textResult.Text, "SLO with ID nonexistent-slo not found")
		assert.Contains(t, textResult.Text, "suggested-slo")
	})

	t.Run("Repository error", func(t *testing.T) {
		sloID := "error-slo"
		expectedErr := errors.New("database error")
		repoMock.EXPECT().GetSLOByID(ctx, sloID).Return(repo.SLO{}, false, expectedErr).Times(1)

		req := createRequest("get_slo", map[string]interface{}{"slo_id": sloID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "error"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "error getting slo error-slo: database error")
	})

	t.Run("Missing application_id", func(t *testing.T) {
		req := createRequest("get_slo", map[string]interface{}{})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "invalid_input"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "Missing slo_id")
	})
}
