package project

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/MarcGrol/learnmcp/mystore"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ProjectService struct {
	server *server.MCPServer
	store  mystore.Store[ProjectConfig]
}

type ProjectConfig struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Authors      []string          `json:"authors"`
	Dependencies map[string]string `json:"dependencies"`
	CreatedAt    time.Time         `json:"created_at"`
	Tasks        []TaskItem        `json:"tasks"`
}

type TaskItem struct {
	ProjectName string     `projectName:"id"`
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

func New(s *server.MCPServer, store mystore.Store[ProjectConfig]) *ProjectService {
	return &ProjectService{
		server: s,
		store:  store,
	}
}

func (p *ProjectService) Initialize(ctx context.Context) error {
	err := p.preprovision(ctx)
	if err != nil {
		return err
	}

	// Add tools
	p.setupTools()

	// Add resources
	p.setupResources()

	// Add prompts
	p.setupPrompts()

	log.Println("Starting Advanced MCP Server with Tools, Resources, and Prompts...")

	return nil
}

func (p *ProjectService) preprovision(c context.Context) error {
	// Start with one sample project

	name := "Sample Project"

	project := ProjectConfig{
		Name:        name,
		Version:     "1.0.0",
		Description: "Initial demo project",
		Authors:     []string{"Demo User"},
		Dependencies: map[string]string{
			"golang": "1.21+",
		},
		CreatedAt: time.Now().AddDate(0, 0, -7), // 7 days ago
		Tasks: []TaskItem{
			{
				ID:          1,
				Title:       "Setup project structure",
				Description: "Initialize the basic project layout",
				Status:      "done",
				Priority:    "high",
				CreatedAt:   time.Now().AddDate(0, 0, -5), // 5 days ago
			},
			{
				ID:          2,
				Title:       "Develop the first feature",
				Description: "First feature creates ...",
				Status:      "open",
				Priority:    "low",
				CreatedAt:   time.Now().AddDate(0, 0, -1), // 1 days ago
			},
		},
	}

	return p.store.Put(c, name, project)
}

func (p *ProjectService) setupTools() {
	// Project management tool
	listProjectTool := mcp.NewTool(
		"list_projects",
		mcp.WithDescription("Lists all projects"),
	)
	p.server.AddTool(listProjectTool, p.listProjectHandler())

	// Project management tool
	createProjectTool := mcp.NewTool(
		"create_project",
		mcp.WithDescription("Create a new project configuration"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
		mcp.WithString("description", mcp.Required(), mcp.Description("Project description")),
		mcp.WithArray("authors", mcp.Description("List of project authors")),
	)
	p.server.AddTool(createProjectTool, p.createProjectHandler())

	// Task management tool
	listTaskTool := mcp.NewTool(
		"list_tasks",
		mcp.WithDescription("Lists all tasks or all tasks of a project"),
		mcp.WithString("project_name", mcp.Description("Project that we want to list the tasks of")),
	)
	p.server.AddTool(listTaskTool, p.listTaskHandler())

	// Task management tool
	createTaskTool := mcp.NewTool(
		"create_task",
		mcp.WithDescription("Create a new task"),
		mcp.WithString("project_name", mcp.Required(), mcp.Description("Project that this task must be added to")),
		mcp.WithString("title", mcp.Required(), mcp.Description("Task title")),
		mcp.WithString("description", mcp.Description("Task description")),
		mcp.WithString("priority", mcp.Description("Task priority: low, medium, high")),
		mcp.WithString("due_date", mcp.Description("Due date in YYYY-MM-DD format")),
	)
	p.server.AddTool(createTaskTool, p.createTaskHandler())

	// Search tool
	searchTool := mcp.NewTool(
		"search_content",
		mcp.WithDescription("Search for content in projects and tasks"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("type", mcp.Description("Content type to search: project, task, all")),
	)
	p.server.AddTool(searchTool, p.searchHandler())

	// Analytics tool
	analyticsTool := mcp.NewTool(
		"generate_analytics",
		mcp.WithDescription("Generate project analytics and reports"),
		mcp.WithString("report_type", mcp.Required(), mcp.Description("Type of report: summary, tasks, timeline")),
	)
	p.server.AddTool(analyticsTool, p.analyticsHandler())
}

func (p *ProjectService) setupResources() {
	// Project configuration resource
	projectResource := mcp.NewResource(
		"project://config",
		"Current project configuration",
		mcp.WithMIMEType("application/json"),
	)
	p.server.AddResource(projectResource, p.projectResourceHandler())

	// Tasks list resource
	tasksResource := mcp.NewResource(
		"tasks://list",
		"List of all tasks in the project",
		mcp.WithMIMEType("application/json"),
	)
	p.server.AddResource(tasksResource, p.tasksResourceHandler())

	// Project statistics resource
	statsResource := mcp.NewResource(
		"stats://project",
		"Project statistics and metrics",
		mcp.WithMIMEType("application/json"),
	)
	p.server.AddResource(statsResource, p.statsResourceHandler())

	// Documentation resource
	docsResource := mcp.NewResource(
		"docs://readme",
		"Project documentation and README",
		mcp.WithMIMEType("text/markdown"),
	)
	p.server.AddResource(docsResource, p.docsResourceHandler())
}

func (p *ProjectService) setupPrompts() {
	// Project planning prompt
	planningPrompt := mcp.NewPrompt(
		"project_planning",
		mcp.WithPromptDescription("Help plan and structure a new project"),
		mcp.WithArgument("project_type", mcp.RequiredArgument(), mcp.ArgumentDescription("Type of project: web, mobile, api, library")),
		mcp.WithArgument("timeline", mcp.ArgumentDescription("Project timeline in weeks")),
		mcp.WithArgument("technologies", mcp.ArgumentDescription("Preferred technologies")),
	)
	p.server.AddPrompt(planningPrompt, p.planningPromptHandler())

	// Code review prompt
	reviewPrompt := mcp.NewPrompt(
		"code_review",
		mcp.WithPromptDescription("Generate code review guidelines and checklist"),
		mcp.WithArgument("language", mcp.RequiredArgument(), mcp.ArgumentDescription("Programming language")),
		mcp.WithArgument("focus", mcp.ArgumentDescription("Review focus: security, performance, style, all")),
	)
	p.server.AddPrompt(reviewPrompt, p.reviewPromptHandler())

	// Sprint planning prompt
	sprintPrompt := mcp.NewPrompt(
		"sprint_planning",
		mcp.WithPromptDescription("Assist with agile sprint planning and task breakdown"),
		mcp.WithArgument("sprint_length", mcp.ArgumentDescription("Sprint length in days")),
		mcp.WithArgument("team_size", mcp.ArgumentDescription("Number of team members")),
	)
	p.server.AddPrompt(sprintPrompt, p.sprintPromptHandler())
}

// Tool handlers
func (p *ProjectService) listProjectHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects, err := p.store.List(ctx)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error listing projects", err), err
		}

		// Simulate search results
		results := []string{}
		for _, p := range projects {
			results = append(results, fmt.Sprintf("%s: %s", p.Name, p.Description))
		}

		result := fmt.Sprintf("Currently available project:\n\n%s",
			strings.Join(results, "\n"))

		return mcp.NewToolResultText(result), nil
	}
}

