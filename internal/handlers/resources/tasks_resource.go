package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewTasksListResource returns the MCP resource contract and handler for the tasks list.
func NewTasksListResource(store mystore.Store[model.Project]) Resouce {
	resource := mcp.NewResource(
		"tasks://list",
		"List of all tasks in the project",
		mcp.WithMIMEType("application/json"),
	)
	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		tasks := []model.TaskItem{}
		projects, err := store.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing projects: %s", err)
		}
		for _, p := range projects {
			tasks = append(tasks, p.Tasks...)
		}
		tasksJSON, err := json.MarshalIndent(map[string]interface{}{
			"total_tasks":  len(tasks),
			"tasks":        tasks,
			"last_updated": time.Now().Format(time.RFC3339),
		}, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error serializing results: %s", err)
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(tasksJSON),
			},
		}, nil
	}
	return Resouce{
		Resource: resource,
		Handler:  handler,
	}
}
