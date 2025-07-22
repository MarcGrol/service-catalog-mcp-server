package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/search"
)

func TestListInterfaceConsumersTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfaceConsumers(gomock.Any(), "interface1").Return([]string{"consumer1", "consumer2"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewListInterfaceConsumersTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interface_consumers", map[string]interface{}{
		"interface_id": "interface1",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "consumer1")
	assert.Contains(t, textResult.Text, "consumer2")
}

func TestListInterfaceConsumersTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfaceConsumers(gomock.Any(), "nonexistent_interface").Return(nil, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_interface", 10).Return(search.SearchResult{Interfaces: []string{"suggested_interface"}})

	tool := NewListInterfaceConsumersTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interface_consumers", map[string]interface{}{
		"interface_id": "nonexistent_interface",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Interface with ID nonexistent_interface not found")
	assert.Contains(t, textResult.Text, "suggested_interface")
}

func TestListInterfaceConsumersTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfaceConsumers(gomock.Any(), "interface_with_error").Return(nil, false, errors.New("failed to list consumers"))

	idx := search.NewMockIndex(ctrl)

	tool := NewListInterfaceConsumersTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interface_consumers", map[string]interface{}{
		"interface_id": "interface_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing consumers of interface interface_with_error: failed to list consumers")
}

func TestListInterfaceConsumersTool_MissingInterfaceID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewListInterfaceConsumersTool(repo, idx)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interface_consumers", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing interface_id")
}
