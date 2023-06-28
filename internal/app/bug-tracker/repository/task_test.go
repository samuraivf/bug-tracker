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
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
)

func Test_CreateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskRepository
	err := errors.New("error")
	performTo := "2023-12-18 10:53:00"

	tests := []struct {
		name           string
		taskData       *dto.CreateTaskDto
		userID         uint64
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
	}{
		{
			name: "Error in admin",
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "IN PROGRESS",
				PerformTo:    performTo,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(taskData.ProjectID, userID).Return(err)

				return &TaskRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name: "Error",
			taskData: &dto.CreateTaskDto{
				Name:         "name",
				Description:  "description",
				TaskPriority: "high",
				ProjectID:    1,
				TaskType:     "IN PROGRESS",
				PerformTo:    performTo,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(taskData.ProjectID, userID).Return(nil)

				mock.ExpectQuery(
					regexp.QuoteMeta(
						"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at, perform_to) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
					),
				).WithArgs(
					taskData.Name,
					taskData.Description,
					taskData.TaskPriority,
					taskData.ProjectID,
					taskData.TaskType,
					sqlmock.AnyArg(),
					performTo,
				).WillReturnError(err)
				log.EXPECT().Error(err)

				return &TaskRepository{db: db, log: log, admin: admin}
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
				PerformTo:    performTo,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(taskData.ProjectID, userID).Return(nil)

				rows := sqlmock.NewRows([]string{"id"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta(
						"INSERT INTO tasks (name, description, task_priority, project_id, task_type, created_at, perform_to) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
					),
				).WithArgs(
					taskData.Name,
					taskData.Description,
					taskData.TaskPriority,
					taskData.ProjectID,
					taskData.TaskType,
					sqlmock.AnyArg(),
					performTo,
				).WillReturnRows(rows)
				log.EXPECT().Infof("Create task: id = %d", uint64(1))

				return &TaskRepository{db: db, log: log, admin: admin}
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.taskData, test.userID)
			res, err := repo.CreateTask(test.taskData, test.userID)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_WorkOnTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		workOnTaskData *dto.WorkOnTaskDto
		userID         uint64
		mockBehaviour  mockBehaviour
		expectedError  error
	}{
		{
			name: "Error no rights",
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskRepository {
				admin := mock_repository.NewMockadmin(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta(
						"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
					),
				).WithArgs(
					workOnTaskData.ProjectID,
					userID,
				).WillReturnError(err)
				admin.EXPECT().IsAdmin(workOnTaskData.ProjectID, userID).Return(err)

				return &TaskRepository{db: db, admin: admin}
			},
			expectedError: ErrNoRights,
		},
		{
			name: "Error cannot work on task",
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskRepository {
				admin := mock_repository.NewMockadmin(c)
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				rows := sqlmock.NewRows([]string{"member_id"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta(
						"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
					),
				).WithArgs(
					workOnTaskData.ProjectID,
					userID,
				).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta(
						"UPDATE tasks SET assignee = $1 WHERE id = $2 AND assignee IS NULL",
					),
				).WithArgs(
					userID,
					workOnTaskData.TaskID,
				).WillReturnError(err)
				log.EXPECT().Error(err)

				return &TaskRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name: "OK",
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
			userID: 1,
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskRepository {
				admin := mock_repository.NewMockadmin(c)
				db, mock, _ := sqlmock.New()
				log := mock_log.NewMockLog(c)

				rows := sqlmock.NewRows([]string{"member_id"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta(
						"SELECT member_id FROM projects_members WHERE project_id = $1 AND member_id = $2",
					),
				).WithArgs(
					workOnTaskData.ProjectID,
					userID,
				).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta(
						"UPDATE tasks SET assignee = $1 WHERE id = $2 AND assignee IS NULL",
					),
				).WithArgs(
					userID,
					workOnTaskData.TaskID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				return &TaskRepository{db: db, log: log, admin: admin}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.workOnTaskData, test.userID)
			err := repo.WorkOnTask(test.workOnTaskData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}
