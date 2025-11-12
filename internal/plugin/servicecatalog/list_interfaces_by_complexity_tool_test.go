package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
)

func TestListInterfacesByComplexityTool_SuccessWithLimit(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := repo.NewMockCataloger(ctrl)
	repository.EXPECT().ListInterfacesByComplexity(gomock.Any(), 5).Return([]repo.Interface{
		{InterfaceID: "interface1", MethodCount: 10},
		{InterfaceID: "interface2", MethodCount: 5},
	}, nil)

	tool := NewMCPHandler(repository, nil).listInterfacesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces_by_complexity", map[string]interface{}{
		"limit_to": 5,
	}))

	// Then
	assert.NoError(t, err)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "interface1")
	assert.Contains(t, textResult.Text, "10")
	assert.Contains(t, textResult.Text, "interface2")
	assert.Contains(t, textResult.Text, "5")
}

func TestListInterfacesByComplexityTool_SuccessWithoutLimit(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := repo.NewMockCataloger(ctrl)
	repository.EXPECT().ListInterfacesByComplexity(gomock.Any(), 20).Return([]repo.Interface{
		{InterfaceID: "interfaceA", MethodCount: 100},
		{InterfaceID: "interfaceB", MethodCount: 50},
	}, nil)

	tool := NewMCPHandler(repository, nil).listInterfacesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces_by_complexity", nil))

	// Then
	assert.NoError(t, err)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "interfaceA")
	assert.Contains(t, textResult.Text, "100")
	assert.Contains(t, textResult.Text, "interfaceB")
	assert.Contains(t, textResult.Text, "50")
}

func TestListInterfacesByComplexityTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfacesByComplexity(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to list interfaces"))

	tool := NewMCPHandler(repo, nil).listInterfacesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces_by_complexity", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing interfaces by complexity: failed to list interfaces")
}
