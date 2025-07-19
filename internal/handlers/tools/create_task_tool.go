package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewCreateTaskTool returns the MCP tool definition and its handler for creating tasks.
func NewCreateTaskTool(store mystore.Store[model.Project]) Tool {
	return Tool{
		Contract: mcp.NewTool(
			"create_task",
			mcp.WithDescription("Create a new task"),
			mcp.WithString("project_name", mcp.Required(), mcp.Description("Project that this task must be added to")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Task title")),
			mcp.WithString("description", mcp.Description("Task description")),
			mcp.WithString("priority", mcp.Description("Task priority: low, medium, high")),
			mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			task := model.TaskItem{
				ProjectName: projectName,
				ID:          int(time.Now().Unix()),
				Title:       title,
				Description: description,
				Status:      "todo",
				Priority:    priority,
				CreatedAt:   time.Now(),
			}
			dueDateStr := request.GetString("due_date", "")
			if dueDateStr != "" {
				if dueDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
					task.DueDate = &dueDate
				}
			}
			proj, exists, err := store.Get(ctx, projectName)
			if err != nil {
				return nil, err
			}
			if !exists {
				return mcp.NewToolResultError(fmt.Sprintf("project %s not found", projectName)), nil
			}
			proj.Tasks = append(proj.Tasks, task)
			err = store.Put(ctx, projectName, proj)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error storing task", err), nil
			}
			taskJSON, err := json.MarshalIndent(task, "", "  ")
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error serializing task", err), nil
			}
			result := fmt.Sprintf("Task created successfully!\n\n%s", string(taskJSON))
			return mcp.NewToolResultText(result), nil
		},
	}
}
