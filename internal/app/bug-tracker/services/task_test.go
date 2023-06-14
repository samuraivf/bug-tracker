package services

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
)

func Test_CreateTask(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
		taskData       *dto.CreateTaskDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().CreateTask(taskData).Return(uint64(0), err)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
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
			mockBehaviour: func(c *gomock.Controller, taskData *dto.CreateTaskDto) *TaskService {
				task := mock_repository.NewMockTask(c)

				task.EXPECT().CreateTask(taskData).Return(uint64(1), nil)

				return &TaskService{repo: repository.Repository{Task: task}}
			},
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

			service := test.mockBehaviour(c, test.taskData)
			user, err := service.CreateTask(test.taskData)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}
