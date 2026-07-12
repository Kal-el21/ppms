package service

import (
	"errors"
	"strings"

	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/budget/errors"
	"github.com/Kal-el21/backend/internal/domain/budget/repository"
	apperrors "github.com/Kal-el21/backend/internal/shared/errors"
	"gorm.io/gorm"
)

type BudgetService interface {
	Create(projectID uint64, createdBy uint64, req dto.CreateBudgetRequest) (*entity.Budget, error)
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

func (s *budgetService) Create(projectID uint64, createdBy uint64, req dto.CreateBudgetRequest) (*entity.Budget, error) {
	existing, _ := s.repo.FindByProjectID(projectID)
	if existing != nil {
		return nil, domainerrors.ErrBudgetAlreadyExists
	}
	if req.AllocatedBudget <= 0 {
		return nil, apperrors.New(apperrors.ErrValidation, "allocated_budget must be greater than 0")
	}

	budgetType, err := normalizeBudgetType(req.BudgetType)
	if err != nil {
		return nil, err
	}
	budgetName, err := normalizeBudgetName(req.BudgetName)
	if err != nil {
		return nil, err
	}

	budget := &entity.Budget{
		ProjectID:       projectID,
		BudgetType:      budgetType,
		BudgetName:      budgetName,
		AllocatedBudget: req.AllocatedBudget,
		CreatedBy:       createdBy,
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
		BudgetType:      budget.BudgetType,
		BudgetName:      budget.BudgetName,
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
	if req.AllocatedBudget <= 0 {
		return nil, apperrors.New(apperrors.ErrValidation, "allocated_budget must be greater than 0")
	}

	budgetType, err := normalizeBudgetType(req.BudgetType)
	if err != nil {
		return nil, err
	}
	budgetName, err := normalizeBudgetName(req.BudgetName)
	if err != nil {
		return nil, err
	}

	budget.AllocatedBudget = req.AllocatedBudget
	budget.BudgetType = budgetType
	budget.BudgetName = budgetName

	rows, err := s.repo.UpdateWithVersionCheck(budget, req.Version)
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, domainerrors.ErrVersionMismatch
	}

	return s.repo.FindByID(id)
}

func normalizeBudgetType(value string) (*string, error) {
	normalized := strings.TrimSpace(value)
	normalized = strings.ToUpper(normalized)
	normalized = strings.NewReplacer(" ", "_", "-", "_").Replace(normalized)
	if normalized == "" {
		return nil, nil
	}
	if normalized != "CAPEX" && normalized != "OPEX" {
		return nil, apperrors.New(apperrors.ErrValidation, "budget_type is invalid")
	}
	return &normalized, nil
}

func normalizeBudgetName(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 200 {
		return "", apperrors.New(apperrors.ErrValidation, "budget_name is too long")
	}
	return value, nil
}
