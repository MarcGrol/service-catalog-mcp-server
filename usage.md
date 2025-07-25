# Available Tools Reference Guide

This document describes all available tools for managing SLOs, services, and system analysis.

## Project Server Tools (SLO/Service Management)

### SLO Management Tools

#### `list_slos()`
Lists all SLOs in the system. **Warning:** Can return very large datasets that may exceed display limits.

#### `get_slo_by_id(slo_id)`
Gets detailed information about a specific SLO including:
- Configuration details (target, duration, category)
- Business impact flags (critical, frontdoor, payment flows)
- Monitoring setup (alerts, dashboards, notifications)
- PromQL queries and metrics

**Parameters:**
- `slo_id` - Unique identifier for the SLO

#### `list_slos_by_team(team_id)`
Lists all SLOs owned by a specific team.

**Parameters:**
- `team_id` - Team identifier (e.g., "access-management", "checkout")

#### `list_slos_by_application(application_id)`
Lists all SLOs for a specific application.

**Parameters:**
- `application_id` - Application identifier (e.g., "partner", "checkout")

#### `suggest_slos(keyword, limit_to)`
Searches for SLOs, teams, and applications matching a keyword. Returns structured results with SLOs, related teams, and applications.

**Parameters:**
- `keyword` - Search term
- `limit_to` - Maximum results per category (optional)

### Module Management Tools

#### `list_modules(filter_keyword)`
Lists modules (services/components) filtered by keyword.

**Parameters:**
- `filter_keyword` - Keyword to filter modules

#### `get_module(module_id)`
Gets detailed information about a specific module including dependencies, interfaces, and configuration.

**Parameters:**
- `module_id` - Module identifier

#### `list_modules_by_complexity(limit_to)`
Lists modules ordered by complexity (most complex first). Useful for identifying high-maintenance services.

**Parameters:**
- `limit_to` - Maximum number of modules to return

#### `list_modules_of_teams(team_id)`
Lists all modules owned by a specific team.

**Parameters:**
- `team_id` - Team identifier

#### `list_modules_with_kind(kind_id)`
Lists modules of a specific type/kind (e.g., "web-service", "database", "terminal-component").

**Parameters:**
- `kind_id` - Module type identifier

### Interface Management Tools

#### `list_interfaces(filter_keyword)`
Lists web APIs and interfaces filtered by keyword.

**Parameters:**
- `filter_keyword` - Keyword to filter interfaces

#### `get_interface(interface_id)`
Gets detailed information about a specific interface/API including endpoints, consumers, and specifications.

**Parameters:**
- `interface_id` - Interface identifier

#### `list_interfaces_by_complexity(limit_to)`
Lists interfaces ordered by complexity.

**Parameters:**
- `limit_to` - Maximum number of interfaces to return

#### `list_interface_consumers(interface_id)`
Lists all modules that consume/depend on a specific interface. Useful for impact analysis.

**Parameters:**
- `interface_id` - Interface identifier

### Database & Dependency Tools

#### `list_database_consumers(database_id)`
Lists all modules that use a specific database. Critical for understanding data dependencies.

**Parameters:**
- `database_id` - Database identifier

### Flow Management Tools

#### `list_flows()`
Lists all critical business flows in the system (e.g., online payments, onboarding).

#### `list_flow_participants(flow_id)`
Lists all modules that participate in a specific business flow.

**Parameters:**
- `flow_id` - Flow identifier

### Catalog Information Tools

#### `list_kinds()`
Lists all available module types/categories in the system.

#### `suggest_candidates(keyword, limit_to)`
General search across modules, interfaces, databases, and teams. Broader than `suggest_slos`.

**Parameters:**
- `keyword` - Search term
- `limit_to` - Maximum results per category

## General Purpose Tools

### Content Creation

#### `artifacts`
Creates and updates structured content including:
- Code snippets and applications
- Technical documentation
- Data visualizations
- Reports and presentations
- Markdown documents

**Common Commands:**
- `create` - Create new artifact
- `update` - Modify existing artifact
- `rewrite` - Complete rewrite of artifact

### Data Analysis

#### `repl`
JavaScript analysis tool for:
- Complex mathematical calculations
- File processing (CSV, Excel, JSON)
- Data manipulation and analysis
- Statistical computations

**Use Cases:**
- Processing uploaded data files
- Complex calculations requiring high precision
- Data transformations and aggregations
- File format conversions

### Web Research

#### `web_search(query)`
Searches the web for current information. Use for:
- Recent developments and news
- Real-time data and metrics
- Current best practices
- Verification of information

**Parameters:**
- `query` - Search terms (keep concise, 1-6 words work best)

#### `web_fetch(url)`
Fetches complete content from specific web pages. Use after `web_search` to get full articles or documentation.

**Parameters:**
- `url` - Exact URL to fetch (must be provided by user or from search results)

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

### Data Analysis Projects
1. **Exploration:** Use `repl` for initial data inspection
2. **Processing:** Use `repl` for complex transformations
3. **Visualization:** Use `artifacts` to create charts and reports
4. **Documentation:** Use `artifacts` for final documentation

### Research and Documentation
1. **Current Info:** Use `web_search` for recent developments
2. **Deep Dive:** Use `web_fetch` for complete content
3. **Analysis:** Use `repl` for data processing if needed
4. **Output:** Use `artifacts` for final documentation

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