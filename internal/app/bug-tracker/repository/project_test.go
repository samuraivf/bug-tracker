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
	mock_repository "github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository/mocks"
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
			name:      "Error in admin",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(err)

				return &ProjectRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name:      "Error cannot delete project",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)
				log := mock_log.NewMockLog(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:      "OK",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects WHERE id = $1"),
				).WithArgs(projectID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: nil, admin: admin}
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
			name:        "Error in admin",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:      1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectData.ProjectID, userID).Return(err)

				return &ProjectRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name:        "Error cannot update project",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:      1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)
				log := mock_log.NewMockLog(c)

				admin.EXPECT().IsAdmin(projectData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET description = $1 WHERE id = $2"),
				).WithArgs(projectData.Description, projectData.ProjectID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:        "OK",
			projectData: &dto.UpdateProjectDto{ProjectID: 1},
			userID:      1,
			mockBehaviour: func(c *gomock.Controller, projectData *dto.UpdateProjectDto, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET description = $1 WHERE id = $2"),
				).WithArgs(projectData.Description, projectData.ProjectID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: nil, admin: admin}
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

func Test_AddMember(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		memberData    *dto.AddMemberDto
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:       "Error in admin",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(err)

				return &ProjectRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name:       "Error cannot add member to project",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)"),
				).WithArgs(memberData.ProjectID, memberData.MemberID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:       "OK",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)"),
				).WithArgs(memberData.ProjectID, memberData.MemberID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.memberData, test.userID)
			err := repo.AddMember(test.memberData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_DeleteMember(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		memberData    *dto.AddMemberDto
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:       "Error in admin",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(err)

				return &ProjectRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name:       "Error cannot delete member from project",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(memberData.ProjectID, memberData.MemberID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:       "OK",
			memberData: &dto.AddMemberDto{ProjectID: 1, MemberID: 2},
			userID:     1,
			mockBehaviour: func(c *gomock.Controller, memberData *dto.AddMemberDto, userID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(memberData.ProjectID, userID).Return(nil)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(memberData.ProjectID, memberData.MemberID).WillReturnResult(sqlmock.NewResult(1, 1))

				log.EXPECT().Infof("Delete member with id=%d from project with id=%d", memberData.MemberID, memberData.ProjectID)

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.memberData, test.userID)
			err := repo.DeleteMember(test.memberData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_GetMembers(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		projectID      uint64
		userID         uint64
		mockBehaviour  mockBehaviour
		expectedError  error
		expectedResult []*models.User
	}{
		{
			name:      "Error in admin",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)
				member := mock_repository.NewMockmember(c)

				member.EXPECT().IsMember(projectID, userID).Return(ErrNoRights)
				admin.EXPECT().IsAdmin(projectID, userID).Return(ErrNoRights)

				return &ProjectRepository{admin: admin, member: member}
			},
			expectedError:  ErrNoRights,
			expectedResult: nil,
		},
		{
			name:      "Error cannot get members",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				member := mock_repository.NewMockmember(c)
				log := mock_log.NewMockLog(c)

				member.EXPECT().IsMember(projectID, userID).Return(nil)

				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM users WHERE users.id IN (SELECT member_id FROM projects_members WHERE projects_members.project_id = $1)"),
				).WithArgs(projectID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, member: member}
			},
			expectedError:  err,
			expectedResult: nil,
		},
		{
			name:      "OK",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				member := mock_repository.NewMockmember(c)

				member.EXPECT().IsMember(projectID, userID).Return(nil)

				rows := sqlmock.NewRows([]string{"id", "name", "username", "password", "email"}).AddRow(uint64(1), "name1", "username1", "password1", "email1")
				mock.ExpectQuery(
					regexp.QuoteMeta("SELECT * FROM users WHERE users.id IN (SELECT member_id FROM projects_members WHERE projects_members.project_id = $1)"),
				).WithArgs(projectID).WillReturnRows(rows)

				return &ProjectRepository{db: db, log: nil, member: member}
			},
			expectedError: nil,
			expectedResult: []*models.User{
				{
					ID:       1,
					Name:     "name1",
					Username: "username1",
					Password: "password1",
					Email:    "email1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.projectID, test.userID)
			res, err := repo.GetMembers(test.projectID, test.userID)

			require.Equal(t, test.expectedError, err)
			require.Equal(t, test.expectedResult, res)
		})
	}
}

