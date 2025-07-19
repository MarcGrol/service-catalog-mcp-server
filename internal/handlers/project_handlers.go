package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/mystore"
	"github.com/MarcGrol/learnmcp/pkg/models"
)

type ProjectHandlers struct {
	mcpServer *server.MCPServer
	store     mystore.Store[models.ProjectConfig]
}

func NewProjectHandlers(s *server.MCPServer, store mystore.Store[models.ProjectConfig]) *ProjectHandlers {
	return &ProjectHandlers{
		mcpServer: s,
		store:     store,
	}
}

func (h *ProjectHandlers) ListProjectHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects, err := h.store.List(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error listing projects", err), err
		}

		// Simulate search results
		results := []string{}
		for _, p := range projects {
			results = append(results, fmt.Sprintf("%s: %s", p.Name, p.Description))
		}

		result := fmt.Sprintf("Currently available project:\n\n%s",
			strings.Join(results, "\n"))

		return mcp.NewToolResultText(result), nil
	}
}

func (h *ProjectHandlers) CreateProjectHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("Missing project name"), nil
		}

		description, err := request.RequireString("description")
		if err != nil {
			return mcp.NewToolResultError("Missing project description"), nil
		}

		authors := request.GetStringSlice("authors", []string{"Anonymous Developer"})

		projectConfig := models.ProjectConfig{
			Name:        name,
			Version:     "1.0.0",
			Description: description,
			Authors:     authors,
			Dependencies: map[string]string{
				"golang": "1.21+",
			},
			CreatedAt: time.Now(),
		}

		// Save project
		err = h.store.Put(ctx, name, projectConfig)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error storing project", err), nil
		}

		projectJSON, err := json.MarshalIndent(projectConfig, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error serializing project", err), nil
		}

		result := fmt.Sprintf("Project '%s' created successfully!\n\nConfiguration:\n%s",
			name, string(projectJSON))

		return mcp.NewToolResultText(result), nil
	}
}

func (h *ProjectHandlers) ListTaskHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		if projectName == "" {
			projects, err := h.store.List(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing projects", err), err
			}
			tasks := []string{}
			for _, p := range projects {
				for _, t := range p.Tasks {
					tasks = append(tasks, fmt.Sprintf("%s: %s - %s - %s", p.Name, t.ID, t.Title, t.Description))
				}
			}

			result := fmt.Sprintf("Currently available tasks:\n\n%s", strings.Join(tasks, "\n"))

			return mcp.NewToolResultText(result), nil
		}

		project, exists, err := h.store.Get(ctx, projectName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("project %s not found", projectName)), nil
		}

		// Simulate search results
		results := []string{}
		for _, t := range project.Tasks {
			results = append(results, fmt.Sprintf("%s: %s - %s - %s", t.ID, projectName, t.Title, t.Description))
		}

		result := fmt.Sprintf("Currently available tasks within project %s:\n\n%s", projectName,
			strings.Join(results, "\n"))

		return mcp.NewToolResultText(result), nil
	}
}

func (h *ProjectHandlers) CreateTaskHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName, err := request.RequireString("project_name")
		if err != nil {
			return mcp.NewToolResultError("Missing project_name"), nil
		}

		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError("Missing task title"), nil
		}

		description := request.GetString("description", "")
		priority := request.GetString("priority", "medium")

		task := models.TaskItem{
			ProjectName: projectName,
			ID:          int(time.Now().Unix()),
			Title:       title,
			Description: description,
			Status:      "todo",
			Priority:    priority,
			CreatedAt:   time.Now(),
		}

		// Parse due date if provided
		dueDateStr := request.GetString("due_date", "")
		if dueDateStr != "" {
			if dueDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
				task.DueDate = &dueDate
			}
		}

		// Search project
		proj, exists, err := h.store.Get(ctx, projectName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("project %s not found", projectName)), nil
		}

		proj.Tasks = append(proj.Tasks, task)

		// Save project
		err = h.store.Put(ctx, projectName, proj)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error storing task", err), nil
		}

		taskJSON, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error serializing task", err), nil
		}
		result := fmt.Sprintf("Task created successfully!\n\n%s", string(taskJSON))

		return mcp.NewToolResultText(result), nil
	}
}
