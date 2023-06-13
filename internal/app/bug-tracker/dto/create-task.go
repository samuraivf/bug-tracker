package dto

type CreateTaskDto struct {
	Name         string `json:"name" validate:"required,min=2"`
	Description  string `json:"description"`
	TaskPriority string `json:"taskPriority" validate:"required"`
	ProjectID    uint64 `json:"projectId" validate:"required"`
	TaskType     string `json:"taskType" validate:"required"`
}
