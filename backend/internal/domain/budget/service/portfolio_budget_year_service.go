package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	"github.com/Kal-el21/backend/internal/domain/budget/repository"
	domainerrors "github.com/Kal-el21/backend/internal/domain/budget/errors"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type PortfolioBudgetYearService interface {
	Create(createdBy uint64, req dto.CreatePortfolioBudgetYearRequest) (*entity.PortfolioBudgetYear, error)
	GetAll() ([]entity.PortfolioBudgetYear, error)
	GetByID(id uint64) (*entity.PortfolioBudgetYear, error)
	Update(id uint64, updatedBy uint64, req dto.UpdatePortfolioBudgetYearRequest) (*entity.PortfolioBudgetYear, error)
	Delete(id uint64) error
}

type portfolioBudgetYearService struct {
	repo repository.PortfolioBudgetYearRepository
}

func NewPortfolioBudgetYearService(repo repository.PortfolioBudgetYearRepository) PortfolioBudgetYearService {
	return &portfolioBudgetYearService{repo: repo}
}

func (s *portfolioBudgetYearService) Create(createdBy uint64, req dto.CreatePortfolioBudgetYearRequest) (*entity.PortfolioBudgetYear, error) {
	if req.Year < 2000 || req.Year > 2100 {
		return nil, apperrors.New(apperrors.ErrValidation, "year must be between 2000 and 2100")
	}
	if req.CapexCeiling < 0 || req.OpexCeiling < 0 {
		return nil, apperrors.New(apperrors.ErrValidation, "ceiling must be >= 0")
	}

	existing, _ := s.repo.FindByYear(req.Year)
	if existing != nil {
		return nil, domainerrors.ErrBudgetAlreadyExists
	}

	creator := createdBy
	year := &entity.PortfolioBudgetYear{
		Year:         req.Year,
		CapexCeiling: req.CapexCeiling,
		OpexCeiling:  req.OpexCeiling,
		CreatedBy:    &creator,
	}
	if err := s.repo.Create(year); err != nil {
		return nil, err
	}
	return year, nil
}

func (s *portfolioBudgetYearService) GetAll() ([]entity.PortfolioBudgetYear, error) {
	return s.repo.FindAll()
}

func (s *portfolioBudgetYearService) GetByID(id uint64) (*entity.PortfolioBudgetYear, error) {
	year, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrBudgetNotFound
		}
		return nil, err
	}
	return year, nil
}

func (s *portfolioBudgetYearService) Update(id uint64, updatedBy uint64, req dto.UpdatePortfolioBudgetYearRequest) (*entity.PortfolioBudgetYear, error) {
	year, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrBudgetNotFound
		}
		return nil, err
	}
	if req.CapexCeiling < 0 || req.OpexCeiling < 0 {
		return nil, apperrors.New(apperrors.ErrValidation, "ceiling must be >= 0")
	}

	year.CapexCeiling = req.CapexCeiling
	year.OpexCeiling = req.OpexCeiling
	updater := updatedBy
	year.UpdatedBy = &updater

	rows, err := s.repo.UpdateWithVersionCheck(year, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}
	return s.repo.FindByID(id)
}

func (s *portfolioBudgetYearService) Delete(id uint64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domainerrors.ErrBudgetNotFound
		}
		return err
	}
	return s.repo.Delete(id)
}
