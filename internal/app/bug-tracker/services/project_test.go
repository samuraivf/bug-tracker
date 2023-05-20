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

func Test_GetProjectById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name           string
		mockBehaviour  mockBehaviour
		id             uint64
		expectedResult *models.Project
		expectedError  error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, id uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().GetProjectById(id).Return(nil, err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			id:             1,
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, id uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().GetProjectById(id).Return(&models.Project{ID: 1}, nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			id:             1,
			expectedResult: &models.Project{ID: 1},
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.id)
			user, err := service.GetProjectById(test.id)

			require.Equal(t, test.expectedResult, user)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		projectID     uint64
		userID        uint64
		expectedError error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().DeleteProject(projectID, userID).Return(err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectID:     1,
			userID:        1,
			expectedError: err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().DeleteProject(projectID, userID).Return(nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectID:     1,
			userID:        1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.projectID, test.userID)
			err := service.DeleteProject(test.projectID, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_UpdateProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		projectData   *dto.UpdateProjectDto
		userID        uint64
		expectedError error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().UpdateProject(projectData, userID).Return(err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectData: &dto.UpdateProjectDto{Description: "description"},
			userID:        1,
			expectedError: err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().UpdateProject(projectData, userID).Return(nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectData: &dto.UpdateProjectDto{Description: "description"},
			userID:        1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.projectData, test.userID)
			err := service.UpdateProject(test.projectData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_AddMember(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		memberData   *dto.AddMemberDto
		userID        uint64
		expectedError error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().AddMember(memberData, userID).Return(err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:        1,
			expectedError: err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().AddMember(memberData, userID).Return(nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:        1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.memberData, test.userID)
			err := service.AddMember(test.memberData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_LeaveProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *ProjectService
	err := errors.New("error")

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		projectID     uint64
		userID        uint64
		expectedError error
	}{
		{
			name: "Error",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().LeaveProject(projectID, userID).Return(err)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectID:     1,
			userID:        1,
			expectedError: err,
		},
		{
			name: "OK",
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectService {
				project := mock_repository.NewMockProject(c)

				project.EXPECT().LeaveProject(projectID, userID).Return(nil)

				return &ProjectService{repo: repository.Repository{Project: project}}
			},
			projectID:     1,
			userID:        1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			service := test.mockBehaviour(c, test.projectID, test.userID)
			err := service.LeaveProject(test.projectID, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}