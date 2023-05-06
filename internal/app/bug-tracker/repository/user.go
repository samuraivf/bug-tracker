package repository

import (
	"database/sql"
	"errors"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type UserRepository struct {
	db  *sql.DB
	log log.Log
}

func NewUserRepo(db *sql.DB, log log.Log) User {
	return &UserRepository{db, log}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)

	row := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error(err)
			return nil, ErrUserNotFound
		}
		r.log.Error(err)
		return nil, err
	}
	r.log.Infof("Get user with email: %s", email)

	return user, nil
}

func (r *UserRepository) GetUserById(id uint64) (*models.User, error) {
	user := new(models.User)

	row := r.db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error(err)
			return nil, ErrUserNotFound
		}
		r.log.Error(err)
		return nil, err
	}
	r.log.Infof("Get user with id: %d", id)

	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := new(models.User)

	row := r.db.QueryRow("SELECT * FROM users WHERE username = $1", username)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error(err)
			return nil, ErrUserNotFound
		}
		r.log.Error(err)
		return nil, err
	}
	r.log.Infof("Get user with username: %s", username)

	return user, nil
}

func (r *UserRepository) CreateUser(userData *dto.SignUpDto) (uint64, error) {
	result := r.db.QueryRow(
		"INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		userData.Name,
		userData.Username,
		userData.Email,
		userData.Password,
	)

	var userID uint64
	if err := result.Scan(&userID); err != nil {
		r.log.Error(err)
		return 0, err
	}
	r.log.Infof("Create user: id = %d", userID)

	return userID, nil
}
