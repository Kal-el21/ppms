package dto

import "time"

type CreateDraftRequest struct {
	Title           string  `json:"title" validate:"required,min=5,max=200"`
	Description     string  `json:"description" validate:"max=2000"`
	BusinessGoal    string  `json:"business_goal" validate:"max=1000"`
	ExpectedOutcome string  `json:"expected_outcome" validate:"max=1000"`
	EstimatedBudget float64 `json:"estimated_budget" validate:"min=0"`
}

type UpdateDraftRequest struct {
	Title           string  `json:"title" validate:"required,min=5,max=200"`
	Description     string  `json:"description" validate:"max=2000"`
	BusinessGoal    string  `json:"business_goal" validate:"max=1000"`
	ExpectedOutcome string  `json:"expected_outcome" validate:"max=1000"`
	EstimatedBudget float64 `json:"estimated_budget" validate:"min=0"`
	Version         int     `json:"version" validate:"required"`
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
}

type ProjectRequestResponse struct {
	ID              uint64     `json:"id"`
	RequesterID     uint64     `json:"requester_id"`
	RequestNumber   string     `json:"request_number,omitempty"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	BusinessGoal    string     `json:"business_goal"`
	ExpectedOutcome string     `json:"expected_outcome"`
	EstimatedBudget float64    `json:"estimated_budget"`
	Status          string     `json:"status"`
	CurrentRevision int        `json:"current_revision"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt      *time.Time `json:"approved_at,omitempty"`
	RejectedAt      *time.Time `json:"rejected_at,omitempty"`
	Version         int        `json:"version"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type RevisionResponse struct {
	ID              uint64    `json:"id"`
	RevisionNumber  int       `json:"revision_number"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	EstimatedBudget float64   `json:"estimated_budget"`
	RevisionReason  string    `json:"revision_reason"`
	RevisedBy       uint64    `json:"revised_by"`
	CreatedAt       time.Time `json:"created_at"`
}

type ApprovalResponse struct {
	ID               uint64    `json:"id"`
	ReviewedBy       uint64    `json:"reviewed_by"`
	Action           string    `json:"action"`
	Comment          string    `json:"comment"`
	ProjectManagerID *uint64   `json:"project_manager_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}
