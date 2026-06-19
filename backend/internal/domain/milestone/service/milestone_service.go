package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/milestone/dto"
	"github.com/Kal-el21/backend/internal/domain/milestone/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/milestone/errors"
	"github.com/Kal-el21/backend/internal/domain/milestone/repository"
	"github.com/Kal-el21/backend/internal/domain/milestone/statemachine"
	"gorm.io/gorm"
)

// TaskProgressProvider menghindari circular dependency milestone <-> task.
// Diimplementasikan oleh task.service di bagian 4.
type TaskProgressProvider interface {
	GetAverageProgressByMilestone(milestoneID uint64) (float64, error)
}

type MilestoneService interface {
	Create(projectID uint64, req dto.CreateMilestoneRequest) (*entity.Milestone, error)
	GetByID(id uint64) (*entity.Milestone, error)
	GetByProjectID(projectID uint64) ([]dto.MilestoneResponse, error)
	Update(id uint64, req dto.UpdateMilestoneRequest) (*entity.Milestone, error)
	ChangeStatus(id uint64, req dto.ChangeMilestoneStatusRequest) (*entity.Milestone, error)
	Reorder(projectID uint64, req dto.ReorderMilestoneRequest) error
	Delete(id uint64, deletedBy uint64) error
	GetAverageProgressByProject(projectID uint64) (float64, error)
}

type milestoneService struct {
	repo         repository.MilestoneRepository
	taskProvider TaskProgressProvider
}

func NewMilestoneService(repo repository.MilestoneRepository, taskProvider TaskProgressProvider) MilestoneService {
	return &milestoneService{repo: repo, taskProvider: taskProvider}
}

func (s *milestoneService) Create(projectID uint64, req dto.CreateMilestoneRequest) (*entity.Milestone, error) {
	existing, _ := s.repo.FindByProjectID(projectID)

	milestone := &entity.Milestone{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      entity.MilestonePlanned,
		OrderIndex:  len(existing), // append di urutan terakhir
	}

	if err := s.repo.Create(milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

func (s *milestoneService) GetByID(id uint64) (*entity.Milestone, error) {
	milestone, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrMilestoneNotFound
		}
		return nil, err
	}
	return milestone, nil
}

func (s *milestoneService) GetByProjectID(projectID uint64) ([]dto.MilestoneResponse, error) {
	milestones, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.MilestoneResponse, len(milestones))
	for i, m := range milestones {
		progress, _ := s.taskProvider.GetAverageProgressByMilestone(m.ID)
		result[i] = dto.MilestoneResponse{
			ID:          m.ID,
			ProjectID:   m.ProjectID,
			Name:        m.Name,
			Description: m.Description,
			OrderIndex:  m.OrderIndex,
			StartDate:   m.StartDate,
			EndDate:     m.EndDate,
			Status:      string(m.Status),
			Progress:    progress,
			Version:     m.Version,
		}
	}

	return result, nil
}

func (s *milestoneService) Update(id uint64, req dto.UpdateMilestoneRequest) (*entity.Milestone, error) {
	milestone, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	milestone.Name = req.Name
	milestone.Description = req.Description
	milestone.StartDate = req.StartDate
	milestone.EndDate = req.EndDate

	rows, err := s.repo.UpdateWithVersionCheck(milestone, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

func (s *milestoneService) ChangeStatus(id uint64, req dto.ChangeMilestoneStatusRequest) (*entity.Milestone, error) {
	milestone, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	newStatus := entity.MilestoneStatus(req.Status)

	if err := statemachine.ValidateTransition(milestone.Status, newStatus); err != nil {
		return nil, err
	}

	milestone.Status = newStatus

	rows, err := s.repo.UpdateWithVersionCheck(milestone, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

func (s *milestoneService) Reorder(projectID uint64, req dto.ReorderMilestoneRequest) error {
	for index, milestoneID := range req.OrderedIDs {
		if err := s.repo.UpdateOrderIndex(milestoneID, index); err != nil {
			return err
		}
	}
	return nil
}

func (s *milestoneService) Delete(id uint64, deletedBy uint64) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id, deletedBy)
}

// GetAverageProgressByProject mengimplementasikan MilestoneProgressProvider
// yang dibutuhkan project.service (Project Progress = Average(Milestone Progress)).
func (s *milestoneService) GetAverageProgressByProject(projectID uint64) (float64, error) {
	milestones, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		return 0, err
	}

	if len(milestones) == 0 {
		return 0, nil
	}

	var total float64
	var counted int

	for _, m := range milestones {
		if m.Status == entity.MilestoneCancelled {
			continue // SDD section 12: cancelled tidak dihitung
		}
		progress, err := s.taskProvider.GetAverageProgressByMilestone(m.ID)
		if err != nil {
			return 0, err
		}
		total += progress
		counted++
	}

	if counted == 0 {
		return 0, nil
	}

	return total / float64(counted), nil
}
