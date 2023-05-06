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

func Test_CreateProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
		projectData    *dto.CreateProjectDto
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().CreateProject(projectData).Return(uint64(0), err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			expectedResult: 0,
			expectedError:  err,
			projectData: &dto.CreateProjectDto{
				Name: "name",
			},
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().CreateProject(projectData).Return(uint64(1), nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			expectedResult: 1,
			expectedError:  nil,
			projectData: &dto.CreateProjectDto{
				Name:        "name",
				Description: "description",
				AdminID:     1,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.projectData)
			user, err := service.CreateProject(test.projectData)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}
