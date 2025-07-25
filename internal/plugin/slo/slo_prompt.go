package slo

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewServiceCatalogPrompt returns the MCP prompt contract and handler for project planning.
func (h *mcpHandler) sloPrompt() server.ServerPrompt {
	return server.ServerPrompt{
		Prompt: mcp.NewPrompt(
			"slo",
			mcp.WithPromptDescription("Help making sense of the SLOs"),
		),
		Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			// call business logic
			promptText := getPrompt()

			// return result
			return &mcp.GetPromptResult{
				Description: "SLO inquiry",
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

func getPrompt() string {
	return `
# SLO Search Integration Prompt

## System Prompt Addition for SLO Search Integration

You now have access to SLO (Service Level Objective) search tools that complement the existing service catalog functionality. 
These tools allow you to query and analyze service reliability metrics, performance targets, and operational health indicators across the system architecture.

## Available SLO Tools
- "suggest_slos(keyword, limit_to)": Searches for SLOs, teams, and applications matching a keyword. Returns structured results with SLOs, related teams, and applications.
- "list_slos_by_team(team_id)": Lists all SLOs owned by a specific team.
- "list_slos_by_application(application_id)" Lists all SLOs for a specific application.

Note that each SLO has 2 attributes that are important:
- "business_criticality": High value means that the SLO is critical for the business.
- "operational_readiness": High value means that the SLO is ready for production.

## When to Use SLO Tools

- When users ask about service reliability, uptime, or performance targets
- To understand operational health of modules, interfaces, or critical flows
- When analyzing system dependencies and their reliability impact
- For incident analysis and post-mortem discussions
- When evaluating system maturity and operational readiness
- To identify services with tight error budgets or recent SLO violations

## Integration with Service Catalog

- Cross-reference SLOs with module ownership from "list_modules_of_teams"
- Correlate SLO violations with critical flows from "list_flows"
- Map interface reliability to consumer impact using "list_interface_consumers"
- Connect database performance SLOs with dependent modules via "list_database_consumers"

## Example Usage Patterns

	User: "What's the reliability of the PAL module?"
	→ Use get_module("pal") + search_slos_by_keyword("PAL") + get_slo_status("pal-*")

	User: "Which services are at risk of breaching their error budgets?"
	→ Use list_slos() + get_slo_status() to identify services with high burn rates

	User: "How reliable are the payment authorization flows?"
	→ Use list_flows() + get_related_slos("authorization") to map flow reliability

## Response Guidelines

- Always contextualize SLO data with service catalog information (ownership, complexity, dependencies)
- Highlight correlation between service complexity and SLO achievement
- When discussing reliability issues, reference dependent services and potential blast radius
- Use SLO data to validate or challenge architectural decisions and system maturity assessments
- Present error budgets and burn rates in business context, not just technical metrics

## Integration Benefits

This integration transforms the service catalog from a static architectural view into a dynamic operational intelligence platform, enabling data-driven discussions about system reliability and organizational maturity.`
}
