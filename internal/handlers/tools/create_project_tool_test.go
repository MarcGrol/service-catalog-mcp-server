package tools

import (
	"context"
	"testing"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateProjectToolAndHandler_InvalidInput(t *testing.T) {
	ctx := context.Background()

	//given
	store, _, _ := mystore.NewInMemoryStore[model.Project](ctx)

	cases := []struct {
		name      string
		args      map[string]interface{}
		expectErr string
	}{
		{"missing name", map[string]interface{}{"description": "desc"}, "Missing project name"},
		{"missing description", map[string]interface{}{"name": "proj"}, "Missing project description"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// when
			tool := NewCreateProjectTool(store)
			result, err := tool.Handler(ctx, createRequest("create_project", tc.args))
			// then
			assert.NoError(t, err)
			assert.True(t, result.IsError)
			assert.NoError(t, err)
			expectError(t, result, "Missing ")
			expectEmpty(t, ctx, store)
		})
	}
}

func TestCreateProjectToolAndHandler_StoreError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	store := mystore.NewMockStore[model.Project](ctrl)
	store.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(assert.AnError)
	store.EXPECT().List(gomock.Any()).Return([]model.Project{}, nil)

	// when
	tool := NewCreateProjectTool(store)
	result, err := tool.Handler(ctx, createRequest("create_project",
		map[string]interface{}{
			"name":        "proj",
			"description": "desc",
		}))

	// then
	assert.NoError(t, err)
	expectError(t, result, "Error storing project")
	expectEmpty(t, ctx, store)
}

func TestCreateProjectToolAndHandler_Success(t *testing.T) {
	ctx := context.Background()

	// given
	store, _, _ := mystore.NewInMemoryStore[model.Project](ctx)
	tool := NewCreateProjectTool(store)

	// when
	result, err := tool.Handler(ctx, createRequest("create_project", map[string]interface{}{
		"name":        "proj",
		"description": "desc",
		"authors":     []string{"A"},
	}))

	// then
	assert.NoError(t, err)
	expectSuccess(t, result, "Project 'proj' created successfully")
	expectProject(t, ctx, store, model.Project{
		Name:        "proj",
		Description: "desc",
		Authors:     []string{"A"},
	})
}
