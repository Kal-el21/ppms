package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/project/dto"
	"github.com/Kal-el21/backend/internal/domain/project/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/project/errors"
	"github.com/Kal-el21/backend/internal/domain/project/repository"
	"github.com/Kal-el21/backend/internal/domain/project/statemachine"
	"gorm.io/gorm"
)

type ProjectService interface {
	GetByID(id uint64) (*entity.Project, error)
	GetAll(page, limit int, statusFilter string) ([]entity.Project, int64, error)
	GetMyProjects(userID uint64, page, limit int) ([]entity.Project, int64, error)
	Update(id uint64, req dto.UpdateProjectRequest) (*entity.Project, error)
	ChangeStatus(id uint64, req dto.ChangeStatusRequest) (*entity.Project, error)
	CalculateProgress(id uint64) (float64, error)
}

type projectService struct {
	repo          repository.ProjectRepository
	milestoneRepo MilestoneProgressProvider
}

// MilestoneProgressProvider adalah interface kecil untuk menghindari
// circular dependency antara domain project <-> milestone.
// Diimplementasikan oleh milestone.service di bagian 3.
type MilestoneProgressProvider interface {
	GetAverageProgressByProject(projectID uint64) (float64, error)
}

func NewProjectService(repo repository.ProjectRepository, milestoneRepo MilestoneProgressProvider) ProjectService {
	return &projectService{repo: repo, milestoneRepo: milestoneRepo}
}

func (s *projectService) GetByID(id uint64) (*entity.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrProjectNotFound
		}
		return nil, err
	}
	return project, nil
}

func (s *projectService) GetAll(page, limit int, statusFilter string) ([]entity.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindAll(page, limit, statusFilter)
}

func (s *projectService) GetMyProjects(userID uint64, page, limit int) ([]entity.Project, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindByMemberUserID(userID, page, limit)
}

func (s *projectService) Update(id uint64, req dto.UpdateProjectRequest) (*entity.Project, error) {
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	project.Name = req.Name
	project.Description = req.Description
	project.StartDate = req.StartDate
	project.EndDate = req.EndDate

	rows, err := s.repo.UpdateWithVersionCheck(project, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

func (s *projectService) ChangeStatus(id uint64, req dto.ChangeStatusRequest) (*entity.Project, error) {
	project, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	newStatus := entity.ProjectStatus(req.Status)

	if err := statemachine.ValidateTransition(project.Status, newStatus); err != nil {
		return nil, err
	}

	project.Status = newStatus

	rows, err := s.repo.UpdateWithVersionCheck(project, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

// CalculateProgress: Project Progress = Average(Milestone Progress) per SDD section 12.
func (s *projectService) CalculateProgress(id uint64) (float64, error) {
	return s.milestoneRepo.GetAverageProgressByProject(id)
}
