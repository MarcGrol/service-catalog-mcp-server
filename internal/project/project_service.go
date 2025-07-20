package project

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/server"

	"github.com/MarcGrol/learnmcp/internal/model"
	"github.com/MarcGrol/learnmcp/internal/mystore"
	"github.com/MarcGrol/learnmcp/internal/project/handlers/prompts"
	"github.com/MarcGrol/learnmcp/internal/project/handlers/resources"
	"github.com/MarcGrol/learnmcp/internal/project/handlers/tools"
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

	p.register()

	return nil
}

func (p *ProjectService) register() {
	p.server.AddTools(
		tools.NewListProjectTool(p.store),
		tools.NewCreateProjectTool(p.store),
		tools.NewListTaskTool(p.store),
		tools.NewCreateTaskTool(p.store),
		tools.NewSearchContentTool(p.store),
		tools.NewGenerateAnalyticsTool(p.store),
	)

	p.server.AddResources(
		resources.NewProjectListResource(p.store),
		resources.NewTasksListResource(p.store),
		resources.NewStatsResource(p.store),
		resources.NewDocsResource(),
	)

	p.server.AddPrompts(
		prompts.NewPlanningPrompt(),
		prompts.NewReviewPrompt(),
		prompts.NewSprintPrompt(),
	)
}

func (p *ProjectService) preprovision(c context.Context) error {
	// Start with one sample project as example

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
