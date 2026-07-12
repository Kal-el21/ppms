package entity

import "time"

type ProjectRequestRevision struct {
	ID               uint64 `gorm:"primaryKey" json:"id"`
	ProjectRequestID uint64 `gorm:"not null" json:"project_request_id"`

	RevisionNumber int `gorm:"not null" json:"revision_number"`

	Title           string  `json:"title"`
	Description     string  `json:"description"`
	BusinessGoal    string  `json:"business_goal"`
	ExpectedOutcome string  `json:"expected_outcome"`
	EstimatedBudget float64 `json:"estimated_budget"`

	Category          string          `json:"category"`
	InitiationType    *InitiationType `json:"initiation_type,omitempty"`
	Priority          RequestPriority `gorm:"not null;default:MEDIUM" json:"priority"`
	ProposedStartDate *time.Time      `json:"proposed_start_date,omitempty"`
	ProposedEndDate   *time.Time      `json:"proposed_end_date,omitempty"`
	BudgetType        *BudgetType     `json:"budget_type,omitempty"`
	BudgetName        string          `json:"budget_name"`
	Notes             string          `json:"notes"`

	RevisionReason string `json:"revision_reason"`

	RevisedBy uint64 `gorm:"not null" json:"revised_by"`

	CreatedAt time.Time `json:"created_at"`
}

func (ProjectRequestRevision) TableName() string {
	return "project_request_revisions"
}
