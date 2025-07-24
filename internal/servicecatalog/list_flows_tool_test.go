package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
)

func TestListFlowsTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListFlows(gomock.Any()).Return([]string{"flow1", "flow2"}, nil)

	tool := NewMCPHandler(repo, nil).listFlowsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flows", nil))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	assert.Len(t, result.Content, 1)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "flow1")
	assert.Contains(t, textResult.Text, "flow2")
}

func TestListFlowsTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListFlows(gomock.Any()).Return(nil, errors.New("failed to list flows"))

	tool := NewMCPHandler(repo, nil).listFlowsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_flows", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "failed to list flows")
}
