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

func TestListFlowParticipantsTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListParticpantsOfFlow(gomock.Any(), "flow1").Return([]string{"participant1", "participant2"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewListFlowParticipantsTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flow_participants", map[string]interface{}{
		"flow_id": "flow1",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "participant1")
	assert.Contains(t, textResult.Text, "participant2")
}

func TestListFlowParticipantsTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListParticpantsOfFlow(gomock.Any(), "nonexistent_flow").Return(nil, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_flow", 10).Return(search.SearchResult{Flows: []string{"suggested_flow"}})

	tool := NewListFlowParticipantsTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flow_participants", map[string]interface{}{
		"flow_id": "nonexistent_flow",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Flow with ID nonexistent_flow not found")
	assert.Contains(t, textResult.Text, "suggested_flow")
}

func TestListFlowParticipantsTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListParticpantsOfFlow(gomock.Any(), "flow_with_error").Return(nil, false, errors.New("failed to list participants"))

	idx := search.NewMockIndex(ctrl)

	tool := NewListFlowParticipantsTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flow_participants", map[string]interface{}{
		"flow_id": "flow_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing participants of flow flow_with_error: failed to list participants")
}

func TestListFlowParticipantsTool_MissingFlowID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewListFlowParticipantsTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flow_participants", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing flow_id")
}
