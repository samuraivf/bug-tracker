package services

import (
	"errors"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidPassword = errors.New("error invalid password")
)

type UserService struct {
	repo repository.User
}

func NewUser(repo repository.User) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserById(id uint64) (*models.User, error) {
	return s.repo.GetUserById(id)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *UserService) CreateUser(userData *dto.SignUpDto) (uint64, error) {
	passwordHash, err := generatePasswordHash(userData.Password)
	if err != nil {
		return 0, err
	}

	userData.Password = string(passwordHash)

	return s.repo.CreateUser(userData)
}

func generatePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func (s *UserService) ValidateUser(email, password string) (*models.User, error) {
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errInvalidPassword
	}

	return user, nil
}
