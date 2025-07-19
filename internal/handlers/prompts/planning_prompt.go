package prompts

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewPlanningPrompt returns the MCP prompt contract and handler for project planning.
func NewPlanningPrompt() Prompt {
	return Prompt{
		Contract: mcp.NewPrompt(
			"project_planning",
			mcp.WithPromptDescription("Help plan and structure a new project"),
			mcp.WithArgument("project_type", mcp.RequiredArgument(), mcp.ArgumentDescription("Type of project: web, mobile, api, library")),
			mcp.WithArgument("timeline", mcp.ArgumentDescription("Project timeline in weeks")),
			mcp.WithArgument("technologies", mcp.ArgumentDescription("Preferred technologies")),
		),
		Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			args := request.Params.Arguments
			projectType := args["project_type"]
			if projectType == "" {
				projectType = "web"
			}
			timeline := args["timeline"]
			if timeline == "" {
				timeline = "4"
			}
			technologiesStr := args["technologies"]
			technologies := []string{"Go", "React"}
			if technologiesStr != "" {
				technologies = strings.Split(technologiesStr, ",")
				for i, tech := range technologies {
					technologies[i] = strings.TrimSpace(tech)
				}
			}
			promptText := fmt.Sprintf(`You are a project planning assistant. Help plan a %s project.

Project Details:
- Type: %s
- Timeline: %s weeks
- Technologies: %s

Please provide:
1. Project structure and architecture recommendations
2. Key milestones and deliverables
3. Risk assessment and mitigation strategies
4. Resource allocation suggestions
5. Technology stack recommendations

Consider best practices for %s development and provide actionable advice.`,
				projectType, projectType, timeline, strings.Join(technologies, ", "), projectType)
			return &mcp.GetPromptResult{
				Description: "Project planning guidance",
				Messages: []mcp.PromptMessage{
					{
						Role: mcp.RoleUser,
						Content: mcp.TextContent{
							Type: "text",
							Text: promptText,
						},
					},
				},
			}, nil
		},
	}
}
