package repository

import (
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	"gorm.io/gorm"
)

type BudgetRepository interface {
	Create(budget *entity.Budget) error
	FindByID(id uint64) (*entity.Budget, error)
	FindByProjectID(projectID uint64) (*entity.Budget, error)
	UpdateWithVersionCheck(budget *entity.Budget, expectedVersion int) (int64, error)
}

type budgetRepository struct {
	db *gorm.DB
}

func NewBudgetRepository(db *gorm.DB) BudgetRepository {
	return &budgetRepository{db: db}
}

func (r *budgetRepository) Create(budget *entity.Budget) error {
	return r.db.Create(budget).Error
}

func (r *budgetRepository) FindByID(id uint64) (*entity.Budget, error) {
	var budget entity.Budget
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) FindByProjectID(projectID uint64) (*entity.Budget, error) {
	var budget entity.Budget
	err := r.db.Where("project_id = ? AND deleted_at IS NULL", projectID).First(&budget).Error
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

func (r *budgetRepository) UpdateWithVersionCheck(budget *entity.Budget, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.Budget{}).
		Where("id = ? AND version = ?", budget.ID, expectedVersion).
		Updates(map[string]interface{}{
			"allocated_budget": budget.AllocatedBudget,
			"budget_type":      budget.BudgetType,
			"budget_name":      budget.BudgetName,
			"version":          gorm.Expr("version + 1"),
		})

	return result.RowsAffected, result.Error
}
