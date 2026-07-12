package repository

import (
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	"gorm.io/gorm"
)

type PortfolioBudgetYearRepository interface {
	Create(year *entity.PortfolioBudgetYear) error
	FindByID(id uint64) (*entity.PortfolioBudgetYear, error)
	FindByYear(year int) (*entity.PortfolioBudgetYear, error)
	FindAll() ([]entity.PortfolioBudgetYear, error)
	UpdateWithVersionCheck(year *entity.PortfolioBudgetYear, expectedVersion int) (int64, error)
	Delete(id uint64) error
}

type portfolioBudgetYearRepository struct {
	db *gorm.DB
}

func NewPortfolioBudgetYearRepository(db *gorm.DB) PortfolioBudgetYearRepository {
	return &portfolioBudgetYearRepository{db: db}
}

func (r *portfolioBudgetYearRepository) Create(year *entity.PortfolioBudgetYear) error {
	return r.db.Create(year).Error
}

func (r *portfolioBudgetYearRepository) FindByID(id uint64) (*entity.PortfolioBudgetYear, error) {
	var year entity.PortfolioBudgetYear
	err := r.db.Where("id = ?", id).First(&year).Error
	if err != nil {
		return nil, err
	}
	return &year, nil
}

func (r *portfolioBudgetYearRepository) FindByYear(year int) (*entity.PortfolioBudgetYear, error) {
	var result entity.PortfolioBudgetYear
	err := r.db.Where("year = ?", year).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *portfolioBudgetYearRepository) FindAll() ([]entity.PortfolioBudgetYear, error) {
	var years []entity.PortfolioBudgetYear
	err := r.db.Order("year DESC").Find(&years).Error
	return years, err
}

func (r *portfolioBudgetYearRepository) UpdateWithVersionCheck(year *entity.PortfolioBudgetYear, expectedVersion int) (int64, error) {
	result := r.db.Model(&entity.PortfolioBudgetYear{}).
		Where("id = ? AND version = ?", year.ID, expectedVersion).
		Updates(map[string]interface{}{
			"capex_ceiling": year.CapexCeiling,
			"opex_ceiling":  year.OpexCeiling,
			"version":       gorm.Expr("version + 1"),
		})
	return result.RowsAffected, result.Error
}

func (r *portfolioBudgetYearRepository) Delete(id uint64) error {
	return r.db.Delete(&entity.PortfolioBudgetYear{}, id).Error
}
