package servicecatalog

import (
	"context"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/catalogrepo"
	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/servicecatalog/search"
)

func TestListDatabaseConsumersTool_Success(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListDatabaseConsumers(gomock.Any(), "db1").Return([]string{"consumer1", "consumer2"}, true, nil)

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listMDatabaseConsumersTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_database_consumers", map[string]interface{}{
		"database_id": "db1",
	}))

	// Then
	assert.NoError(t, err)
	expectSuccess(t, result, `"status": "success"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "consumer1")
	assert.Contains(t, textResult.Text, "consumer2")
}

func TestListDatabaseConsumersTool_NotFound(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListDatabaseConsumers(gomock.Any(), "nonexistent_db").Return(nil, false, nil)

	idx := search.NewMockIndex(ctrl)
	idx.EXPECT().Search(gomock.Any(), "nonexistent_db", 10).Return(search.Result{Databases: []string{"suggested_db"}})

	tool := NewMCPHandler(repo, idx).listMDatabaseConsumersTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_database_consumers", map[string]interface{}{
		"database_id": "nonexistent_db",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "not_found"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Module with ID nonexistent_db not found")
	assert.Contains(t, textResult.Text, "suggested_db")
}

func TestListDatabaseConsumersTool_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	repo.EXPECT().ListDatabaseConsumers(gomock.Any(), "db_with_error").Return(nil, false, errors.New("failed to list consumers"))

	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listMDatabaseConsumersTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_database_consumers", map[string]interface{}{
		"database_id": "db_with_error",
	}))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "error"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "error getting database db_with_error: failed to list consumers")
}

func TestListDatabaseConsumersTool_MissingDatabaseID(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := catalogrepo.NewMockCataloger(ctrl)
	idx := search.NewMockIndex(ctrl)

	tool := NewMCPHandler(repo, idx).listMDatabaseConsumersTool()

	// When
	result, err := tool.Handler(context.Background(), createRequest("list_database_consumers", nil))

	// Then
	assert.NoError(t, err)
	expectError(t, result, `"status": "invalid_input"`)
	textResult := result.Content[0].(mcp.TextContent)
	assert.Contains(t, textResult.Text, "Missing database_id")
}
