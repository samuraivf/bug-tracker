package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
	"github.com/stretchr/testify/require"
)

func Test_new_adminStrategy(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	log := mock_log.NewMockLog(c)
	db, _, _ := sqlmock.New()

	expected := &adminStrategy{db, log}

	require.Equal(t, expected, new_adminStrategy(db, log))
}

func Test_IsAdmin(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *adminStrategy
	err := errors.New("error")

	tests := []struct {
		name          string
		projectID     uint64
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name: "Error cannot get admin",
			projectID: 1,
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *adminStrategy {
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnError(err)
				log.EXPECT().Error(err).Return()

				return &adminStrategy{db, log}
			},
			expectedError: err,
		},
		{
			name: "Error no rights",
			projectID: 1,
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *adminStrategy {
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(2))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnRows(rows)

				return &adminStrategy{db, log}
			},
			expectedError: ErrNoRights,
		},
		{
			name: "OK",
			projectID: 1,
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *adminStrategy {
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnRows(rows)

				return &adminStrategy{db, log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			admin := test.mockBehaviour(c, test.projectID, test.userID)

			err := admin.IsAdmin(test.projectID, test.userID)
			require.Equal(t, test.expectedError, err)
		})
	}
}
