package tools

import (
	"context"
	"fmt"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewGenerateAnalyticsTool returns the MCP tool definition and its handler for generating analytics.
func NewGenerateAnalyticsTool(store mystore.Store[model.Project]) Tool {
	tool := mcp.NewTool(
		"generate_analytics",
		mcp.WithDescription("Generate project analytics and reports"),
		mcp.WithString("report_type", mcp.Required(), mcp.Description("Type of report: summary, tasks, timeline")),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reportType, err := request.RequireString("report_type")
		if err != nil {
			return mcp.NewToolResultError("Missing report type"), nil
		}
		var report string
		switch reportType {
		case "summary":
			report = `Project Summary Report\n========================\n- Total Projects: 5\n- Active Tasks: 12\n- Completed Tasks: 8\n- Team Members: 4\n- Sprint Progress: 75%\n- Code Coverage: 85%`
		case "tasks":
			report = `Task Analysis Report\n===================\nHigh Priority: 3 tasks\nMedium Priority: 7 tasks\nLow Priority: 2 tasks\n\nStatus Distribution:\n- Todo: 5 tasks\n- In Progress: 4 tasks\n- In Review: 2 tasks\n- Done: 1 task`
		case "timeline":
			report = `Timeline Report\n==============\nWeek 1: Project setup and initial planning\nWeek 2: Core feature development (75% complete)\nWeek 3: Testing and refinement (planned)\nWeek 4: Documentation and deployment (planned)\n\nMilestones:\n✓ Project kickoff\n✓ Architecture design\n⧗ MVP completion (in progress)\n◯ Beta release (upcoming)`
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown report type: %s", reportType)), nil
		}
		return mcp.NewToolResultText(report), nil
	}
	return Tool{
		Tool:    tool,
		Handler: handler,
	}
}
