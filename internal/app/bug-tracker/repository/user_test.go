package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/stretchr/testify/require"
)

func Test_GetUserByEmail(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, email string) *UserRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		email          string
		expectedResult *models.User
		expectedError  error
	}{
		{
			name:  "Error no rows",
			email: "email",
			mockBehaviour: func(c *gomock.Controller, email string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).WithArgs(email).WillReturnError(sql.ErrNoRows)
				log.EXPECT().Error(sql.ErrNoRows)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  ErrUserNotFound,
		},
		{
			name:  "Error",
			email: "email",
			mockBehaviour: func(c *gomock.Controller, email string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).WithArgs(email).WillReturnError(err)
				log.EXPECT().Error(err)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name:  "OK",
			email: "email",
			mockBehaviour: func(c *gomock.Controller, email string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "email"}).AddRow(uint64(1), "name", "username", "password", "email")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = $1")).WithArgs(email).WillReturnRows(rows)
				log.EXPECT().Infof("Get user with email: %s", email)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: &models.User{
				ID:       1,
				Name:     "name",
				Username: "username",
				Password: "password",
				Email:    "email",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.email)
			res, err := repo.GetUserByEmail(test.email)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetUserById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *UserRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		id             uint64
		expectedResult *models.User
		expectedError  error
	}{
		{
			name: "Error no rows",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1")).WithArgs(id).WillReturnError(sql.ErrNoRows)
				log.EXPECT().Error(sql.ErrNoRows)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  ErrUserNotFound,
		},
		{
			name: "Error",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1")).WithArgs(id).WillReturnError(err)
				log.EXPECT().Error(err)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "OK",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "email"}).AddRow(uint64(1), "name", "username", "password", "email")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE id = $1")).WithArgs(id).WillReturnRows(rows)
				log.EXPECT().Infof("Get user with id: %d", id)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: &models.User{
				ID:       1,
				Name:     "name",
				Username: "username",
				Password: "password",
				Email:    "email",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.id)
			res, err := repo.GetUserById(test.id)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetUserByUsername(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, username string) *UserRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		username       string
		expectedResult *models.User
		expectedError  error
	}{
		{
			name:     "Error no rows",
			username: "username",
			mockBehaviour: func(c *gomock.Controller, username string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = $1")).WithArgs(username).WillReturnError(sql.ErrNoRows)
				log.EXPECT().Error(sql.ErrNoRows)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  ErrUserNotFound,
		},
		{
			name:     "Error",
			username: "username",
			mockBehaviour: func(c *gomock.Controller, username string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = $1")).WithArgs(username).WillReturnError(err)
				log.EXPECT().Error(err)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name:     "OK",
			username: "username",
			mockBehaviour: func(c *gomock.Controller, username string) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "email"}).AddRow(uint64(1), "name", "username", "password", "email")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = $1")).WithArgs(username).WillReturnRows(rows)
				log.EXPECT().Infof("Get user with username: %s", username)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: &models.User{
				ID:       1,
				Name:     "name",
				Username: "username",
				Password: "password",
				Email:    "email",
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.username)
			res, err := repo.GetUserByUsername(test.username)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_CreateUser(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, data *dto.SignUpDto) *UserRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		data           *dto.SignUpDto
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
	}{
		{
			name: "Error",
			data: &dto.SignUpDto{
				Name: "name",
				Username: "username",
				Password: "password",
				Email: "email",
			},
			mockBehaviour: func(c *gomock.Controller, data *dto.SignUpDto) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id"),
				).WithArgs(data.Name, data.Username, data.Email, data.Password).WillReturnError(err)
				log.EXPECT().Error(err)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: 0,
			expectedError: err,
		},
		{
			name: "OK",
			data: &dto.SignUpDto{
				Name: "name",
				Username: "username",
				Password: "password",
				Email: "email",
			},
			mockBehaviour: func(c *gomock.Controller, data *dto.SignUpDto) *UserRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				userID := uint64(1)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(userID)
				mock.ExpectQuery(
					regexp.QuoteMeta("INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id"),
				).WithArgs(data.Name, data.Username, data.Email, data.Password).WillReturnRows(rows)
				log.EXPECT().Infof("Create user: id = %d", userID)

				return &UserRepository{db: db, log: log}
			},
			expectedResult: 1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.data)
			res, err := repo.CreateUser(test.data)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}
