package repository

import (
	"database/sql"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(userData *dto.SignUpDto) (uint64, error)
}

type Repository struct {
	User
}

func NewRepository(db *sql.DB, log log.Log) *Repository {
	return &Repository{
		User: NewUserRepo(db, log),
	}
}
