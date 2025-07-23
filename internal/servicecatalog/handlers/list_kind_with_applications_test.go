package handlers

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

func TestListKindWithApplicationsTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesWithKind(gomock.Any(), "kind1").Return([]string{"participant1", "participant2"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewListModulesWithKindTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_kind_participants", map[string]interface{}{
		"kind_id": "kind1",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "participant1")
	assert.Contains(t, textResult.Text, "participant2")
}

func TestListKindWithApplicationsTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesWithKind(gomock.Any(), "nonexistent_kind").Return(nil, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_kind", 10).Return(search.SearchResult{Kinds: []string{"suggested_kind"}})

	tool := NewListModulesWithKindTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_kind_participants", map[string]interface{}{
		"kind_id": "nonexistent_kind",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "No modules found for kind with ID nonexistent_kind")
	assert.Contains(t, textResult.Text, "suggested_kind")
}

func TestListKindParticipantsTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesWithKind(gomock.Any(), "kind_with_error").Return(nil, false, errors.New("failed to list kinds"))

	idx := search.NewMockIndex(ctrl)

	tool := NewListModulesWithKindTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_with_kind", map[string]interface{}{
		"kind_id": "kind_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing modules with kind")
}

func TestListKindParticipantsTool_MissingKindID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewListModulesWithKindTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_kind_participants", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing kind_id")
}
