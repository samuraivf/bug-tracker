package dto

type UpdateTaskDto struct {
	Name         string `json:"name" validate:"required,min=2"`
	Description  string `json:"description" validate:"required"`
	TaskPriority string `json:"taskPriority" validate:"required"`
	TaskID       uint64 `json:"taskId" validate:"required"`
	ProjectID    uint64 `json:"projectId" validate:"required"`
	TaskType     string `json:"taskType" validate:"required"`
	PerformTo    string `json:"performTo"`
}
