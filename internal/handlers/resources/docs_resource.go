package resources

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

// NewDocsResource returns the MCP resource contract and handler for documentation.
func NewDocsResource() Resource {
	return Resource{
		Contract: mcp.NewResource(
			"docs://readme",
			"Project documentation and README",
			mcp.WithMIMEType("text/markdown"),
		),
		Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			readme := `# Advanced Project Manager MCP Server

This is an advanced example of an MCP server built with Go that demonstrates:

## Features

- **Tools**: Project and task management
- **Resources**: Configuration, task lists, statistics
- **Prompts**: Planning and review assistance

## Usage

1. Create projects with "create_project"
2. Add tasks with "create_task"
3. Search content with "search_content"
4. Generate reports with "generate_analytics"

## Resources Available

- "project://config" - Current project configuration
- "tasks://list" - All project tasks
- "stats://project" - Project statistics
- "docs://readme" - This documentation

## Getting Started

Connect this server to your MCP client (Claude Desktop, Cursor, etc.) and start managing your projects with AI assistance!
`
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "text/markdown",
					Text:     readme,
				},
			}, nil
		},
	}
}
