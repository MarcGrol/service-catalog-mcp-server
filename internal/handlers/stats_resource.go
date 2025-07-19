package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewStatsResourceAndHandler returns the MCP resource contract and handler for project statistics.
func NewStatsResourceAndHandler(store mystore.Store[model.Project]) (mcp.Resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)) {
	resource := mcp.NewResource(
		"stats://project",
		"Project statistics and metrics",
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
		todoTasks := 0
		inProgressTasks := 0
		doneTasks := 0
		highPriorityTasks := 0
		for _, task := range tasks {
			switch task.Status {
			case "todo":
				todoTasks++
			case "in_progress":
				inProgressTasks++
			case "done":
				doneTasks++
			}
			if task.Priority == "high" {
				highPriorityTasks++
			}
		}
		stats := map[string]interface{}{
			"total_projects":      len(projects),
			"total_tasks":         len(tasks),
			"tasks_todo":          todoTasks,
			"tasks_in_progress":   inProgressTasks,
			"tasks_done":          doneTasks,
			"high_priority_tasks": highPriorityTasks,
			"completion_rate":     float64(doneTasks) / float64(len(tasks)) * 100,
			"last_calculated":     time.Now().Format(time.RFC3339),
		}
		statsJSON, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error serializing results: %s", err)
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(statsJSON),
			},
		}, nil
	}
	return resource, handler
}
