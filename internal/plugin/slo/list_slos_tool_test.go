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

func TestListSLOTool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockSLORepo(ctrl)
	idxMock := slosearch.NewMockIndex(ctrl)

	h := NewMCPHandler(repoMock, idxMock)
	tool := h.listSLOTool()
	ctx := context.Background()

	t.Run("Successful list", func(t *testing.T) {
		expectedSLOs := []repo.SLO{
			{UID: "slo1", Team: "test-team"},
			{UID: "slo2", Team: "test-team"},
		}
		repoMock.EXPECT().ListSLOs(ctx).Return(expectedSLOs, nil).Times(1)

		req := createRequest("list_slos", map[string]interface{}{})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectSuccess(t, result, `"status": "success"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, `"uid": "slo1",`)
	})

	t.Run("Repository error", func(t *testing.T) {
		expectedErr := errors.New("database error")
		repoMock.EXPECT().ListSLOs(ctx).Return(nil, expectedErr).Times(1)

		req := createRequest("list_slos_by_team", map[string]interface{}{})
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "error"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "error listing slos: database error")
	})
}