func Test_LeaveProject(t *testing.T) {
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
			name:      "Error no error in admin",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(nil)

				return &ProjectRepository{admin: admin}
			},
			expectedError: ErrNoRights,
		},
		{
			name:      "Error cannot leave project",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)
				log := mock_log.NewMockLog(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(err)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(projectID, userID).WillReturnError(err)

				log.EXPECT().Error(err).Return()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:      "OK",
			projectID: 1,
			userID:    1,
			mockBehaviour: func(c *gomock.Controller, projectID, userID uint64) *ProjectRepository {
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(projectID, userID).Return(err)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(projectID, userID).WillReturnResult(sqlmock.NewResult(1, 1))

				return &ProjectRepository{db: db, log: nil, admin: admin}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.projectID, test.userID)
			err := repo.LeaveProject(test.projectID, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_SetNewAdmin(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name          string
		newAdminData  *dto.NewAdminDto
		userID        uint64
		mockBehaviour mockBehaviour
		expectedError error
	}{
		{
			name:         "Error in admin",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(err)

				return &ProjectRepository{admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "Error cannot begin transaction",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin().WillReturnError(err)

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "Error cannot set new admin",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET admin = $1 WHERE id = $2"),
				).WithArgs(
					newAdminData.NewAdminID,
					newAdminData.ProjectID,
				).WillReturnError(err)

				log.EXPECT().Error(err)
				mock.ExpectRollback()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "Error cannot delete from projects_members",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET admin = $1 WHERE id = $2"),
				).WithArgs(
					newAdminData.NewAdminID,
					newAdminData.ProjectID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				log.EXPECT().Infof("Set new admin = %d in project = %d", newAdminData.NewAdminID, newAdminData.ProjectID)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(
					newAdminData.ProjectID,
					newAdminData.NewAdminID,
				).WillReturnError(err)

				log.EXPECT().Error(err)
				mock.ExpectRollback()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "Error cannot insert into projects_members",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET admin = $1 WHERE id = $2"),
				).WithArgs(
					newAdminData.NewAdminID,
					newAdminData.ProjectID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				log.EXPECT().Infof("Set new admin = %d in project = %d", newAdminData.NewAdminID, newAdminData.ProjectID)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(
					newAdminData.ProjectID,
					newAdminData.NewAdminID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)"),
				).WithArgs(
					newAdminData.ProjectID,
					adminID,
				).WillReturnError(err)

				log.EXPECT().Error(err)
				mock.ExpectRollback()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "Error cannot commit transaction",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET admin = $1 WHERE id = $2"),
				).WithArgs(
					newAdminData.NewAdminID,
					newAdminData.ProjectID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				log.EXPECT().Infof("Set new admin = %d in project = %d", newAdminData.NewAdminID, newAdminData.ProjectID)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(
					newAdminData.ProjectID,
					newAdminData.NewAdminID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)"),
				).WithArgs(
					newAdminData.ProjectID,
					adminID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(err)
				log.EXPECT().Error(err)

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: err,
		},
		{
			name:         "OK",
			newAdminData: &dto.NewAdminDto{ProjectID: 1, NewAdminID: 2},
			userID:       1,
			mockBehaviour: func(c *gomock.Controller, newAdminData *dto.NewAdminDto, adminID uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()
				admin := mock_repository.NewMockadmin(c)

				admin.EXPECT().IsAdmin(newAdminData.ProjectID, adminID).Return(nil)

				mock.ExpectBegin()
				mock.ExpectExec(
					regexp.QuoteMeta("UPDATE projects SET admin = $1 WHERE id = $2"),
				).WithArgs(
					newAdminData.NewAdminID,
					newAdminData.ProjectID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				log.EXPECT().Infof("Set new admin = %d in project = %d", newAdminData.NewAdminID, newAdminData.ProjectID)

				mock.ExpectExec(
					regexp.QuoteMeta("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2"),
				).WithArgs(
					newAdminData.ProjectID,
					newAdminData.NewAdminID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec(
					regexp.QuoteMeta("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)"),
				).WithArgs(
					newAdminData.ProjectID,
					adminID,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()

				return &ProjectRepository{db: db, log: log, admin: admin}
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.newAdminData, test.userID)
			err := repo.SetNewAdmin(test.newAdminData, test.userID)

			require.Equal(t, test.expectedError, err)
		})
	}
}

func Test_getProjectsByUserId(t *testing.T) {
	type mockBehaviour func(c *gomock.Controller, id uint64) *ProjectRepository
	err := errors.New("error")

	tests := []struct {
		name           string
		id             uint64
		mockBehaviour  mockBehaviour
		expectedResult []*models.Project
		expectedError  error
	}{
		{
			name: "Error",
			id:   1,
			mockBehaviour: func(c *gomock.Controller, id uint64) *ProjectRepository {
				log := mock_log.NewMockLog(c)
				db, mock, _ := sqlmock.New()

				mock.ExpectQuery(
					regexp.QuoteMeta(
						`SELECT * FROM projects WHERE projects.id IN (
							SELECT project_id FROM projects_members WHERE member_id = $1
						) UNION SELECT * FROM projects WHERE admin = $1`,
					),
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
					regexp.QuoteMeta(
						`SELECT * FROM projects WHERE projects.id IN (
							SELECT project_id FROM projects_members WHERE member_id = $1
						) UNION SELECT * FROM projects WHERE admin = $1`,
					),
				).WithArgs(id).WillReturnRows(rows)

				return &ProjectRepository{db: db, log: log}
			},
			expectedResult: []*models.Project{{
				ID:          1,
				Name:        "name",
				Description: "",
				AdminID:     1,
			}},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := test.mockBehaviour(c, test.id)
			res, err := repo.GetProjectsByUserId(test.id)

			require.Equal(t, test.expectedResult, res)
			require.Equal(t, test.expectedError, err)
		})
	}
}
