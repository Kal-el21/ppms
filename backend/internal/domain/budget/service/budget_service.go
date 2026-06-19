package service

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/budget/errors"
	"github.com/Kal-el21/backend/internal/domain/budget/repository"
	"gorm.io/gorm"
)

type BudgetService interface {
	Create(projectID uint64, req dto.CreateBudgetRequest) (*entity.Budget, error)
	GetByProjectID(projectID uint64) (*dto.BudgetResponse, error)
	Update(id uint64, req dto.UpdateBudgetRequest) (*entity.Budget, error)
}

type budgetService struct {
	repo   repository.BudgetRepository
	txRepo repository.TransactionRepository
}

func NewBudgetService(repo repository.BudgetRepository, txRepo repository.TransactionRepository) BudgetService {
	return &budgetService{repo: repo, txRepo: txRepo}
}

func (s *budgetService) Create(projectID uint64, req dto.CreateBudgetRequest) (*entity.Budget, error) {
	existing, _ := s.repo.FindByProjectID(projectID)
	if existing != nil {
		return nil, domainerrors.ErrBudgetAlreadyExists
	}

	budget := &entity.Budget{
		ProjectID:       projectID,
		AllocatedBudget: req.AllocatedBudget,
	}

	if err := s.repo.Create(budget); err != nil {
		return nil, err
	}

	return budget, nil
}

func (s *budgetService) GetByProjectID(projectID uint64) (*dto.BudgetResponse, error) {
	budget, err := s.repo.FindByProjectID(projectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrBudgetNotFound
		}
		return nil, err
	}

	used, err := s.txRepo.SumByBudgetID(budget.ID)
	if err != nil {
		return nil, err
	}

	remaining := budget.AllocatedBudget - used
	var usagePercentage float64
	if budget.AllocatedBudget > 0 {
		usagePercentage = (used / budget.AllocatedBudget) * 100
	}

	return &dto.BudgetResponse{
		ID:              budget.ID,
		ProjectID:       budget.ProjectID,
		AllocatedBudget: budget.AllocatedBudget,
		UsedBudget:      used,
		RemainingBudget: remaining,
		UsagePercentage: usagePercentage,
		Version:         budget.Version,
	}, nil
}

func (s *budgetService) Update(id uint64, req dto.UpdateBudgetRequest) (*entity.Budget, error) {
	budget, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrBudgetNotFound
		}
		return nil, err
	}

	budget.AllocatedBudget = req.AllocatedBudget

	rows, err := s.repo.UpdateWithVersionCheck(budget, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}
