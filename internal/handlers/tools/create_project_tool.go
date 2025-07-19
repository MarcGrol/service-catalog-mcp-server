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

// NewCreateProjectTool returns the MCP tool definition and its handler implementation together.
func NewCreateProjectTool(store mystore.Store[model.Project]) Tool {
	return Tool{
		Contract: mcp.NewTool(
			"create_project",
			mcp.WithDescription("Create a new project configuration"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
			mcp.WithString("description", mcp.Required(), mcp.Description("Project description")),
			mcp.WithArray("authors", mcp.Description("List of project authors")),
		),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			name, err := request.RequireString("name")
			if err != nil {
				return mcp.NewToolResultError("Missing project name"), nil
			}

			description, err := request.RequireString("description")
			if err != nil {
				return mcp.NewToolResultError("Missing project description"), nil
			}

			authors := request.GetStringSlice("authors", []string{"Anonymous Developer"})

			projectConfig := model.Project{
				Name:        name,
				Version:     "1.0.0",
				Description: description,
				Authors:     authors,
				Dependencies: map[string]string{
					"golang": "1.21+",
				},
				CreatedAt: time.Now(),
			}

			err = store.Put(ctx, name, projectConfig)
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
		},
	}
}
