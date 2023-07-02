package services

import (
	"context"
	"time"

	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/redis"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

//go:generate mockgen -source=services.go -destination=mocks/services.go

type Auth interface {
	GetRefreshTokenTTL() time.Duration
	GenerateAccessToken(username string, userID uint64) (string, error)
	GenerateRefreshToken(username string, userID uint64) (*RefreshTokenData, error)
	ParseAccessToken(accessToken string) (*TokenData, error)
	ParseRefreshToken(refreshToken string) (*TokenData, error)
}

type User interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id uint64) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(userData *dto.SignUpDto) (uint64, error)
	ValidateUser(email, password string) (*models.User, error)
}

type Redis interface {
	Set(ctx context.Context, key, val string, exp time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	SetRefreshToken(ctx context.Context, key, refreshToken string) error
	GetRefreshToken(ctx context.Context, key string) (string, error)
	DeleteRefreshToken(ctx context.Context, key string) error
	Close() error
}

type Project interface {
	CreateProject(projectData *dto.CreateProjectDto) (uint64, error)
	GetProjectById(id uint64) (*models.Project, error)
	DeleteProject(projectID, userID uint64) error
	UpdateProject(projectData *dto.UpdateProjectDto, userID uint64) error
	AddMember(memberData *dto.AddMemberDto, userID uint64) error
	DeleteMember(memberData *dto.AddMemberDto, userID uint64) error
	LeaveProject(projectID, userID uint64) error
	SetNewAdmin(newAdmintData *dto.NewAdminDto, adminID uint64) error
}

type Task interface {
	CreateTask(taskData *dto.CreateTaskDto, userID uint64) (uint64, error)
	WorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error
	StopWorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error
	UpdateTask(taskData *dto.UpdateTaskDto, userID uint64) (uint64, error)
	GetTaskById(id uint64) (*models.Task, error)
	DeleteTask(taskData *dto.DeleteTaskDto, userID uint64) error
}

type Service struct {
	Auth
	User
	Redis
	Project
	Task
}

func NewService(repo *repository.Repository, redisRepo redis.Redis) *Service {
	return &Service{
		Auth:    NewAuth(),
		User:    NewUser(repo.User),
		Redis:   NewRedis(redisRepo),
		Project: NewProject(repo.Project),
		Task:    NewTask(repo.Task),
	}
}
