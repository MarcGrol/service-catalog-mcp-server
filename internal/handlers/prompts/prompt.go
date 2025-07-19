package prompts

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

type Prompt struct {
	Prompt  mcp.Prompt
	Handler func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
}

func (h Prompt) Funcs() (mcp.Prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)) {
	return h.Prompt, h.Handler
}
