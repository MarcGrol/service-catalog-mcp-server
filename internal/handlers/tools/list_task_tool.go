package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewListTaskTool returns the MCP tool definition and its handler for listing tasks.
func NewListTaskTool(store mystore.Store[model.Project]) Tool {
	return Tool{
		Contract: mcp.NewTool(
			"list_tasks",
			mcp.WithDescription("Lists all tasks or all tasks of a project"),
			mcp.WithString("project_name", mcp.Description("Project that we want to list the tasks of")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projectName := request.GetString("project_name", "")
			if projectName == "" {
				projects, err := store.List(ctx)
				if err != nil {
					return mcp.NewToolResultErrorFromErr("Error listing tasks", err), nil
				}
				tasks := []string{}
				for _, p := range projects {
					for _, t := range p.Tasks {
						tasks = append(tasks, fmt.Sprintf("%s: %d - %s - %s", p.Name, t.ID, t.Title, t.Description))
					}
				}
				result := fmt.Sprintf("Currently available tasks:\n\n%s", strings.Join(tasks, "\n"))
				return mcp.NewToolResultText(result), nil
			}
			project, exists, err := store.Get(ctx, projectName)
			if err != nil {
				return mcp.NewToolResultText(fmt.Sprintf("error getting project %s: %s", projectName, err)), nil
			}
			if !exists {
				return mcp.NewToolResultError(fmt.Sprintf("project %s not found", projectName)), nil
			}
			results := []string{}
			for _, t := range project.Tasks {
				results = append(results, fmt.Sprintf("%d: %s - %s - %s", t.ID, projectName, t.Title, t.Description))
			}
			result := fmt.Sprintf("Currently available tasks within project %s:\n\n%s", projectName, strings.Join(results, "\n"))
			return mcp.NewToolResultText(result), nil
		},
	}
}
