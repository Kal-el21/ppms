package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/division/dto"
	"github.com/Kal-el21/backend/internal/domain/division/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/division/errors"
	"github.com/Kal-el21/backend/internal/domain/division/repository"
	"gorm.io/gorm"
)

type DivisionService interface {
	Create(req dto.CreateDivisionRequest) (*entity.Division, error)
	GetByID(id uint64) (*entity.Division, error)
	GetAll() ([]entity.Division, error)
	Update(id uint64, req dto.UpdateDivisionRequest) (*entity.Division, error)
	Delete(id uint64, deletedBy uint64) error
}

type divisionService struct {
	repo repository.DivisionRepository
}

func NewDivisionService(repo repository.DivisionRepository) DivisionService {
	return &divisionService{repo: repo}
}

func (s *divisionService) Create(req dto.CreateDivisionRequest) (*entity.Division, error) {
	existing, err := s.repo.FindByName(req.Name)
	if err == nil && existing != nil {
		return nil, domainerrors.ErrDivisionNameTaken
	}

	division := &entity.Division{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(division); err != nil {
		return nil, err
	}

	return division, nil
}

func (s *divisionService) GetByID(id uint64) (*entity.Division, error) {
	division, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrDivisionNotFound
		}
		return nil, err
	}
	return division, nil
}

func (s *divisionService) GetAll() ([]entity.Division, error) {
	return s.repo.FindAll()
}

func (s *divisionService) Update(id uint64, req dto.UpdateDivisionRequest) (*entity.Division, error) {
	division, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if division.Name != req.Name {
		existing, err := s.repo.FindByName(req.Name)
		if err == nil && existing != nil && existing.ID != id {
			return nil, domainerrors.ErrDivisionNameTaken
		}
	}

	division.Name = req.Name
	division.Description = req.Description

	if err := s.repo.Update(division); err != nil {
		return nil, err
	}

	return division, nil
}

func (s *divisionService) Delete(id uint64, deletedBy uint64) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id, deletedBy)
}
