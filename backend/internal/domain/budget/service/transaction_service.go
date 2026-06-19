package service

import (
	"github.com/Kal-el21/backend/internal/domain/budget/dto"
	"github.com/Kal-el21/backend/internal/domain/budget/entity"
	domainerrors "github.com/Kal-el21/backend/internal/domain/budget/errors"
	"github.com/Kal-el21/backend/internal/domain/budget/repository"
	"github.com/Kal-el21/backend/internal/events"
)

const budgetWarningThreshold = 80.0 // persen, sesuai SDD event budget.warning

type TransactionService interface {
	Create(budgetID uint64, createdBy uint64, req dto.CreateTransactionRequest) (*entity.BudgetTransaction, error)
	GetByBudgetID(budgetID uint64) ([]entity.BudgetTransaction, error)
	GetByBudgetIDPaginated(budgetID uint64, page, limit int) ([]entity.BudgetTransaction, int64, error)
}

type transactionService struct {
	txRepo     repository.TransactionRepository
	budgetRepo repository.BudgetRepository
	eventBus   *events.Bus
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	budgetRepo repository.BudgetRepository,
	eventBus *events.Bus,
) TransactionService {
	return &transactionService{txRepo: txRepo, budgetRepo: budgetRepo, eventBus: eventBus}
}

// Create menerapkan beberapa business rules penting:
//  1. Idempotency check (mencegah duplicate akibat retry/double-click)
//  2. Validasi ADJUSTMENT wajib adjustment_type + reason (FR-08.04, DB CHECK constraint sbg safety net kedua)
//  3. Konvensi amount: EXPENSE disimpan sebagai POSITIVE (menambah used_budget),
//     REFUND disimpan sebagai NEGATIVE (mengurangi used_budget),
//     ADJUSTMENT bisa POSITIVE/NEGATIVE tergantung reason
//     (caller/frontend mengirim amount sesuai konvensi ini; service tidak membalik sign secara implisit
//     untuk EXPENSE/REFUND agar histori sesuai input asli dan auditable).
//  4. Publish event budget.warning / budget.over_limit berdasarkan threshold
func (s *transactionService) Create(budgetID uint64, createdBy uint64, req dto.CreateTransactionRequest) (*entity.BudgetTransaction, error) {
	// Idempotency check — jika key sudah pernah dipakai, kembalikan transaksi yang sudah ada
	// alih-alih membuat duplikat (FR-17 spirit: defend against duplicate submission)
	existing, err := s.txRepo.FindByIdempotencyKey(req.IdempotencyKey)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	if req.Type == "ADJUSTMENT" && (req.AdjustmentType == "" || req.Reason == "") {
		return nil, domainerrors.ErrAdjustmentRequiresType
	}

	// Normalisasi sign sesuai konvensi: REFUND selalu disimpan negative
	// agar SumByBudgetID otomatis mengurangi used_budget tanpa logic tambahan.
	amount := req.Amount
	if req.Type == "REFUND" {
		amount = -amount
	}

	tx := &entity.BudgetTransaction{
		BudgetID:       budgetID,
		Type:           entity.TransactionType(req.Type),
		Amount:         amount,
		Reason:         req.Reason,
		Description:    req.Description,
		IdempotencyKey: &req.IdempotencyKey,
		CreatedBy:      createdBy,
	}

	if req.AdjustmentType != "" {
		adjType := entity.AdjustmentType(req.AdjustmentType)
		tx.AdjustmentType = &adjType
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}

	s.checkBudgetThreshold(budgetID)

	return tx, nil
}

func (s *transactionService) checkBudgetThreshold(budgetID uint64) {
	budget, err := s.budgetRepo.FindByID(budgetID)
	if err != nil || budget.AllocatedBudget <= 0 {
		return
	}

	used, err := s.txRepo.SumByBudgetID(budgetID)
	if err != nil {
		return
	}

	usagePercentage := (used / budget.AllocatedBudget) * 100

	if usagePercentage >= 100 {
		s.eventBus.Publish(events.Event{
			Name: "budget.over_limit",
			Data: map[string]interface{}{
				"budget_id":  budgetID,
				"project_id": budget.ProjectID,
				"usage_pct":  usagePercentage,
			},
		})
	} else if usagePercentage >= budgetWarningThreshold {
		s.eventBus.Publish(events.Event{
			Name: "budget.warning",
			Data: map[string]interface{}{
				"budget_id":  budgetID,
				"project_id": budget.ProjectID,
				"usage_pct":  usagePercentage,
			},
		})
	}
}

func (s *transactionService) GetByBudgetID(budgetID uint64) ([]entity.BudgetTransaction, error) {
	return s.txRepo.FindByBudgetID(budgetID)
}

func (s *transactionService) GetByBudgetIDPaginated(budgetID uint64, page, limit int) ([]entity.BudgetTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.txRepo.FindByBudgetIDPaginated(budgetID, page, limit)
}
