package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogrepo"
)

func TestListModulesTool_SuccessWithKeyword(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModules(gomock.Any(), "test").Return([]catalogrepo.Module{
		{ModuleID: "module1", Name: "Module One", Description: "Desc One"},
		{ModuleID: "module2", Name: "Module Two", Description: "Desc Two"},
	}, nil)

	tool := NewMCPHandler(repo, nil).listModulesTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules", map[string]interface{}{
		"filter_keyword": "test",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "module1")
	assert.Contains(t, textResult.Text, "Module One")
	assert.Contains(t, textResult.Text, "module2")
	assert.Contains(t, textResult.Text, "Module Two")
}

func TestListModulesTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListModules(gomock.Any(), "error").Return(nil, errors.New("failed to list modules"))

	tool := NewMCPHandler(repo, nil).listModulesTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules", map[string]interface{}{
		"filter_keyword": "error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing modules with keyword error: failed to list modules")
}

func TestListModulesTool_MissingKeyword(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)

	tool := NewMCPHandler(repo, nil).listModulesTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_modules", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing filter_keyword")
}
