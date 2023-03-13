package dto

type SignInDto struct {
	Email    string `json:"email" validate:"reqired,email"`
	Password string `json:"password" validate:"reqired,min=8"`
}
