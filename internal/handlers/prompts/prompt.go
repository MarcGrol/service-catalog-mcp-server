package prompts

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Prompt struct {
	Contract mcp.Prompt
	Handler  func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
}

func (h Prompt) Attrs() (mcp.Prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)) {
	return h.Contract, h.Handler
}
