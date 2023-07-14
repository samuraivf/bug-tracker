package repository

import (
	"database/sql"
	"errors"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

var (
	ErrNoRights = errors.New("error no rights to do this operation")
)

type ProjectRepository struct {
	db     *sql.DB
	log    log.Log
	admin  admin
	member member
}

func NewProjectRepo(db *sql.DB, log log.Log, admin admin, member member) Project {
	return &ProjectRepository{
		db:     db,
		log:    log,
		admin:  admin,
		member: member,
	}
}

func (r *ProjectRepository) CreateProject(projectDto *dto.CreateProjectDto) (uint64, error) {
	result := r.db.QueryRow(
		"INSERT INTO projects (name, description, admin) VALUES ($1, $2, $3) RETURNING id",
		projectDto.Name,
		projectDto.Description,
		projectDto.AdminID,
	)

	var projectID uint64
	if err := result.Scan(&projectID); err != nil {
		r.log.Error(err)
		return 0, err
	}
	r.log.Infof("Create project: id = %d", projectID)

	return projectID, nil
}

func (r *ProjectRepository) GetProjectById(id uint64) (*models.Project, error) {
	result := r.db.QueryRow("SELECT * FROM projects WHERE id = $1", id)

	project := new(models.Project)
	if err := result.Scan(&project.ID, &project.Name, &project.Description, &project.AdminID); err != nil {
		r.log.Error(err)
		return nil, err
	}
	r.log.Infof("Get project: id = %d", id)

	return project, nil
}

func (r *ProjectRepository) DeleteProject(projectID, userID uint64) error {
	if err := r.admin.IsAdmin(projectID, userID); err != nil {
		return err
	}

	_, err := r.db.Exec("DELETE FROM projects WHERE id = $1", projectID)

	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *ProjectRepository) UpdateProject(projectData *dto.UpdateProjectDto, userID uint64) error {
	if err := r.admin.IsAdmin(projectData.ProjectID, userID); err != nil {
		return err
	}

	_, err := r.db.Exec(
		"UPDATE projects SET description = $1 WHERE id = $2",
		projectData.Description,
		projectData.ProjectID,
	)

	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *ProjectRepository) AddMember(memberData *dto.AddMemberDto, userID uint64) error {
	if err := r.admin.IsAdmin(memberData.ProjectID, userID); err != nil {
		return err
	}

	_, err := r.db.Exec(
		"INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)",
		memberData.ProjectID,
		memberData.MemberID,
	)

	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *ProjectRepository) DeleteMember(memberData *dto.AddMemberDto, userID uint64) error {
	if err := r.admin.IsAdmin(memberData.ProjectID, userID); err != nil {
		return err
	}

	_, err := r.db.Exec(
		"DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2",
		memberData.ProjectID,
		memberData.MemberID,
	)

	if err != nil {
		r.log.Error(err)
		return err
	}
	r.log.Infof("Delete member with id=%d from project with id=%d", memberData.MemberID, memberData.ProjectID)

	return nil
}

func (r *ProjectRepository) GetMembers(projectID, userID uint64) ([]*models.User, error) {
	if r.member.IsMember(projectID, userID) != nil && r.admin.IsAdmin(projectID, userID) != nil {
		return nil, ErrNoRights
	}

	rows, err := r.db.Query("SELECT * FROM users WHERE users.id IN (SELECT member_id FROM projects_members WHERE projects_members.project_id = $1)", projectID)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	members := make([]*models.User, 0)
	for rows.Next() {
		member := new(models.User)
		err := rows.Scan(&member.ID, &member.Name, &member.Username, &member.Password, &member.Email)
		if err != nil {
			r.log.Error(err)
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (r *ProjectRepository) LeaveProject(projectID, userID uint64) error {
	if err := r.admin.IsAdmin(projectID, userID); err == nil {
		return ErrNoRights
	}

	_, err := r.db.Exec("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2", projectID, userID)
	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *ProjectRepository) SetNewAdmin(newAdminData *dto.NewAdminDto, adminID uint64) error {
	if err := r.admin.IsAdmin(newAdminData.ProjectID, adminID); err != nil {
		return err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE projects SET admin = $1 WHERE id = $2", newAdminData.NewAdminID, newAdminData.ProjectID)
	if err != nil {
		r.log.Error(err)
		tx.Rollback()
		return err
	}
	r.log.Infof("Set new admin = %d in project = %d", newAdminData.NewAdminID, newAdminData.ProjectID)

	_, err = tx.Exec("DELETE FROM projects_members WHERE project_id = $1 AND member_id = $2", newAdminData.ProjectID, newAdminData.NewAdminID)
	if err != nil {
		r.log.Error(err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO projects_members (project_id, member_id) VALUES ($1, $2)", newAdminData.ProjectID, adminID)
	if err != nil {
		r.log.Error(err)
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		r.log.Error(err)
	}

	return err
}

func (r *ProjectRepository) GetProjectsByUserId(id uint64) ([]*models.Project, error) {
	rows, err := r.db.Query(
		`SELECT * FROM projects WHERE projects.id IN (
			SELECT project_id FROM projects_members WHERE member_id = $1
		) UNION SELECT * FROM projects WHERE admin = $1`,
		id,
	)
	if err != nil {
		r.log.Error(err)
		return nil, err
	}

	projects := make([]*models.Project, 0)
	for rows.Next() {
		project := new(models.Project)
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.AdminID); err != nil {
			r.log.Error(err)
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}
