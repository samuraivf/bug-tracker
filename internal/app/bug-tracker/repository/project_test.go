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
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

func Test_CreateProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		projectData    *dto.CreateProjectDto
		mockBehaviour  mockBehaviour
		expectedResult uint64
		expectedError  error
	}{
		{
			name: "Error",
			projectData: &dto.CreateProjectDto{
				Name:        "name",
				Description: "description",
				AdminID:     1,
			},
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta("INSERT INTO projects (name, description, admin) VALUES ($1, $2, $3) RETURNING id"),
				).WithArgs(projectData.Name, projectData.Description, projectData.AdminID).WillReturnError(err)
				log.EXPECT().Error(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedResult: 0,
			expectedError:  err,
		},
		{
			name: "OK",
			projectData: &dto.CreateProjectDto{
				Name:        "name",
				Description: "description",
				AdminID:     1,
			},
			mockBehaviour: func(c *gomock.Controller, projectData *dto.CreateProjectDto) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				projectID := uint64(1)
				rows := sqlmock.NewRows([]string{"id"}).AddRow(projectID)
				mock.ExpectQuery(
					regexp.QuoteMeta("INSERT INTO projects (name, description, admin) VALUES ($1, $2, $3) RETURNING id"),
				).WithArgs(projectData.Name, projectData.Description, projectData.AdminID).WillReturnRows(rows)
				log.EXPECT().Infof("Create project: id = %d", projectID)

				return &ProjectRepository{db: db, log: log}
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.projectData)
			res, err := repo.CreateProject(test.projectData)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetProjectById(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		id             uint64
		mockBehaviour  mockBehaviour
		expectedResult *models.Project
		expectedError  error
	}{
		{
			name: "Error",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM projects WHERE id = $1"),
				).WithArgs(id).WillReturnError(err)
				log.EXPECT().Error(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedResult: nil,
			expectedError:  err,
		},
		{
			name: "OK",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"id", "name", "description", "admin"}).AddRow(uint64(1), "name", "", uint64(1))

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM projects WHERE id = $1"),
				).WithArgs(id).WillReturnRows(rows)
				log.EXPECT().Infof("Get project: id = %d", uint64(1))

				return &ProjectRepository{db: db, log: log}
			},
			expectedResult: &models.Project{
				ID:          1,
				Name:        "name",
				Description: "",
				AdminID:     1,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.id)
			res, err := repo.GetProjectById(test.id)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		projectID     uint64
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:      "Error cannot get admin",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnError(err)
				log.EXPECT().Error(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: err,
		},
		{
			name:      "Error no rights",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(2))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnRows(rows)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: ErrNoRights,
		},
		{
			name:      "Error cannot delete project",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnError(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: err,
		},
		{
			name:      "OK",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.projectID, test.userID)
			err := repo.DeleteProject(test.projectID, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_UpdateProject(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		projectData   *dto.UpdateProjectDto
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:      "Error cannot get admin",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectData.ProjectID).WillReturnError(err)
				log.EXPECT().Error(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: err,
		},
		{
			name:      "Error no rights",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(2))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectData.ProjectID).WillReturnRows(rows)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: ErrNoRights,
		},
		{
			name:      "Error cannot update project",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectData.ProjectID).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET description = $1 WHERE id = $2"),
				).WithArgs(projectData.Description, projectData.ProjectID).WillReturnError(err)

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: err,
		},
		{
			name:      "OK",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				rows := sqlmock.NewRows([]string{"admin"}).AddRow(uint64(1))
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT admin FROM projects WHERE id = $1"),
				).WithArgs(projectData.ProjectID).WillReturnRows(rows)

				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET description = $1 WHERE id = $2"),
				).WithArgs(projectData.Description, projectData.ProjectID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: log}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.projectData, test.userID)
			err := repo.UpdateProject(test.projectData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}
