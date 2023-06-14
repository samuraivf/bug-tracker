package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	mock_log "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log/mocks"
)

func Test_CreateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		taskData       *dto.CreateTaskDto
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
	}{
		{
			name: "Error",
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "IN PROGRESS",
			},
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta(
						"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
					),
				).WithArgs(
					taskData.Name,
					taskData.Description,
					taskData.TaskPriority,
					taskData.ProjectID,
					taskData.TaskType,
					sqlmock.AnyArg(),
				).WillReturnError(err)
				log.EXPECT().Error(err)

				return &TaskRepository{db: db, log: log}
			},
			expectedResult: 0,
			expectedError:  err,
		},
		{
			name: "OK",
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "IN PROGRESS",
			},
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta(
						"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
					),
				).WithArgs(
					taskData.Name,
					taskData.Description,
					taskData.TaskPriority,
					taskData.ProjectID,
					taskData.TaskType,
					sqlmock.AnyArg(),
				).WillReturnRows(rows)
				log.EXPECT().Infof("Create task: id = %d", uint64(1))

				return &TaskRepository{db: db, log: log}
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.taskData)
			res, err := repo.CreateTask(test.taskData)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}
