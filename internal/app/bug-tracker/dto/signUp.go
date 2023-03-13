package dto

type SignUpDto struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}
