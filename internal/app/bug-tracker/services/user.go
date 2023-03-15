package services

import "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"

type UserService struct {
	repo repository.User
}

func NewUser(repo repository.User) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetUserByEmail() {

}

func (s *UserService) GetUserById() {

}

func (s *UserService) CreateUser() {

}
