package models

import (
	"database/sql"
)

type Task struct {
	ID          uint64        `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Description string        `json:"description" db:"description"`
	Priority    string        `json:"priority" db:"priority"`
	ProjectID   uint64        `json:"projectId" db:"project_id"`
	TaskType    string        `json:"taskType" db:"task_type"`
	Assignee    sql.NullInt64 `json:"assignee" db:"assignee"`
	CreatedAt   sql.NullTime  `json:"createdAt" db:"created_at"`
	PerformTo   sql.NullTime  `json:"performTo" db:"perform_to"`
}
