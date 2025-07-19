package tools

import (
	"context"
	"testing"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListProjectToolAndHandler_StoreError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// given
	store := mystore.NewMockStore[model.Project](ctrl)
	store.EXPECT().List(ctx).Return(nil, assert.AnError)

	// when
	tool := NewListProjectTool(store)
	result, err := tool.Handler(ctx, createRequest("list_projects", nil))

	// then
	assert.NoError(t, err)
	expectError(t, result, "Error listing projects")
}

func TestListProjectToolAndHandler_Success(t *testing.T) {
	ctx := context.Background()

	// given
	store, _, _ := mystore.NewInMemoryStore[model.Project](ctx)
	store.Put(ctx, "A", model.Project{Name: "A", Description: "descA"})

	// when
	tool := NewListProjectTool(store)
	result, err := tool.Handler(ctx, createRequest("list_projects", nil))

	// then
	assert.NoError(t, err)
	expectSuccess(t, result, "A: descA")
}