func (p *ProjectService) createProjectHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("Missing project name"), nil
		}

		description, err := request.RequireString("description")
		if err != nil {
			return mcp.NewToolResultError("Missing project description"), nil
		}

		authors := request.GetStringSlice("authors", []string{"Anonymous Developer"})

		project := ProjectConfig{
			Name:        name,
			Version:     "1.0.0",
			Description: description,
			Authors:     authors,
			Dependencies: map[string]string{
				"golang": "1.21+",
			},
			CreatedAt: time.Now(),
		}

		// Save project
		err = p.store.Put(ctx, name, project)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error storing project", err), nil
		}

		projectJSON, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error serializing project", err), nil
		}

		result := fmt.Sprintf("Project '%s' created successfully!\n\nConfiguration:\n%s",
			name, string(projectJSON))

		return mcp.NewToolResultText(result), nil
	}
}

// Tool handlers
func (p *ProjectService) listTaskHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName := request.GetString("project_name", "")
		if projectName == "" {
			projects, err := p.store.List(ctx)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("Error listing projects", err), err
			}
			tasks := []string{}
			for _, p := range projects {
				for _, t := range p.Tasks {
					tasks = append(tasks, fmt.Sprintf("%s: %s - %s - %s", p.Name, t.ID, t.Title, t.Description))
				}
			}

			result := fmt.Sprintf("Currently available tasks:\n\n%s", strings.Join(tasks, "\n"))

			return mcp.NewToolResultText(result), nil
		}

		p, exists, err := p.store.Get(ctx, projectName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return mcp.NewToolResultErrorFromErr(fmt.Sprintf("project %s not found", projectName), err), nil
		}

		// Simulate search results
		results := []string{}
		for _, t := range p.Tasks {
			results = append(results, fmt.Sprintf("%s: %s - %s - %s", t.ID, projectName, t.Title, t.Description))
		}

		result := fmt.Sprintf("Currently available tasks within project %s:\n\n%s", projectName,
			strings.Join(results, "\n"))

		return mcp.NewToolResultText(result), nil
	}
}

