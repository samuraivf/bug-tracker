package repository

import (
	"database/sql"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/log"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository.go

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(userData *dto.SignUpDto) (uint64, error)
}

type Project interface {
	CreateProject(projectData *dto.CreateProjectDto) (uint64, error)
	GetProjectById(id uint64) (*models.Project, error)
	DeleteProject(projectID, userID uint64) error
	UpdateProject(projectData *dto.UpdateProjectDto, userID uint64) error
	AddMember(memberData *dto.AddMemberDto, userID uint64) error
	DeleteMember(memberData *dto.AddMemberDto, userID uint64) error
	LeaveProject(projectID, userID uint64) error
	SetNewAdmin(newAdminData *dto.NewAdminDto, adminID uint64) error
}

type Task interface {
	CreateTask(taskData *dto.CreateTaskDto, userID uint64) (uint64, error)
	WorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error
	UpdateTask(taskData *dto.UpdateTaskDto, userID uint64) (uint64, error)
	GetTaskById(it uint64) (*models.Task, error)
}

type Repository struct {
	User
	Project
	Task
}

func NewRepository(db *sql.DB, log log.Log) *Repository {
	admin := new_adminStrategy(db, log)

	return &Repository{
		User:    NewUserRepo(db, log),
		Project: NewProjectRepo(db, log, admin),
		Task:    NewTaskRepo(db, log, admin),
	}
}
