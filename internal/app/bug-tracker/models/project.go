package models

type Project struct {
	ID          uint64 `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	AdminID     uint64 `json:"admin" db:"admin"`
}
