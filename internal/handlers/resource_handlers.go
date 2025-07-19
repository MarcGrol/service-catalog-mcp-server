package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
)

type ResourceHandlers struct {
	store mystore.Store[model.Project]
}

func NewResourceHandlers(store mystore.Store[model.Project]) *ResourceHandlers {
	return &ResourceHandlers{
		store: store,
	}
}

func (h *ResourceHandlers) ProjectResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		projects, err := h.store.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing projects: %s", err)
		}

		projectsJson, err := json.MarshalIndent(
			map[string]interface{}{
				"total_projects": len(projects),
				"projects":       projects,
				"last_updated":   time.Now().Format(time.RFC3339),
			}, "", "  ")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(projectsJson),
			},
		}, nil
	}
}

func (h *ResourceHandlers) TasksResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		tasks := []model.TaskItem{}

		projects, err := h.store.List(ctx)
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
}

func (h *ResourceHandlers) StatsResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		tasks := []model.TaskItem{}

		projects, err := h.store.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing projects: %s", err)
		}

		for _, p := range projects {
			tasks = append(tasks, p.Tasks...)
		}

		// CALCULATE from current data
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
}

func (h *ResourceHandlers) DocsResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		readme := `# Advanced Project Manager MCP Server

This is an advanced example of an MCP server built with Go that demonstrates:

## Features

- **Tools**: Project and task management
- **Resources**: Configuration, task lists, statistics
- **Prompts**: Planning and review assistance

## Usage

1. Create projects with "create_project"
2. Add tasks with "create_task"
3. Search content with "search_content"
4. Generate reports with "generate_analytics"

## Resources Available

- "project://config" - Current project configuration
- "tasks://list" - All project tasks
- "stats://project" - Project statistics
- "docs://readme" - This documentation

## Getting Started

Connect this server to your MCP client (Claude Desktop, Cursor, etc.) and start managing your projects with AI assistance!
`

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "text/markdown",
				Text:     readme,
			},
		}, nil
	}
}
