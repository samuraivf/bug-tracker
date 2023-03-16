package repository

import (
	"database/sql"
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

var (
	ErrUserNotFound = errors.New("error user not found")
)

type UserRepository struct {
	db  *sql.DB
	log *zerolog.Logger
}

func NewUserRepo(db *sql.DB, log *zerolog.Logger) *UserRepository {
	return &UserRepository{db, log}
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user *models.User

	row := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error().Err(err)
			return nil, ErrUserNotFound
		}
		r.log.Error().Err(err)
		return nil, err
	}
	r.log.Info().Msgf("Get user with email: %s", email)

	return user, nil
}

func (r *UserRepository) GetUserById(id uint64) (*models.User, error) {
	var user *models.User

	row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error().Err(err)
			return nil, ErrUserNotFound
		}
		r.log.Error().Err(err)
		return nil, err
	}
	r.log.Info().Msgf("Get user with id: %d", id)

	return user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user *models.User

	row := r.db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	err := row.Scan(&user.ID, &user.Name, &user.Username, &user.Password, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error().Err(err)
			return nil, ErrUserNotFound
		}
		r.log.Error().Err(err)
		return nil, err
	}
	r.log.Info().Msgf("Get user with username: %s", username)

	return user, nil
}

func (r *UserRepository) CreateUser(userData *dto.SignUpDto) (uint64, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4)",
		userData.Name,
		userData.Username,
		userData.Email,
		userData.Password,
	)
	if err != nil {
		log.Error().Err(err)
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Error().Err(err)
		return 0, err
	}
	log.Info().Msgf("Create user: id = %d", userID)

	return uint64(userID), nil
}
