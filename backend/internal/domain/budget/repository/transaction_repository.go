package repository

import (
	"errors"

	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	"gorm.io/gorm"
)

// TransactionRepository SENGAJA tidak menyediakan method Update atau Delete.
// Ini menegakkan business rule "Budget transaction IMMUTABLE" di level kode,
// bukan hanya di level dokumentasi — tidak ada cara memanggil UPDATE/DELETE
// pada budget_transactions melalui repository ini.
type TransactionRepository interface {
	Create(tx *entity.BudgetTransaction) error
	FindByID(id uint64) (*entity.BudgetTransaction, error)
	FindByBudgetID(budgetID uint64) ([]entity.BudgetTransaction, error)
	FindByIdempotencyKey(key string) (*entity.BudgetTransaction, error)
	FindByBudgetIDPaginated(budgetID uint64, page, limit int) ([]entity.BudgetTransaction, int64, error)
	SumByBudgetID(budgetID uint64) (float64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *entity.BudgetTransaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) FindByID(id uint64) (*entity.BudgetTransaction, error) {
	var tx entity.BudgetTransaction
	err := r.db.Where("id = ?", id).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindByBudgetID(budgetID uint64) ([]entity.BudgetTransaction, error) {
	var txs []entity.BudgetTransaction
	err := r.db.Where("budget_id = ?", budgetID).Order("transaction_date DESC").Find(&txs).Error
	return txs, err
}

func (r *transactionRepository) FindByIdempotencyKey(key string) (*entity.BudgetTransaction, error) {
	var tx entity.BudgetTransaction
	err := r.db.Where("idempotency_key = ?", key).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindByBudgetIDPaginated(budgetID uint64, page, limit int) ([]entity.BudgetTransaction, int64, error) {
	var txs []entity.BudgetTransaction
	var total int64
	offset := (page - 1) * limit

	if err := r.db.Model(&entity.BudgetTransaction{}).Where("budget_id = ?", budgetID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where("budget_id = ?", budgetID).
		Order("transaction_date DESC").
		Limit(limit).Offset(offset).
		Find(&txs).Error

	return txs, total, err
}

// SumByBudgetID menghitung used_budget: EXPENSE menambah, REFUND & ADJUSTMENT
// dengan amount negatif mengurangi, sesuai konvensi amount yang ditentukan
// di service layer (lihat penjelasan di transaction_service.go).
func (r *transactionRepository) SumByBudgetID(budgetID uint64) (float64, error) {
	var total float64
	err := r.db.Model(&entity.BudgetTransaction{}).
		Where("budget_id = ?", budgetID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}