func (p *ProjectService) createTaskHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectName, err := request.RequireString("project_name")
		if err != nil {
			return mcp.NewToolResultError("Missing project_name"), nil
		}

		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError("Missing task title"), nil
		}

		description := request.GetString("description", "")
		priority := request.GetString("priority", "medium")

		task := TaskItem{
			ProjectName: projectName,
			ID:          int(time.Now().Unix()),
			Title:       title,
			Description: description,
			Status:      "todo",
			Priority:    priority,
			CreatedAt:   time.Now(),
		}

		// Parse due date if provided
		dueDateStr := request.GetString("due_date", "")
		if dueDateStr != "" {
			if dueDate, err := time.Parse("2006-01-02", dueDateStr); err == nil {
				task.DueDate = &dueDate
			}
		}

		// Search project
		project, exists, err := p.store.Get(ctx, projectName)
		if err != nil {
			return nil, err
		}
		if !exists {
			return mcp.NewToolResultError(fmt.Sprintf("project %s not found", projectName)), nil
		}

		project.Tasks = append(project.Tasks, task)

		// Save project
		err = p.store.Put(ctx, projectName, project)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error storing task", err), nil
		}

		taskJSON, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Error serializing task", err), nil
		}
		result := fmt.Sprintf("Task created successfully!\n\n%s", string(taskJSON))

		return mcp.NewToolResultText(result), nil
	}
}

func (p *ProjectService) searchHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("Missing search query"), nil
		}

		searchType := request.GetString("type", "all")

		// Simulate search results
		results := []string{
			fmt.Sprintf("Found in project config: %s", strings.ToLower(query)),
			fmt.Sprintf("Found in task #123: %s related item", query),
			fmt.Sprintf("Found in documentation: %s reference", query),
		}

		result := fmt.Sprintf("Search Results for '%s' (type: %s):\n\n%s",
			query, searchType, strings.Join(results, "\n"))

		return mcp.NewToolResultText(result), nil
	}
}

func (p *ProjectService) analyticsHandler() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		reportType, err := request.RequireString("report_type")
		if err != nil {
			return mcp.NewToolResultError("Missing report type"), nil
		}

		var report string
		switch reportType {
		case "summary":
			report = `Project Summary Report
	========================
- Total Projects: 5
- Active Tasks: 12
- Completed Tasks: 8
- Team Members: 4
- Sprint Progress: 75%
- Code Coverage: 85%`

		case "tasks":
			report = `Task Analysis Report
===================
High Priority: 3 tasks
Medium Priority: 7 tasks
Low Priority: 2 tasks

Status Distribution:
- Todo: 5 tasks
- In Progress: 4 tasks
- In Review: 2 tasks
- Done: 1 task`

		case "timeline":
			report = `Timeline Report
==============
Week 1: Project setup and initial planning
Week 2: Core feature development (75% complete)
Week 3: Testing and refinement (planned)
Week 4: Documentation and deployment (planned)

Milestones:
✓ Project kickoff
✓ Architecture design
⧗ MVP completion (in progress)
◯ Beta release (upcoming)`

		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unknown report type: %s", reportType)), nil
		}

		return mcp.NewToolResultText(report), nil
	}
}

// Resource handlers
func (p *ProjectService) projectResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		projects, err := p.store.List(ctx)
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
	}
}

func (p *ProjectService) tasksResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		tasks := []TaskItem{}

		projects, err := p.store.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing projects: %s", err)
		}

		for _, p := range projects {
			tasks = append(tasks, p.Tasks...)
		}

		tasksJSON, err := json.MarshalIndent(map[string]interface{}{
			"total_tasks":  len(tasks),
			"tasks":        tasks,
			"last_updated": time.Now().Format(time.RFC3339),
		}, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error serializing results: %s", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(tasksJSON),
			},
		}, nil
	}
}

