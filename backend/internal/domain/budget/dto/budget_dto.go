package dto

import "time"

type CreateBudgetRequest struct {
	AllocatedBudget float64 `json:"allocated_budget" validate:"required,gt=0"`
	BudgetType      string  `json:"budget_type" validate:"omitempty,oneof=CAPEX OPEX"`
	BudgetName      string  `json:"budget_name" validate:"max=200"`
}

type UpdateBudgetRequest struct {
	AllocatedBudget float64 `json:"allocated_budget" validate:"required,gt=0"`
	BudgetType      string  `json:"budget_type" validate:"omitempty,oneof=CAPEX OPEX"`
	BudgetName      string  `json:"budget_name" validate:"max=200"`
	Version         int     `json:"version" validate:"required"`
}

type CreateTransactionRequest struct {
	Type           string  `json:"type" validate:"required,oneof=EXPENSE REFUND ADJUSTMENT"`
	AdjustmentType string  `json:"adjustment_type" validate:"omitempty,oneof=ERROR_CORRECTION BUDGET_REALLOCATION AUDIT_CORRECTION MANUAL_OVERRIDE"`
	Amount         float64 `json:"amount" validate:"required,gt=0"`
	Reason         string  `json:"reason"`
	Description    string  `json:"description" validate:"max=1000"`

	// IdempotencyKey dikirim client (umumnya UUID v4 yang digenerate saat user klik submit)
	// untuk mencegah duplicate transaction akibat retry/double-click.
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`
}

type BudgetResponse struct {
	ID              uint64  `json:"id"`
	ProjectID       uint64  `json:"project_id"`
	BudgetType      *string `json:"budget_type,omitempty"`
	BudgetName      string  `json:"budget_name"`
	AllocatedBudget float64 `json:"allocated_budget"`
	UsedBudget      float64 `json:"used_budget"`
	RemainingBudget float64 `json:"remaining_budget"`
	UsagePercentage float64 `json:"usage_percentage"`
	Version         int     `json:"version"`
}

type TransactionResponse struct {
	ID              uint64    `json:"id"`
	BudgetID        uint64    `json:"budget_id"`
	Type            string    `json:"type"`
	AdjustmentType  *string   `json:"adjustment_type,omitempty"`
	Amount          float64   `json:"amount"`
	Reason          string    `json:"reason"`
	Description     string    `json:"description"`
	TransactionDate time.Time `json:"transaction_date"`
	CreatedBy       uint64    `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}
