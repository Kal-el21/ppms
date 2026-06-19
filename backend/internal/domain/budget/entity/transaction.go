package entity

import "time"

type TransactionType string
type AdjustmentType string

const (
	TypeExpense    TransactionType = "EXPENSE"
	TypeRefund     TransactionType = "REFUND"
	TypeAdjustment TransactionType = "ADJUSTMENT"

	AdjustmentErrorCorrection    AdjustmentType = "ERROR_CORRECTION"
	AdjustmentBudgetReallocation AdjustmentType = "BUDGET_REALLOCATION"
	AdjustmentAuditCorrection    AdjustmentType = "AUDIT_CORRECTION"
	AdjustmentManualOverride     AdjustmentType = "MANUAL_OVERRIDE"
)

// BudgetTransaction bersifat IMMUTABLE: tidak ada method Update/Delete
// di repository layer untuk entity ini (lihat repository.go).
type BudgetTransaction struct {
	ID       uint64 `gorm:"primaryKey" json:"id"`
	BudgetID uint64 `gorm:"not null" json:"budget_id"`

	Type           TransactionType `gorm:"not null" json:"type"`
	AdjustmentType *AdjustmentType `json:"adjustment_type,omitempty"`

	Amount float64 `gorm:"not null" json:"amount"`

	Reason      string `json:"reason"`
	Description string `json:"description"`

	TransactionDate time.Time `json:"transaction_date"`

	IdempotencyKey *string `gorm:"unique" json:"idempotency_key,omitempty"`

	CreatedBy uint64 `gorm:"not null" json:"created_by"`

	CreatedAt time.Time `json:"created_at"`
}

func (BudgetTransaction) TableName() string {
	return "budget_transactions"
}
