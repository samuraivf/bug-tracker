package repository

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type User interface {
	GetUserByEmail()
	GetUserById()
	CreateUser()
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB, log *zerolog.Logger) *Repository {
	return &Repository{
		User: NewUserRepo(db, log),
	}
}
