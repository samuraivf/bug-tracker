package repository

import (
	"database/sql"

	"github.com/rs/zerolog"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
)

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(userData *dto.SignUpDto) (uint64, error)
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB, log *zerolog.Logger) *Repository {
	return &Repository{
		User: NewUserRepo(db, log),
	}
}