func (p *ProjectService) statsResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		tasks := []TaskItem{}

		projects, err := p.store.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("error listing projects: %s", err)
		}

		for _, p := range projects {
			tasks = append(tasks, p.Tasks...)
		}

		// CALCULATE from current data
		todoTasks := 0
		inProgressTasks := 0
		doneTasks := 0
		highPriorityTasks := 0

		for _, task := range tasks {
			switch task.Status {
			case "todo":
				todoTasks++
			case "in_progress":
				inProgressTasks++
			case "done":
				doneTasks++
			}

			if task.Priority == "high" {
				highPriorityTasks++
			}
		}

		stats := map[string]interface{}{
			"total_projects":      len(projects),
			"total_tasks":         len(tasks),
			"tasks_todo":          todoTasks,
			"tasks_in_progress":   inProgressTasks,
			"tasks_done":          doneTasks,
			"high_priority_tasks": highPriorityTasks,
			"completion_rate":     float64(doneTasks) / float64(len(tasks)) * 100,
			"last_calculated":     time.Now().Format(time.RFC3339),
		}

		statsJSON, err := json.MarshalIndent(stats, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error serializing results: %s", err)
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(statsJSON),
			},
		}, nil
	}
}

func (p *ProjectService) docsResourceHandler() func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
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
	}
}

// Prompt handlers
func (p *ProjectService) planningPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		var args map[string]string = request.Params.Arguments

		projectType := "web" // default
		if pt, ok := args["project_type"]; ok {
			projectType = pt
		}

		timeline := "4" // default weeks
		if tl, ok := args["timeline"]; ok {
			timeline = tl
		}

		// For comma-separated technologies
		technologiesStr := args["technologies"]
		technologies := []string{"Go", "React"}
		if technologiesStr != "" {
			technologies = strings.Split(technologiesStr, ",")
			// Trim whitespace
			for i, tech := range technologies {
				technologies[i] = strings.TrimSpace(tech)
			}
		}
		prompt := fmt.Sprintf(`You are a project planning assistant. Help plan a %s project.

Project Details:
- Type: %s
- Timeline: %s weeks
- Technologies: %s

Please provide:
1. Project structure and architecture recommendations
2. Key milestones and deliverables
3. Risk assessment and mitigation strategies
4. Resource allocation suggestions
5. Technology stack recommendations

Consider best practices for %s development and provide actionable advice.`,
			projectType, projectType, timeline, strings.Join(technologies, ", "), projectType)

		return &mcp.GetPromptResult{
			Description: "Code review guidance",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		}, nil
	}
}

func (p *ProjectService) reviewPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments

		language := args["language"]
		if language == "" {
			language = "Go"
		}

		focus := args["focus"]
		if focus == "" {
			focus = "all"
		}
		prompt := fmt.Sprintf(`You are a code review assistant for %s code.

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
						Text: prompt,
					},
				},
			},
		}, nil
	}
}

func (p *ProjectService) sprintPromptHandler() func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		args := request.Params.Arguments

		sprintLengthStr := args["sprint_length"]
		sprintLength := 14.0
		if sprintLengthStr != "" {
			if parsed, err := strconv.ParseFloat(sprintLengthStr, 64); err == nil {
				sprintLength = parsed
			}
		}

		teamSizeStr := args["team_size"]
		teamSize := 4.0
		if teamSizeStr != "" {
			if parsed, err := strconv.ParseFloat(teamSizeStr, 64); err == nil {
				teamSize = parsed
			}
		}

		prompt := fmt.Sprintf(`You are a sprint planning assistant for an agile team.

Sprint Configuration:
- Sprint Length: %.0f days
- Team Size: %.0f members

Help plan the upcoming sprint by:

1. **Capacity Planning**: Calculate team velocity and available story points
2. **Task Breakdown**: Break epics into manageable user stories
3. **Priority Assessment**: Rank tasks by business value and dependencies
4. **Risk Identification**: Identify potential blockers and mitigation strategies
5. **Timeline Creation**: Suggest daily standup goals and milestone checkpoints

Consider team capacity at 70%% for %.0f days = %.1f effective working days.
With %.0f team members, total capacity is approximately %.0f story points (assuming 1 point = 0.5 days).

Focus on creating realistic, achievable sprint goals that deliver value to stakeholders.`,
			sprintLength, teamSize, sprintLength, sprintLength*0.7, teamSize, teamSize*sprintLength*0.7*2)

		return &mcp.GetPromptResult{
			Description: "Sprint planning guidance",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Type: "text",
						Text: prompt,
					},
				},
			},
		}, nil
	}
}
