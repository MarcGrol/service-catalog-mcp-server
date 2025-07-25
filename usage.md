# Available Tools Reference Guide

This document describes all available tools for managing SLOs, services, and system analysis.

## Project Server Tools (SLO/Service Management)

### SLO Management Tools

#### `suggest_slos(keyword, limit_to)`
Searches for SLOs, teams, and applications matching a keyword. Returns structured results with SLOs, related teams, and applications.

#### `list_slos_by_team(team_id)`
Lists all SLOs owned by a specific team.

#### `list_slos_by_application(application_id)`
Lists all SLOs for a specific application.

#### `get_slo_by_id(slo_id)`
Gets detailed information about a specific SLO including:
- Configuration details (target, duration, category)
- Business impact flags (critical, frontdoor, payment flows)
- Monitoring setup (alerts, dashboards, notifications)
- PromQL queries and metrics

### Module Management Tools

#### `suggest_candidates(keyword, limit_to)`
General search across modules, interfaces, databases, and teams. Broader than `suggest_slos`.

#### `list_modules(filter_keyword)`
Lists modules (services/components) filtered by keyword.

#### `list_modules_by_complexity(limit_to)`
Lists modules ordered by complexity (most complex first). Useful for identifying high-maintenance services.

#### `list_modules_of_teams(team_id)`
Lists all modules owned by a specific team.

#### `list_modules_with_kind(kind_id)`
Lists modules of a specific type/kind (e.g., "web-service", "database", "terminal-component").

#### `get_module(module_id)`
Gets detailed information about a specific module including dependencies, interfaces, and configuration.

### Interface Management Tools

#### `list_interfaces(filter_keyword)`
Lists web APIs and interfaces filtered by keyword.

#### `list_interfaces_by_complexity(limit_to)`
Lists interfaces ordered by complexity.

#### `list_interface_consumers(interface_id)`
Lists all modules that consume/depend on a specific interface. Useful for impact analysis.

#### `get_interface(interface_id)`
Gets detailed information about a specific interface/API including endpoints, consumers, and specifications.

### Database & Dependency Tools

#### `list_database_consumers(database_id)`
Lists all modules that use a specific database. Critical for understanding data dependencies.

### Flow Management Tools

#### `list_flows()`
Lists all critical business flows in the system (e.g., online payments, onboarding).

#### `list_flow_participants(flow_id)`
Lists all modules that participate in a specific business flow.

### Catalog Information Tools

#### `list_kinds()`
Lists all available module types/categories in the system.

## Common Usage Patterns

### SLO Analysis Workflow
1. **Discovery:** Use `suggest_slos(keyword)` to find relevant SLOs
2. **Scoping:** Use `list_slos_by_team()` or `list_slos_by_application()` for specific areas
3. **Details:** Use `get_slo_by_id()` for comprehensive SLO information
4. **Documentation:** Use `artifacts` to create reports or summaries

### Service Architecture Analysis
1. **Search:** Use `suggest_candidates(keyword)` for general exploration
2. **Structure:** Use `list_modules()` and `get_module()` to understand services
3. **Dependencies:** Use `list_interface_consumers()` and `list_database_consumers()`
4. **Impact:** Use `list_flow_participants()` to understand business impact

## Best Practices

### Search Strategy
- Start with broad searches using `suggest_` functions
- Narrow down with specific `list_` functions
- Get details with `get_` functions

### Performance Considerations
- `list_slos()` without filters can be very large - use team/application filters
- Use `limit_to` parameters to control result sizes
- `suggest_` functions are optimized for discovery

### Error Handling
- Tools will suggest correct parameter names if you use wrong ones
- Team/application IDs are case-sensitive
- Use `suggest_candidates()` to find correct identifiers

### Documentation
- Use `artifacts` for any content meant to be saved or referenced
- Include both technical details and business context
- Structure documentation with clear headings and sections
