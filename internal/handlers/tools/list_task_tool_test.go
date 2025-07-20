package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
)

func TestListTaskToolAndHandler_StoreError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	store := mystore.NewMockStore[model.Project](ctrl)
	store.EXPECT().List(ctx).Return(nil, assert.AnError)
	//store.EXPECT().Get(ctx, gomock.Any()).Return(model.Project{}, false, assert.AnError)

	// when
	tool := NewListTaskTool(store)
	result, err := tool.Handler(ctx, createRequest("list_tasks", nil))

	// then
	assert.NoError(t, err)
	expectError(t, result, "Error listing tasks")
}

func TestListTaskToolAndHandler_MissingProject(t *testing.T) {
	ctx := context.Background()

	// given
	store, _, _ := mystore.NewInMemoryStore[model.Project](ctx)
	store.Put(ctx, "A", model.Project{Name: "A", Description: "descA",
		Tasks: []model.TaskItem{
			{ProjectName: "B", ID: 1, Title: "taskA", Description: "descA"},
		},
	})

	// when
	tool := NewListTaskTool(store)
	result, err := tool.Handler(ctx, createRequest("list_tasks",
		map[string]interface{}{
			"project_name": "B",
		}))

	// then
	assert.NoError(t, err)
	expectError(t, result, "project B not found")
}

func TestListTaskToolAndHandler_Success(t *testing.T) {
	ctx := context.Background()

	// given
	store, _, _ := mystore.NewInMemoryStore[model.Project](ctx)
	store.Put(ctx, "A", model.Project{Name: "A", Description: "descA",
		Tasks: []model.TaskItem{
			{ProjectName: "A", ID: 1, Title: "taskA", Description: "descA"},
		},
	})

	// when
	tool := NewListTaskTool(store)
	result, err := tool.Handler(ctx, createRequest("list_tasks", nil))

	// then
	assert.NoError(t, err)
	expectSuccess(t, result, "taskA")
}
