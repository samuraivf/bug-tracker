package services

import (
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/dto"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/models"
	"github.com/samuraivf/bug-tracker/internal/app/bug-tracker/repository"
)

type ProjectService struct {
	repo repository.Project
}

func NewProject(repo repository.Project) Project {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(projectData *dto.CreateProjectDto) (uint64, error) {
	return s.repo.CreateProject(projectData)
}

func (s *ProjectService) GetProjectById(id uint64) (*models.Project, error) {
	return s.repo.GetProjectById(id)
}

func (s *ProjectService) DeleteProject(projectID, userID uint64) error {
	return s.repo.DeleteProject(projectID, userID)
}

func (s *ProjectService) UpdateProject(projectData *dto.UpdateProjectDto, userID uint64) error {
	return s.repo.UpdateProject(projectData, userID)
}

func (s *ProjectService) AddMember(memberData *dto.AddMemberDto, userID uint64) error {
	return s.repo.AddMember(memberData, userID)
}

func (s *ProjectService) DeleteMember(memberData *dto.AddMemberDto, userID uint64) error {
	return s.repo.DeleteMember(memberData, userID)
}

func (s *ProjectService) LeaveProject(projectID, userID uint64) error {
	return s.repo.LeaveProject(projectID, userID)
}

func (s *ProjectService) SetNewAdmin(newAdminData *dto.NewAdminDto, adminID uint64) error {
	return s.repo.SetNewAdmin(newAdminData, adminID)
}
