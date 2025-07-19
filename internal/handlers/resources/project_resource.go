package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/mark3labs/mcp-go/mcp"
)

// NewProjectConfigResourceAndHandler returns the MCP resource contract and handler for project configuration.
func NewProjectListResource(store mystore.Store[model.Project]) Resource {
	return Resource{
		Contract: mcp.NewResource(
			"project://config",
			"Current project configuration",
			mcp.WithMIMEType("application/json"),
		),
		Handler: func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			projects, err := store.List(ctx)
			if err != nil {
				return nil, fmt.Errorf("error listing projects: %s", err)
			}
			projectsJson, err := json.MarshalIndent(
				map[string]interface{}{
					"total_projects": len(projects),
					"projects":       projects,
					"last_updated":   time.Now().Format(time.RFC3339),
				}, "", "  ")
			if err != nil {
				return nil, err
			}
			return []mcp.ResourceContents{
				mcp.TextResourceContents{
					URI:      request.Params.URI,
					MIMEType: "application/json",
					Text:     string(projectsJson),
				},
			}, nil
		},
	}
}
