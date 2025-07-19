package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

type PromptHandlers struct {
}

func NewPromptHandlers() *PromptHandlers {
	return &PromptHandlers{}
}

func (h *PromptHandlers) PlanningPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		var args map[string]string = request.Params.Arguments

		projectType := "web" // default
		if pt, ok := args["project_type"]; ok {
			projectType = pt
		}

		timeline := "4" // default weeks
		if tl, ok := args["timeline"]; ok {
			timeline = tl
		}

		// For comma-separated technologies
		technologiesStr := args["technologies"]
		technologies := []string{"Go", "React"}
		if technologiesStr != "" {
			technologies = strings.Split(technologiesStr, ",")
			// Trim whitespace
			for i, tech := range technologies {
				technologies[i] = strings.TrimSpace(tech)
			}
		}
		prompt := fmt.Sprintf(`You are a project planning assistant. Help plan a %s project.

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
			Description: "Code review guidance",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		}, nil
	}
}

func (h *PromptHandlers) ReviewPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments

		language := args["language"]
		if language == "" {
			language = "Go"
		}

		focus := args["focus"]
		if focus == "" {
			focus = "all"
		}
		prompt := fmt.Sprintf(`You are a code review assistant for %s code.

Review Focus: %s

Please review the following code and provide feedback on:

1. **Code Quality**: Structure, readability, and maintainability
2. **Best Practices**: Language-specific conventions and patterns
3. **Performance**: Potential bottlenecks and optimization opportunities
4. **Security**: Vulnerability assessment and security best practices
5. **Testing**: Test coverage and testing strategies

For %s specifically, pay attention to:
- Error handling patterns
- Memory management
- Concurrency safety
- Package structure and naming conventions

Provide constructive feedback with specific suggestions for improvement.`,
			language, focus, language)

		return &mcp.GetPromptResult{
			Description: "Code review guidance",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		}, nil
	}
}

func (h *PromptHandlers) SprintPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments

		sprintLengthStr := args["sprint_length"]
		sprintLength := 14.0
		if sprintLengthStr != "" {
			if parsed, err := strconv.ParseFloat(sprintLengthStr, 64); err == nil {
				sprintLength = parsed
			}
		}

		teamSizeStr := args["team_size"]
		teamSize := 4.0
		if teamSizeStr != "" {
			if parsed, err := strconv.ParseFloat(teamSizeStr, 64); err == nil {
				teamSize = parsed
			}
		}

		prompt := fmt.Sprintf(`You are a sprint planning assistant for an agile team.

Sprint Configuration:
- Sprint Length: %.0f days
- Team Size: %.0f members

Help plan the upcoming sprint by:

1. **Capacity Planning**: Calculate team velocity and available story points
2. **Task Breakdown**: Break epics into manageable user stories
3. **Priority Assessment**: Rank tasks by business value and dependencies
4. **Risk Identification**: Identify potential blockers and mitigation strategies
5. **Timeline Creation**: Suggest daily standup goals and milestone checkpoints

Consider team capacity at 70%%%% for %.0f days = %.1f effective working days.
With %.0f team members, total capacity is approximately %.0f story points (assuming 1 point = 0.5 days).

Focus on creating realistic, achievable sprint goals that deliver value to stakeholders.`,
			sprintLength, teamSize, sprintLength, sprintLength*0.7, teamSize, teamSize*sprintLength*0.7*2)

		return &mcp.GetPromptResult{
			Description: "Sprint planning guidance",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		}, nil
	}
}
