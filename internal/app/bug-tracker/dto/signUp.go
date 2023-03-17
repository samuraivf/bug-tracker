package dto

type SignUpDto struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Email    string `json:"email" validate:"required,email"`
}
