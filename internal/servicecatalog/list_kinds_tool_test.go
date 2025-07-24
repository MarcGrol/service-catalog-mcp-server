package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestListKindsTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListKinds(gomock.Any()).Return([]string{"kind1", "kind2"}, nil)

	tool := NewMCPHandler(repo, nil).NewListKindsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_kinds", nil))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	assert.Len(t, result.Content, 1)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "kind1")
	assert.Contains(t, textResult.Text, "kind2")
}

func TestListKindsTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListKinds(gomock.Any()).Return(nil, errors.New("failed to list types"))

	tool := NewMCPHandler(repo, nil).NewListKindsTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_kinds", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error listing kinds")
}
