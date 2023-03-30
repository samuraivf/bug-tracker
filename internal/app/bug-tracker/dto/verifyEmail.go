package dto

type VerifyEmail struct {
	Email string `json:"email" form:"email" validate:"required,email"`
}
