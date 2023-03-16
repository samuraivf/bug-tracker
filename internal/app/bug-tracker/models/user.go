package models

type User struct {
	ID       uint64 `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
	Email    string `json:"email" db:"email"`
}
