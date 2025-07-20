package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewReviewPrompt returns the MCP prompt contract and handler for code review.
func NewReviewPrompt() server.ServerPrompt {
	return server.ServerPrompt{
		Prompt: mcp.NewPrompt(
			"code_review",
			mcp.WithPromptDescription("Generate code review guidelines and checklist"),
			mcp.WithArgument("language", mcp.RequiredArgument(), mcp.ArgumentDescription("Programming language")),
			mcp.WithArgument("focus", mcp.ArgumentDescription("Review focus: security, performance, style, all")),
		),
		Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			args := request.Params.Arguments
			language := args["language"]
			if language == "" {
				language = "Go"
			}
			focus := args["focus"]
			if focus == "" {
				focus = "all"
			}
			promptText := fmt.Sprintf(`You are a code review assistant for %s code.

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
							Text: promptText,
						},
					},
				},
			}, nil
		},
	}
}
