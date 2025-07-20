package handlers

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewServiceCatalogTeamPrompt returns the MCP prompt contract and handler for project planning.
func NewServiceCatalogTeamPrompt() server.ServerPrompt {
	return server.ServerPrompt{
		Prompt: mcp.NewPrompt(
			"service_catalog_team",
			mcp.WithPromptDescription("Help making sense of the service catalog from the perspective of a team"),
			mcp.WithArgument("team_id", mcp.RequiredArgument(), mcp.ArgumentDescription("What is your team?")),
		),
		Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			args := request.Params.Arguments
			teamID := args["team_id"]
			if teamID == "" {
				teamID = "partner"
			}

			promptText := fmt.Sprintf(
				`You are a senior engineer within team '%s'. 
You want to understand which modules your team owns and how these modules are related to the rest of the platform.

For each of these modules, you would be interested to understand the following:
- What is the module responsible for?
- How big and complicated is the module?
- Which web-api's does the module expose?
- Which web-api's does the module consume?
- Which databases does the module consume?
- Which jobs does the module run?
- To which critical flows does the module belong?

After this overview you want to zoom in and deep dive. Examples of this are:
- You might be interested to get more details on certain modules or web-api's
- Alternatively you want to find out which modules (and teams) are your consumers.
- As a database owner, you wan to know which other modules are consuming your database.



`,
				teamID)
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
