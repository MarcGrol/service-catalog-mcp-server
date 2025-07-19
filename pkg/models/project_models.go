package models

import (
	"time"
)

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
	ProjectName string     `json:"project_name"`
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}
