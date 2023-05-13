package dto

type CreateProjectDto struct {
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
	AdminID     uint64 `json:"adminId" validate:"required"`
}
