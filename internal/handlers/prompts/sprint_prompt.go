package prompts

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewSprintPrompt returns the MCP prompt contract and handler for sprint planning.
func NewSprintPrompt() server.ServerPrompt {
	return server.ServerPrompt{
		Prompt: mcp.NewPrompt(
			"sprint_planning",
			mcp.WithPromptDescription("Assist with agile sprint planning and task breakdown"),
			mcp.WithArgument("sprint_length", mcp.ArgumentDescription("Sprint length in days")),
			mcp.WithArgument("team_size", mcp.ArgumentDescription("Number of team members")),
		),
		Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
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
			promptText := fmt.Sprintf(`You are a sprint planning assistant for an agile team.

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
							Text: promptText,
						},
					},
				},
			}, nil
		},
	}
}
