package dto

import "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"

type ProjectWithTasks struct {
	Project *models.Project `json:"project"`
	Tasks   []*models.Task  `json:"tasks"`
}
