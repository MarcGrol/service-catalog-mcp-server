package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/repo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
)

func TestGetModuleTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := repo.NewMockCataloger(ctrl)
	repository.EXPECT().GetModuleOnID(gomock.Any(), "module1").Return(repo.Module{ModuleID: "module1", Name: "Test Module"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repository, idx).getSingleModuleTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("get_module", map[string]interface{}{
		"module_id": "module1",
	}))

	// Then
	assert.NoError(t, err)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "module1")
	assert.Contains(t, textResult.Text, "Test Module")
}

func TestGetModuleTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := repo.NewMockCataloger(ctrl)
	repository.EXPECT().GetModuleOnID(gomock.Any(), "nonexistent_module").Return(repo.Module{}, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_module", 10).Return(search.Result{Modules: []string{"suggested_module"}})

	tool := NewMCPHandler(repository, idx).getSingleModuleTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("get_module", map[string]interface{}{
		"module_id": "nonexistent_module",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Module with ID nonexistent_module not found")
	assert.Contains(t, textResult.Text, "suggested_module")
}

func TestGetModuleTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := repo.NewMockCataloger(ctrl)
	repository.EXPECT().GetModuleOnID(gomock.Any(), "module_with_error").Return(repo.Module{}, false, errors.New("failed to get module"))

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repository, idx).getSingleModuleTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("get_module", map[string]interface{}{
		"module_id": "module_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error getting module module_with_error: failed to get module")
}

func TestGetModuleTool_MissingModuleID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).getSingleModuleTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("get_module", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing module_id")
}
