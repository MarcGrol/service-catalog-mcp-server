package handlers

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

func TestListInterfacesTool_SuccessWithKeyword(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfaces(gomock.Any(), "test").Return([]catalogrepo.Interface{
		{InterfaceID: "interface1", Description: "desc1", Kind: "kind1"},
		{InterfaceID: "interface2", Description: "desc2", Kind: "kind2"},
	}, nil)

	tool := NewListInterfacesTool(repo)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces", map[string]interface{}{
		"filter_keyword": "test",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "interface1")
	assert.Contains(t, textResult.Text, "desc1")
	assert.Contains(t, textResult.Text, "kind1")
	assert.Contains(t, textResult.Text, "interface2")
	assert.Contains(t, textResult.Text, "desc2")
	assert.Contains(t, textResult.Text, "kind2")
}

func TestListInterfacesTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListInterfaces(gomock.Any(), "error").Return(nil, errors.New("failed to list interfaces"))

	tool := NewListInterfacesTool(repo)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces", map[string]interface{}{
		"filter_keyword": "error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing interfaces with keyword: failed to list interfaces")
}

func TestListInterfacesTool_MissingKeyword(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)

	tool := NewListInterfacesTool(repo)

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_interfaces", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing filter_keyword")
}
