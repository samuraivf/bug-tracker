package repository

import (
	"database/sql"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

type ProjectRepository struct {
	db  *sql.DB
	log log.Log
}

func NewProjectRepo(db *sql.DB, log log.Log) Project {
	return &ProjectRepository{db, log}
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
