package dto

type WorkOnTaskDto struct {
	TaskID    uint64 `json:"taskId" validate:"required"`
	ProjectID uint64 `json:"projectId" validate:"required"`
}
