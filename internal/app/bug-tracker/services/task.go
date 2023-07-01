package services

import (
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type TaskService struct {
	repo repository.Task
}

func NewTask(repo repository.Task) Task {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(taskData *dto.CreateTaskDto, userID uint64) (uint64, error) {
	return s.repo.CreateTask(taskData, userID)
}

func (s *TaskService) WorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	return s.repo.WorkOnTask(workOnTaskData, userID)
}

func (s *TaskService) StopWorkOnTask(workOnTaskData *dto.WorkOnTaskDto, userID uint64) error {
	return s.repo.StopWorkOnTask(workOnTaskData, userID)
}

func (s *TaskService) UpdateTask(taskData *dto.UpdateTaskDto, userID uint64) (uint64, error) {
	return s.repo.UpdateTask(taskData, userID)
}

func (s *TaskService) GetTaskById(id uint64) (*models.Task, error) {
	return s.repo.GetTaskById(id)
}
