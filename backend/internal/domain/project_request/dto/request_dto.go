package dto

import "time"

type ProjectRequestMetadata struct {
	Category          string     `json:"category" validate:"max=100"`
	InitiationType    string     `json:"initiation_type" validate:"omitempty,oneof=NEW_INITIATIVE RENEWAL ENHANCEMENT"`
	Priority          string     `json:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH URGENT"`
	ProposedStartDate *time.Time `json:"proposed_start_date"`
	ProposedEndDate   *time.Time `json:"proposed_end_date"`
	BudgetType        string     `json:"budget_type" validate:"omitempty,oneof=CAPEX OPEX"`
	BudgetName        string     `json:"budget_name" validate:"max=200"`
	Notes             string     `json:"notes" validate:"max=2000"`
}

type CreateDraftRequest struct {
	Title           string  `json:"title" validate:"required,min=5,max=200"`
	Description     string  `json:"description" validate:"max=2000"`
	BusinessGoal    string  `json:"business_goal" validate:"max=1000"`
	ExpectedOutcome string  `json:"expected_outcome" validate:"max=1000"`
	EstimatedBudget float64 `json:"estimated_budget" validate:"min=0"`
	ProjectRequestMetadata
}

type UpdateDraftRequest struct {
	Title           string  `json:"title" validate:"required,min=5,max=200"`
	Description     string  `json:"description" validate:"max=2000"`
	BusinessGoal    string  `json:"business_goal" validate:"max=1000"`
	ExpectedOutcome string  `json:"expected_outcome" validate:"max=1000"`
	EstimatedBudget float64 `json:"estimated_budget" validate:"min=0"`
	Version         int     `json:"version" validate:"required"`
	ProjectRequestMetadata
}

type ReviewRequest struct {
	Action           string  `json:"action" validate:"required,oneof=APPROVED REJECTED REQUEST_REVISION"`
	Comment          string  `json:"comment" validate:"max=1000"`
	ProjectManagerID *uint64 `json:"project_manager_id,omitempty"`
}

type ReviseRequest struct {
	Title           string  `json:"title" validate:"required,min=5,max=200"`
	Description     string  `json:"description" validate:"max=2000"`
	BusinessGoal    string  `json:"business_goal" validate:"max=1000"`
	ExpectedOutcome string  `json:"expected_outcome" validate:"max=1000"`
	EstimatedBudget float64 `json:"estimated_budget" validate:"min=0"`
	RevisionReason  string  `json:"revision_reason" validate:"required,min=5,max=1000"`
	ProjectRequestMetadata
}

type ProjectRequestResponse struct {
	ID                uint64     `json:"id"`
	RequesterID       uint64     `json:"requester_id"`
	RequestNumber     string     `json:"request_number,omitempty"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	BusinessGoal      string     `json:"business_goal"`
	ExpectedOutcome   string     `json:"expected_outcome"`
	EstimatedBudget   float64    `json:"estimated_budget"`
	Category          string     `json:"category"`
	InitiationType    *string    `json:"initiation_type,omitempty"`
	Priority          string     `json:"priority"`
	ProposedStartDate *time.Time `json:"proposed_start_date,omitempty"`
	ProposedEndDate   *time.Time `json:"proposed_end_date,omitempty"`
	BudgetType        *string    `json:"budget_type,omitempty"`
	BudgetName        string     `json:"budget_name"`
	Notes             string     `json:"notes"`
	Status            string     `json:"status"`
	CurrentRevision   int        `json:"current_revision"`
	SubmittedAt       *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	RejectedAt        *time.Time `json:"rejected_at,omitempty"`
	Version           int        `json:"version"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type RevisionResponse struct {
	ID                uint64     `json:"id"`
	RevisionNumber    int        `json:"revision_number"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	BusinessGoal      string     `json:"business_goal"`
	ExpectedOutcome   string     `json:"expected_outcome"`
	EstimatedBudget   float64    `json:"estimated_budget"`
	Category          string     `json:"category"`
	InitiationType    *string    `json:"initiation_type,omitempty"`
	Priority          string     `json:"priority"`
	ProposedStartDate *time.Time `json:"proposed_start_date,omitempty"`
	ProposedEndDate   *time.Time `json:"proposed_end_date,omitempty"`
	BudgetType        *string    `json:"budget_type,omitempty"`
	BudgetName        string     `json:"budget_name"`
	Notes             string     `json:"notes"`
	RevisionReason    string     `json:"revision_reason"`
	RevisedBy         uint64     `json:"revised_by"`
	CreatedAt         time.Time  `json:"created_at"`
}

type ApprovalResponse struct {
	ID               uint64    `json:"id"`
	ReviewedBy       uint64    `json:"reviewed_by"`
	Action           string    `json:"action"`
	Comment          string    `json:"comment"`
	ProjectManagerID *uint64   `json:"project_manager_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
