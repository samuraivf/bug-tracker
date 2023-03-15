package repository

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type UserRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewUserRepo(db *sql.DB, log *zerolog.Logger) *UserRepository {
	return &UserRepository{db, log}
}

func (r *UserRepository) GetUserByEmail() {

}

func (r *UserRepository) GetUserById() {

}

func (r *UserRepository) CreateUser() {

}
