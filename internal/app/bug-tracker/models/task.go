package models

import "time"

type Task struct {
	ID          uint64        `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	Priority    string        `json:"priority" db:"priority"`
	ProjectID   uint64        `json:"projectId" db:"project_id"`
	TaskType    string        `json:"taskType" db:"task_type"`
	PerformTo   time.Duration `json:"performTo" db:"perform_to"`
}
