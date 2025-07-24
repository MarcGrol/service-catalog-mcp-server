package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestListModulesOfTeamsTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesOfTeam(gomock.Any(), "team1").Return([]string{"module1", "module2"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listModulesOfTeamsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_of_teams", map[string]interface{}{
		"team_id": "team1",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "module1")
	assert.Contains(t, textResult.Text, "module2")
}

func TestListModulesOfTeamsTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesOfTeam(gomock.Any(), "nonexistent_team").Return(nil, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_team", 10).Return(search.Result{Teams: []string{"suggested_team"}})

	tool := NewMCPHandler(repo, idx).listModulesOfTeamsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_of_teams", map[string]interface{}{
		"team_id": "nonexistent_team",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Team with ID nonexistent_team not found")
	assert.Contains(t, textResult.Text, "suggested_team")
}

func TestListModulesOfTeamsTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesOfTeam(gomock.Any(), "team_with_error").Return(nil, false, errors.New("failed to list modules"))

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listModulesOfTeamsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_of_teams", map[string]interface{}{
		"team_id": "team_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing modules of team team_with_error: failed to list modules")
}

func TestListModulesOfTeamsTool_MissingTeamID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listModulesOfTeamsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_of_teams", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing team_id")
}
