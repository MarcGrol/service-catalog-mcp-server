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

func TestListSLOByTeamTool(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := repo.NewMockSLORepo(ctrl)
	idxMock := slosearch.NewMockIndex(ctrl)

	h := NewMCPHandler(repoMock, idxMock)
	tool := h.listSLOByTeamTool()
	ctx := context.Background()

	t.Run("Successful list", func(t *testing.T) {
		teamID := "test-team"
		expectedSLOs := []repo.SLO{
			{UID: "slo1", Team: teamID},
			{UID: "slo2", Team: teamID},
		}
		repoMock.EXPECT().ListSLOsByTeam(ctx, teamID).Return(expectedSLOs, true, nil).Times(1)

		req := createRequest("list_slos_by_team",
			map[string]interface{}{"team_id": teamID},
		)
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectSuccess(t, result, `"status": "success"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, `"uid": "slo1",`)
	})

	t.Run("Team not found", func(t *testing.T) {
		teamID := "nonexistent-team"
		repoMock.EXPECT().ListSLOsByTeam(ctx, teamID).Return([]repo.SLO{}, false, nil).Times(1)
		idxMock.EXPECT().Search(ctx, teamID, gomock.Any()).Return(slosearch.Result{Teams: []string{"suggested-team"}}).Times(1)

		req := createRequest("list_slos_by_team",
			map[string]interface{}{"team_id": teamID},
		)
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "not_found"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "Team with ID nonexistent-team not found")
		assert.Contains(t, textResult.Text, "suggested-team")
	})

	t.Run("Repository error", func(t *testing.T) {
		teamID := "error-team"
		expectedErr := errors.New("database error")
		repoMock.EXPECT().ListSLOsByTeam(ctx, teamID).Return(nil, false, expectedErr).Times(1)

		req := createRequest("list_slos_by_team",
			map[string]interface{}{"team_id": teamID},
		)
		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "error"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "error listing slos of team error-team: database error")
	})

	t.Run("Missing team_id", func(t *testing.T) {
		req := createRequest("list_slos_by_team",
			map[string]interface{}{},
		)

		result, err := tool.Handler(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		expectError(t, result, `"status": "invalid_input"`)
		textResult := result.Content[0].(mcp.TextContent)
		assert.Contains(t, textResult.Text, "Missing team_id")
	})
}
