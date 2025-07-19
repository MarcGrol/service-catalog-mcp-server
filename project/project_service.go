package project

import (
	"context"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/handlers/prompts"
	"github.com/MarcGrol/learnmcp/internal/handlers/resources"
	"github.com/MarcGrol/learnmcp/internal/handlers/tools"
	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
)

type ProjectService struct {
	server *server.MCPServer
	store  mystore.Store[model.Project]
}

func New(s *server.MCPServer, store mystore.Store[model.Project]) *ProjectService {
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

func (p *ProjectService) setupTools() {
	// Project management tool
	p.server.AddTool(tools.NewListProjectTool(p.store).Attrs())
	p.server.AddTool(tools.NewCreateProjectTool(p.store).Attrs())
	p.server.AddTool(tools.NewListTaskTool(p.store).Attrs())
	p.server.AddTool(tools.NewCreateTaskTool(p.store).Attrs())

	// Search tool
	p.server.AddTool(tools.NewSearchContentTool(p.store).Attrs())

	// Analytics tool
	p.server.AddTool(tools.NewGenerateAnalyticsTool(p.store).Attrs())
}

func (p *ProjectService) setupResources() {
	// Project management resource
	p.server.AddResource(resources.NewProjectListResource(p.store).Attrs())
	p.server.AddResource(resources.NewTasksListResource(p.store).Attrs())

	// Project statistics resource
	p.server.AddResource(resources.NewStatsResource(p.store).Attrs())

	// Documentation resource
	p.server.AddResource(resources.NewDocsResource().Attrs())
}

func (p *ProjectService) setupPrompts() {
	// Project planning prompt
	p.server.AddPrompt(prompts.NewPlanningPrompt().Attrs())

	// Code review prompt
	p.server.AddPrompt(prompts.NewReviewPrompt().Attrs())

	// Sprint planning prompt
	p.server.AddPrompt(prompts.NewSprintPrompt().Attrs())
}

func (p *ProjectService) preprovision(c context.Context) error {
	// Start with one sample project

	name := "Sample Project"

	project := model.Project{
		Name:        name,
		Version:     "1.0.0",
		Description: "Initial demo project",
		Authors:     []string{"Demo User"},
		Dependencies: map[string]string{
			"golang": "1.21+",
		},
		CreatedAt: time.Now().AddDate(0, 0, -7), // 7 days ago
		Tasks: []model.TaskItem{
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
