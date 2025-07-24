package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
)

func TestListModulesByComplexityTool_SuccessWithLimit(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesByCompexity(gomock.Any(), 5).Return([]catalogrepo.Module{
		{ModuleID: "module1", Name: "Module One", Description: "Desc One", ComplexityScore: 10.5},
		{ModuleID: "module2", Name: "Module Two", Description: "Desc Two", ComplexityScore: 8.2},
	}, nil)

	tool := NewMCPHandler(repo, nil).listModulesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_by_complexity", map[string]interface{}{
		"limit_to": 5,
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "module1")
	assert.Contains(t, textResult.Text, "Module One")
	assert.Contains(t, textResult.Text, "10.5")
	assert.Contains(t, textResult.Text, "module2")
	assert.Contains(t, textResult.Text, "Module Two")
	assert.Contains(t, textResult.Text, "8.2")
}

func TestListModulesByComplexityTool_SuccessWithoutLimit(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesByCompexity(gomock.Any(), 20).Return([]catalogrepo.Module{
		{ModuleID: "moduleA", Name: "Module A", Description: "Desc A", ComplexityScore: 50.1},
		{ModuleID: "moduleB", Name: "Module B", Description: "Desc B", ComplexityScore: 30.9},
	}, nil)

	tool := NewMCPHandler(repo, nil).listModulesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_by_complexity", nil))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "moduleA")
	assert.Contains(t, textResult.Text, "Module A")
	assert.Contains(t, textResult.Text, "50.1")
	assert.Contains(t, textResult.Text, "moduleB")
	assert.Contains(t, textResult.Text, "Module B")
	assert.Contains(t, textResult.Text, "30.9")
}

func TestListModulesByComplexityTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModulesByCompexity(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to list modules"))

	tool := NewMCPHandler(repo, nil).listModulesByComplexityTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules_by_complexity", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing modules by complexity: failed to list modules")
}
