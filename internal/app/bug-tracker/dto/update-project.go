package dto

type UpdateProjectDto struct {
	ProjectID   uint64 `json:"projectId"`
	Description string `json:"description"`
}
