package services

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
)

func Test_CreateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		userID         uint64
		expectedResult uint64
		expectedError  error
		taskData       *dto.CreateTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().CreateTask(taskData, userID).Return(uint64(0), err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:         1,
			expectedResult: 0,
			expectedError:  err,
			taskData: &dto.CreateTaskDto{
				Name:        "name",
				Description: "description",
				ProjectID:   1,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().CreateTask(taskData, userID).Return(uint64(1), nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:         1,
			expectedResult: 1,
			expectedError:  nil,
			taskData: &dto.CreateTaskDto{
				Name:        "name",
				Description: "description",
				ProjectID:   1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.taskData, test.userID)
			user, err := service.CreateTask(test.taskData, test.userID)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_WorkOnTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		userID         uint64
		expectedError  error
		workOnTaskData *dto.WorkOnTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().WorkOnTask(workOnTaskData, userID).Return(err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: err,
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().WorkOnTask(workOnTaskData, userID).Return(nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: nil,
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.workOnTaskData, test.userID)
			err := service.WorkOnTask(test.workOnTaskData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_StopWorkOnTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		userID         uint64
		expectedError  error
		workOnTaskData *dto.WorkOnTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().StopWorkOnTask(workOnTaskData, userID).Return(err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: err,
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, workOnTaskData *dto.WorkOnTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().StopWorkOnTask(workOnTaskData, userID).Return(nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: nil,
			workOnTaskData: &dto.WorkOnTaskDto{
				TaskID:    1,
				ProjectID: 1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.workOnTaskData, test.userID)
			err := service.StopWorkOnTask(test.workOnTaskData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_UpdateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		userID         uint64
		expectedResult uint64
		expectedError  error
		taskData       *dto.UpdateTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().UpdateTask(taskData, userID).Return(uint64(0), err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:         1,
			expectedResult: 0,
			expectedError:  err,
			taskData: &dto.UpdateTaskDto{
				Name:        "name",
				Description: "description",
				ProjectID:   1,
				TaskID:      1,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.UpdateTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().UpdateTask(taskData, userID).Return(uint64(1), nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:         1,
			expectedResult: 1,
			expectedError:  nil,
			taskData: &dto.UpdateTaskDto{
				Name:        "name",
				Description: "description",
				ProjectID:   1,
				TaskID:      1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.taskData, test.userID)
			user, err := service.UpdateTask(test.taskData, test.userID)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetTaskById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		id             uint64
		expectedResult *models.Task
		expectedError  error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, id uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().GetTaskById(id).Return(nil, err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			id:             1,
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().GetTaskById(id).Return(&models.Task{ID: 1}, nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			id:             1,
			expectedResult: &models.Task{ID: 1},
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.id)
			user, err := service.GetTaskById(test.id)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetTasksByProjectId(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		id             uint64
		expectedResult []*models.Task
		expectedError  error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, id uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().GetTasksByProjectId(id).Return(nil, err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			id:             1,
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().GetTasksByProjectId(id).Return([]*models.Task{{ID: 1}}, nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			id:             1,
			expectedResult: []*models.Task{{ID: 1}},
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.id)
			user, err := service.GetTasksByProjectId(test.id)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *TaskService
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		userID        uint64
		expectedError error
		taskData      *dto.DeleteTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().DeleteTask(taskData, userID).Return(err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: err,
			taskData: &dto.DeleteTaskDto{
				ProjectID: 1,
				TaskID:    1,
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.DeleteTaskDto, userID uint64) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().DeleteTask(taskData, userID).Return(nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
			userID:        1,
			expectedError: nil,
			taskData: &dto.DeleteTaskDto{
				ProjectID: 1,
				TaskID:    1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.taskData, test.userID)
			err := service.DeleteTask(test.taskData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}
