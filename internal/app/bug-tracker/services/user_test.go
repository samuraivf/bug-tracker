package services

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
)

func Test_GetUserByEmail(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, email string) *UserService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult *models.User
		expectedError  error
		email          string
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, email string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByEmail(email).Return(nil, err)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: nil,
			expectedError:  err,
			email:          "email@gmail.com",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, email string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByEmail(email).Return(&models.User{ID: 1}, nil)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: &models.User{ID: 1},
			expectedError:  nil,
			email:          "email@gmail.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.email)
			user, err := service.GetUserByEmail(test.email)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetUserById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *UserService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult *models.User
		expectedError  error
		id             uint64
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, id uint64) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserById(id).Return(nil, err)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: nil,
			expectedError:  err,
			id:             1,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserById(id).Return(&models.User{ID: 1}, nil)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: &models.User{ID: 1},
			expectedError:  nil,
			id:             1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.id)
			user, err := service.GetUserById(test.id)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetUserByUsername(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, username string) *UserService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult *models.User
		expectedError  error
		username       string
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, username string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByUsername(username).Return(nil, err)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: nil,
			expectedError:  err,
			username:       "username",
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, username string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByUsername(username).Return(&models.User{ID: 1}, nil)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: &models.User{ID: 1},
			expectedError:  nil,
			username:       "username",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.username)
			user, err := service.GetUserByUsername(test.username)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_CreateUser(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, userData *dto.SignUpDto) *UserService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
		userData       *dto.SignUpDto
	}{
		{
			name: "Error password too long",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *UserService {
				user := mock_repository.NewMockUser(c)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: 0,
			expectedError:  bcrypt.ErrPasswordTooLong,
			userData: &dto.SignUpDto{
				Password: "password11111111111111111111111111111111111111111111111111111111111111111111",
			},
		},
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().CreateUser(userData).Return(uint64(0), err)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: 0,
			expectedError:  err,
			userData: &dto.SignUpDto{
				Password: "password",
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, userData *dto.SignUpDto) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().CreateUser(userData).Return(uint64(1), nil)

				return &UserService{repo: repository.Repository{User: user}}
			},
			expectedResult: 1,
			expectedError:  nil,
			userData: &dto.SignUpDto{
				Password: "password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.userData)
			user, err := service.CreateUser(test.userData)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_generatePasswordHash(t *testing.T) {
	password := "password"
	hash, err := generatePasswordHash(password)

	require.True(t, len(hash) > 0)
	require.NoError(t, err)
	require.NoError(t, bcrypt.CompareHashAndPassword(hash, []byte(password)))
}

func Test_ValidateUser(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, email, password string) *UserService
	err := errors.New("error")
	hash, _ := bcrypt.GenerateFromPassword([]byte("password"), 2)
	userModel := &models.User{
		Password: string(hash),
	}

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		email          string
		password       string
		expectedResult *models.User
		expectedError  error
	}{
		{
			name: "Error in GetUserByEmail",
			mockBehaviour: func(c *gomock.Controller, email, password string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByEmail(email).Return(nil, err)

				return &UserService{repo: repository.Repository{User: user}}

			},
			email:          "email@gmail.com",
			password:       "password",
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "Error in comparing password",
			mockBehaviour: func(c *gomock.Controller, email, password string) *UserService {
				user := mock_repository.NewMockUser(c)
				hash, _ := bcrypt.GenerateFromPassword([]byte("password1"), 2)
				userModel := &models.User{
					Password: string(hash),
				}

				user.EXPECT().GetUserByEmail(email).Return(userModel, nil)

				return &UserService{repo: repository.Repository{User: user}}

			},
			email:          "email@gmail.com",
			password:       "password",
			expectedResult: nil,
			expectedError:  errInvalidPassword,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, email, password string) *UserService {
				user := mock_repository.NewMockUser(c)

				user.EXPECT().GetUserByEmail(email).Return(userModel, nil)

				return &UserService{repo: repository.Repository{User: user}}

			},
			email:          "email@gmail.com",
			password:       "password",
			expectedResult: userModel,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.email, test.password)
			user, err := service.ValidateUser(test.email, test.password)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}
