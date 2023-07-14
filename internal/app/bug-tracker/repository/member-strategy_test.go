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

func Test_new_memberStrategy(t *testing.T) {
	c := gomock.NewController(t)
	defer c.Finish()

	log := mock_log.NewMockLog(c)
	db, _, _ := sqlmock.New()

	expected := &memberStrategy{db, log}

	require.Equal(t, expected, new_memberStrategy(db, log))
}

func Test_IsMember(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *memberStrategy
	err := errors.New("error")

	tests := []struct {
		name          string
		projectID     uint64
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:      "Error cannot get member",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *memberStrategy {
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(projectID, userID).WillReturnError(err)
				log.EXPECT().Error(err).Return()

				return &memberStrategy{db, log}
			},
			expectedError: ErrNoRights,
		},
		{
			name:      "OK",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *memberStrategy {
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				rows := sqlmock.NewRows([]string{"member_id"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(projectID, userID).WillReturnRows(rows)

				return &memberStrategy{db, log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			member := test.mockBehaviour(c, test.projectID, test.userID)

			err := member.IsMember(test.projectID, test.userID)
			require.Equal(t, test.expectedError, err)
		})
	}
}
