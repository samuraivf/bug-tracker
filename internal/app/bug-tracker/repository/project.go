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
	db    *sql.DB
	log   log.Log
	admin admin
}

func NewProjectRepo(db *sql.DB, log log.Log) Project {
	return &ProjectRepository{
		db:    db,
		log:   log,
		admin: new_adminStrategy(db, log),
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

func (r *ProjectRepository) LeaveProject(projectID, userID uint64) error {
	if err := r.admin.IsAdmin(projectID, userID); err == nil {
		return ErrNoRights
	}

	_, err := r.db.Exec("DELETE FROM projects_members WHERE member_id = $1", userID)
	if err != nil {
		r.log.Error(err)
	}

	return err
}
