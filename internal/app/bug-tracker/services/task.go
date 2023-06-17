package services

import (
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
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
