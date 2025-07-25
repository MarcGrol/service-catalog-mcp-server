package slo

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/slosearch"
)

func TestListSLOByApplicationTool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockSLORepo(ctrl)
	idxMock := slosearch.NewMockIndex(ctrl)

	h := NewMCPHandler(repoMock, idxMock)
	tool := h.listSLOByApplicationTool()
	ctx := context.Background()

	t.Run("Successful list", func(t *testing.T) {
		applicationID := "test-app"
		expectedSLOs := []repo.SLO{
			{UID: "slo1", Application: applicationID},
			{UID: "slo2", Application: applicationID},
		}
		repoMock.EXPECT().ListSLOsByApplication(ctx, applicationID).Return(expectedSLOs, true, nil).Times(1)

		req := createRequest("list_slos_by_application", map[string]interface{}{"application_id": applicationID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)

		expectSuccess(t, result, `"status": "success"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, `"uid": "slo1",`)
	})

	t.Run("Application not found", func(t *testing.T) {
		applicationID := "nonexistent-app"
		repoMock.EXPECT().ListSLOsByApplication(ctx, applicationID).Return([]repo.SLO{}, false, nil).Times(1)
		idxMock.EXPECT().Search(ctx, applicationID, gomock.Any()).Return(slosearch.Result{Applications: []string{"suggested-app"}}).Times(1)

		req := createRequest("list_slos_by_application", map[string]interface{}{"application_id": applicationID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "Application with ID nonexistent-app not found")
		assert.Contains(t, textResult.Text, "suggested-app")
	})

	t.Run("Repository error", func(t *testing.T) {
		applicationID := "error-app"
		expectedErr := errors.New("database error")
		repoMock.EXPECT().ListSLOsByApplication(ctx, applicationID).Return(nil, false, expectedErr).Times(1)

		req := createRequest("list_slos_by_application", map[string]interface{}{"application_id": applicationID})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "error"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "error listing slos of application error-app: database error")
	})

	t.Run("Missing application_id", func(t *testing.T) {
		req := createRequest("list_slos_by_application", map[string]interface{}{})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "invalid_input"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "Missing application_id")
	})
}
