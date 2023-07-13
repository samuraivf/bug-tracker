package dto

import "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"

type TaskByIdWithAssignee struct {
	Task     *models.Task `json:"task"`
	Assignee *models.User `json:"assignee"`
}
